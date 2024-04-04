package mempool

import (
	"errors"
	"fmt"
	"sync/atomic"
	"time"

	"golang.org/x/sync/semaphore"

	abci "github.com/cometbft/cometbft/abci/types"
	protomem "github.com/cometbft/cometbft/api/cometbft/mempool/v1"
	cfg "github.com/cometbft/cometbft/config"
	"github.com/cometbft/cometbft/internal/clist"
	cmtsync "github.com/cometbft/cometbft/internal/sync"
	"github.com/cometbft/cometbft/libs/log"
	"github.com/cometbft/cometbft/p2p"
	"github.com/cometbft/cometbft/types"
)

// Reactor handles mempool tx broadcasting amongst peers.
// It maintains a map from peer ID to counter, to prevent gossiping txs to the
// peers you received it from.
type Reactor struct {
	p2p.BaseReactor
	config  *cfg.MempoolConfig
	mempool *CListMempool

	waitSync   atomic.Bool
	waitSyncCh chan struct{} // for signaling when to start receiving and sending txs

	// `txSenders` maps every received transaction to the set of peer IDs that
	// have sent the transaction to this node. Sender IDs are used during
	// transaction propagation to avoid sending a transaction to a peer that
	// already has it.
	txSenders    map[types.TxKey]map[p2p.ID]bool
	txSendersMtx cmtsync.Mutex

	broadcastRoutins atomic.Int32 // number of broadcast rountine started

	// similar to txSenders, but is set before a tx is validated, and is removed no maater the result
	txSendersUnchecked                map[types.TxKey]map[p2p.ID]struct{}
	txSendersUncheckedMtx             cmtsync.Mutex
	txSendersUncheckedRemoveThreshold map[types.TxKey]int32 // how many broadcast routines to wait before removing the tx from the unchecked map
	txSendersUncheckedRemoveCount     map[types.TxKey]int32 // how many broadcast routines has finished chekcing the tx unchecked map
	txSendersUncheckedRemoveCountMtx  cmtsync.Mutex

	// Semaphores to keep track of how many connections to peers are active for broadcasting
	// transactions. Each semaphore has a capacity that puts an upper bound on the number of
	// connections for different groups of peers.
	activePersistentPeersSemaphore    *semaphore.Weighted
	activeNonPersistentPeersSemaphore *semaphore.Weighted

	// record all peers, so we know whether a TXs comes from a peer
	peers *p2p.PeerSet
}

// NewReactor returns a new Reactor with the given config and mempool.
func NewReactor(config *cfg.MempoolConfig, mempool *CListMempool, waitSync bool) *Reactor {
	memR := &Reactor{
		config:                            config,
		mempool:                           mempool,
		waitSync:                          atomic.Bool{},
		txSenders:                         make(map[types.TxKey]map[p2p.ID]bool),
		txSendersUnchecked:                make(map[types.TxKey]map[p2p.ID]struct{}),
		txSendersUncheckedRemoveThreshold: make(map[types.TxKey]int32),
		txSendersUncheckedRemoveCount:     make(map[types.TxKey]int32),
		peers:                             p2p.NewPeerSet(), // initialize an empty peerSet
	}
	memR.BaseReactor = *p2p.NewBaseReactor("Mempool", memR)
	if waitSync {
		memR.waitSync.Store(true)
		memR.waitSyncCh = make(chan struct{})
	}
	memR.mempool.SetTxRemovedCallback(func(txKey types.TxKey) { memR.removeSenders(txKey) })
	memR.activePersistentPeersSemaphore = semaphore.NewWeighted(int64(memR.config.ExperimentalMaxGossipConnectionsToPersistentPeers))
	memR.activeNonPersistentPeersSemaphore = semaphore.NewWeighted(int64(memR.config.ExperimentalMaxGossipConnectionsToNonPersistentPeers))

	return memR
}

// SetLogger sets the Logger on the reactor and the underlying mempool.
func (memR *Reactor) SetLogger(l log.Logger) {
	memR.Logger = l
	memR.mempool.SetLogger(l)
}

// OnStart implements p2p.BaseReactor.
func (memR *Reactor) OnStart() error {
	if memR.WaitSync() {
		memR.Logger.Info("Starting reactor in sync mode: tx propagation will start once sync completes")
	}
	if !memR.config.Broadcast {
		memR.Logger.Info("Tx broadcasting is disabled")
	}
	return nil
}

// GetChannels implements Reactor by returning the list of channels for this
// reactor.
func (memR *Reactor) GetChannels() []*p2p.ChannelDescriptor {
	largestTx := make([]byte, memR.config.MaxTxBytes)
	batchMsg := protomem.Message{
		Sum: &protomem.Message_Txs{
			Txs: &protomem.Txs{Txs: [][]byte{largestTx}},
		},
	}

	return []*p2p.ChannelDescriptor{
		{
			ID:                  MempoolChannel,
			Priority:            5,
			RecvMessageCapacity: batchMsg.Size(),
			MessageType:         &protomem.Message{},
		},
	}
}

// AddPeer implements Reactor.
// It starts a broadcast routine ensuring all txs are forwarded to the given peer.
func (memR *Reactor) AddPeer(peer p2p.Peer) {
	if memR.config.Broadcast {
		go func() {
			memR.mempool.metrics.ActiveOutboundConnections.Add(1)
			defer memR.mempool.metrics.ActiveOutboundConnections.Add(-1)
			memR.broadcastRoutins.Add(1)
			memR.broadcastTxRoutine(peer)
			memR.broadcastRoutins.Add(-1)
		}()

		if err := memR.peers.Add(peer); err != nil {
			memR.Logger.Error("")
		}
	}
}

// Receive implements Reactor.
// It adds any received transactions to the mempool.
func (memR *Reactor) Receive(e p2p.Envelope) {
	memR.Logger.Debug("Receive", "src", "e.Src", "chId", e.ChannelID, "msg", "e.Message")
	switch msg := e.Message.(type) {
	case *protomem.Txs:
		if memR.WaitSync() {
			memR.Logger.Debug("Ignored message received while syncing", "msg", msg)
			return
		}

		protoTxs := msg.GetTxs()
		if len(protoTxs) == 0 {
			memR.Logger.Error("received empty txs from peer", "src", e.Src)
			return
		}

		for _, txBytes := range protoTxs {
			tx := types.Tx(txBytes)

			// first add sender preset
			memR.addSenderUnchecked(tx.Key(), e.Src.ID())
			// memR.txSendersUncheckedRemoveThreshold[tx.Key()] = memR.broadcastRoutins.Load()

			// check if the src is our peer
			reqRes, err := memR.mempool.CheckTx(tx)
			switch {
			case errors.Is(err, ErrTxInCache):
				memR.Logger.Debug("Tx already exists in cache", "tx", "tx.Hash()")
			case err != nil:
				memR.Logger.Info("Could not check tx", "tx", "tx.Hash()", "err", err)
			default:
				// Record the sender only when the transaction is valid and, as
				// a consequence, added to the mempool. Senders are stored until
				// the transaction is removed from the mempool. Note that it's
				// possible a tx is still in the cache but no longer in the
				// mempool. For example, after committing a block, txs are
				// removed from mempool but not the cache.
				reqRes.SetCallback(func(res *abci.Response) {
					if res.GetCheckTx().Code == abci.CodeTypeOK {
						memR.addSender(tx.Key(), e.Src.ID())
					}
				})
			}
		}
	default:
		memR.Logger.Error("unknown message type", "src", e.Src, "chId", e.ChannelID, "msg", e.Message)
		memR.Switch.StopPeerForError(e.Src, fmt.Errorf("mempool cannot handle message of type: %T", e.Message))
		return
	}

	// broadcasting happens from go routines per peer
}

func (memR *Reactor) EnableInOutTxs() {
	memR.Logger.Info("enabling inbound and outbound transactions")
	if !memR.waitSync.CompareAndSwap(true, false) {
		return
	}

	// Releases all the blocked broadcastTxRoutine instances.
	if memR.config.Broadcast {
		close(memR.waitSyncCh)
	}
}

func (memR *Reactor) WaitSync() bool {
	return memR.waitSync.Load()
}

// PeerState describes the state of a peer.
type PeerState interface {
	GetHeight() int64
}

// Send new mempool txs to peer.
func (memR *Reactor) broadcastTxRoutine(peer p2p.Peer) {
	var next *clist.CElement

	// If the node is catching up, don't start this routine immediately.
	if memR.WaitSync() {
		select {
		case <-memR.waitSyncCh:
			// EnableInOutTxs() has set WaitSync() to false.
		case <-memR.Quit():
			return
		}
	}

	for {
		// In case of both next.NextWaitChan() and peer.Quit() are variable at the same time
		if !memR.IsRunning() || !peer.IsRunning() {
			return
		}

		// This happens because the CElement we were looking at got garbage
		// collected (removed). That is, .NextWait() returned nil. Go ahead and
		// start from the beginning.
		if next == nil {
			select {
			case <-memR.mempool.TxsWaitChan(): // Wait until a tx is available
				if next = memR.mempool.TxsFront(); next == nil {
					continue
				}
			case <-peer.Quit():
				return
			case <-memR.Quit():
				return
			}
		}

		// Make sure the peer is up to date.
		peerState, ok := peer.Get(types.PeerStateKey).(PeerState)
		if !ok {
			// Peer does not have a state yet. We set it in the consensus reactor, but
			// when we add peer in Switch, the order we call reactors#AddPeer is
			// different every time due to us using a map. Sometimes other reactors
			// will be initialized before the consensus reactor. We should wait a few
			// milliseconds and retry.
			time.Sleep(PeerCatchupSleepIntervalMS * time.Millisecond)
			continue
		}

		// If we suspect that the peer is lagging behind, at least by more than
		// one block, we don't send the transaction immediately. This code
		// reduces the mempool size and the recheck-tx rate of the receiving
		// node. See [RFC 103] for an analysis on this optimization.
		//
		// [RFC 103]: https://github.com/cometbft/cometbft/pull/735
		memTx := next.Value.(*mempoolTx)
		if peerState.GetHeight() < memTx.Height()-1 {
			time.Sleep(PeerCatchupSleepIntervalMS * time.Millisecond)
			continue
		}

		// check if memTx comes from a peer
		isFromPeer := false
		var senders p2p.ID

		// memR.Logger.Debug("checking if tx comes from a peer", "tx", memTx.tx.Hash()[:8])

		memR.txSendersUncheckedMtx.Lock()
		txSenderUnchecked := memR.txSendersUnchecked[memTx.tx.Key()]
		// memR.Logger.Debug("number of senders", "num", len(txSenderUnchecked), "tx", memTx.tx.Hash()[:8])
		for peerID := range txSenderUnchecked {
			senders += peerID + ","
			if memR.peers.Has(peerID) {
				isFromPeer = true
				break
			}
		}
		memR.txSendersUncheckedMtx.Unlock()

		// add txSendersUncheckedRemoveCount by 1
		removeCount := memR.increaseSendersUncheckedRemoveCount(memTx.tx.Key())
		// remove sendersUnchecked if removeCount >= threshold
		// if removeCount >= memR.txSendersUncheckedRemoveThreshold[memTx.tx.Key()] {
		if removeCount >= 3 { // fixed threshold since we are only testing for 4 nodes
			memR.removeSendersUnchecked(memTx.tx.Key())
		}

		// if isFromPeer {
		// 	memR.Logger.Debug("tx comes from a peer, skip sending",
		// 		"tx", memTx.tx.Hash()[:8], "senders", senders)
		// } else {
		// 	// print all peers from memR
		// 	memR.Logger.Debug("tx comes from a non-peer, keep sending",
		// 		"tx", memTx.tx.Hash()[:8], "senders", senders)
		// 	memR.peers.ForEach(func(peer p2p.Peer) {
		// 		memR.Logger.Debug("peer", "peer", peer.ID())
		// 	})
		// }

		// NOTE: Transaction batching was disabled due to
		// https://github.com/tendermint/tendermint/issues/5796

		if !memR.isSender(memTx.tx.Key(), peer.ID()) && !isFromPeer {

			// memR.Logger.Info("sending tx to peer", "peer", peer.ID(),
			// 	"tx", memTx.tx.Hash()[:8], "height", memTx.Height(),
			// 	"txs", memR.mempool.Size(), "peerHeight", peerState.GetHeight(),
			// 	"peerPersistent", peer.IsPersistent())

			success := peer.Send(p2p.Envelope{
				ChannelID: MempoolChannel,
				Message:   &protomem.Txs{Txs: [][]byte{memTx.tx}},
			})
			if !success {
				time.Sleep(PeerCatchupSleepIntervalMS * time.Millisecond)
				continue
			}
		}

		select {
		case <-next.NextWaitChan():
			// see the start of the for loop for nil check
			next = next.Next()
		case <-peer.Quit():
			return
		case <-memR.Quit():
			return
		}
	}
}

func (memR *Reactor) isSender(txKey types.TxKey, peerID p2p.ID) bool {
	memR.txSendersMtx.Lock()
	defer memR.txSendersMtx.Unlock()

	sendersSet, ok := memR.txSenders[txKey]
	return ok && sendersSet[peerID]
}

func (memR *Reactor) addSender(txKey types.TxKey, senderID p2p.ID) {
	memR.txSendersMtx.Lock()
	defer memR.txSendersMtx.Unlock()

	if sendersSet, ok := memR.txSenders[txKey]; ok {
		sendersSet[senderID] = true
		return
	}
	memR.txSenders[txKey] = map[p2p.ID]bool{senderID: true}
}

func (memR *Reactor) removeSenders(txKey types.TxKey) {
	memR.txSendersMtx.Lock()
	defer memR.txSendersMtx.Unlock()

	if memR.txSenders != nil {
		delete(memR.txSenders, txKey)
	}
}

func (memR *Reactor) addSenderUnchecked(txKey types.TxKey, senderID p2p.ID) {
	memR.txSendersUncheckedMtx.Lock()
	defer memR.txSendersUncheckedMtx.Unlock()

	// memR.Logger.Debug("adding sender unchecked", "tx", txKey, "sender", senderID)

	if sendersSet, ok := memR.txSendersUnchecked[txKey]; ok {
		sendersSet[senderID] = struct{}{}
		return
	}
	memR.txSendersUnchecked[txKey] = map[p2p.ID]struct{}{senderID: struct{}{}}

	// memR.Logger.Debug("added sender unchecked", "tx", txKey, "sender", senderID)
}

func (memR *Reactor) removeSendersUnchecked(txKey types.TxKey) {
	memR.txSendersUncheckedMtx.Lock()
	defer memR.txSendersUncheckedMtx.Unlock()

	if memR.txSendersUnchecked != nil {
		delete(memR.txSendersUnchecked, txKey)
	}
}

func (memR *Reactor) increaseSendersUncheckedRemoveCount(txKey types.TxKey) int32 {
	memR.txSendersUncheckedRemoveCountMtx.Lock()
	defer memR.txSendersUncheckedRemoveCountMtx.Unlock()

	memR.txSendersUncheckedRemoveCount[txKey]++
	return memR.txSendersUncheckedRemoveCount[txKey]
}
