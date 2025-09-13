package ansible

import (
	"context"
	"sync"

	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/lvyonghuan/Ubik-Util/uerr"
	"github.com/lvyonghuan/mobiles/hainish"
	"github.com/lvyonghuan/mobiles/runtime"
	"github.com/lvyonghuan/mobiles/util"
)

type Ansible interface {
	host() host.Host
	setupDiscovery() error
	setProtocolAndHandel()

	setLeader(peerID peer.ID)
	getLeader() peer.ID
	getPluginMetadata() hainish.Plugin
	getRuntime() *runtime.Runtime

	initWorkflowListener(workflowID int) (chan any, chan error, chan hainish.Edge)
	getWorkflowListener(workflowID int) *workflowListener
	getPeerManager() *peerManager
}

type ImplAnsible struct {
	h         host.Host
	peerStore *peerManager

	leaderID peer.ID

	pluginMetadata hainish.Plugin
	r              *runtime.Runtime

	wfListener map[int]workflowListener
}

func (asb *ImplAnsible) host() host.Host {
	return asb.h
}

type peerManager struct {
	peers map[peer.ID]ansiblePeer
	mu    sync.Mutex

	ansible Ansible
}

type ansiblePeer struct {
	addr peer.AddrInfo

	linkCount int // Number of active connections.
	ctx       context.Context
	cancel    context.CancelFunc

	resetHeartbeatTimeoutChan chan struct{}
}

func Init(pluginMetadata hainish.Plugin, runtime *runtime.Runtime) (*Ansible, error) {
	var ansible ImplAnsible

	ansible.pluginMetadata = pluginMetadata
	ansible.r = runtime

	// Initialize libp2p host
	h, err := initLibp2p()
	if err != nil {
		return nil, err
	}
	ansible.h = h

	// Initialize peer manager
	ansible.peerStore = &peerManager{
		peers:   make(map[peer.ID]ansiblePeer),
		mu:      sync.Mutex{},
		ansible: &ansible,
	}
	// Add self to peer store
	ansible.peerStore.peers[h.ID()] = ansiblePeer{
		addr:      peer.AddrInfo{ID: h.ID(), Addrs: h.Addrs()},
		linkCount: 0,
	}

	// Start peer discovery
	err = ansible.setupDiscovery()
	if err != nil {
		return nil, err
	}

	//Set protocol and handle
	ansible.setProtocolAndHandel()

	return nil, err
}

func initLibp2p() (host.Host, error) {
	h, err := libp2p.New()
	if err != nil {
		return nil, uerr.NewError(err)
	}

	return h, nil
}

// Set protocol and handle functions
func (asb *ImplAnsible) setProtocolAndHandel() {
	// Heartbeat protocol
	asb.h.SetStreamHandler(heartbeatProtocol, asb.peerStore.handelHeartbeat)
	// Identity confirmation protocol
	asb.h.SetStreamHandler(identityConfirmationProtocol, asb.peerStore.handelIdentityConfirmation)
	// Create workflow protocol
	asb.h.SetStreamHandler(createWorkflow, asb.peerStore.handelCreateWorkflow)
	// Create node protocol
	asb.h.SetStreamHandler(createNodeProtocol, asb.peerStore.handelCreateNodeProtocol)
	// Create edge protocol
	asb.h.SetStreamHandler(createEdgeProtocol, asb.peerStore.handelCreateEdgeProtocol)
	// Passing data protocol
	asb.h.SetStreamHandler(runWorkflowProtocol, asb.peerStore.handelRunWorkflow)
	// Delete workflow protocol
	asb.h.SetStreamHandler(deleteWorkflow, asb.peerStore.handelDeleteWorkflow)
	// Delete node protocol
	asb.h.SetStreamHandler(deleteNodeProtocol, asb.peerStore.handelDeleteNodeProtocol)
	// Set param protocol
	asb.h.SetStreamHandler(setParamProtocol, asb.peerStore.handelSetParamProtocol)
	// Delete edge protocol
	asb.h.SetStreamHandler(deleteEdgeProtocol, asb.peerStore.handelDeleteEdgeProtocol)
	// Stop workflow protocol
	asb.h.SetStreamHandler(stopWorkflowProtocol, asb.peerStore.handelStopWorkflow)
}

// Add link based on workflow as the scale
// Or if the peer is the leader
func (p *peerManager) addLink(peerID peer.ID) error {
	p.mu.Lock()
	defer p.mu.Unlock()
	if pr, exists := p.peers[peerID]; exists {
		pr.linkCount++
		// If this is the first link, start heartbeat
		if pr.linkCount == 1 {
			// Make a context to stop the heartbeat when needed
			pr.ctx, pr.cancel = context.WithCancel(context.Background())
			pr.resetHeartbeatTimeoutChan = make(chan struct{}, 1)
			// Start heartbeat goroutine
			go p.startSendHeartbeat(pr)
		}
		p.peers[peerID] = pr
	} else {
		return uerr.NewError(util.ErrPeerNotExist)
	}
	return nil
}

// Subtract link based on workflow as the scale
func (p *peerManager) subLink(peerID peer.ID) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	if pr, exists := p.peers[peerID]; exists {
		pr.linkCount--
		// If no more links, stop heartbeat
		if pr.linkCount <= 0 {
			pr.linkCount = 0
			pr.cancel()
			// Close the reset channel to avoid goroutine leak
			close(pr.resetHeartbeatTimeoutChan)
			pr.resetHeartbeatTimeoutChan = nil // Avoid further use
		}
		p.peers[peerID] = pr
	} else {
		return uerr.NewError(util.ErrPeerNotExist)
	}
	return nil
}

func (asb *ImplAnsible) setLeader(peerID peer.ID) {
	asb.leaderID = peerID
}

func (asb *ImplAnsible) getLeader() peer.ID {
	return asb.leaderID
}

func (asb *ImplAnsible) getPluginMetadata() hainish.Plugin {
	return asb.pluginMetadata
}

func (asb *ImplAnsible) getRuntime() *runtime.Runtime {
	return asb.r
}

func (asb *ImplAnsible) getPeerManager() *peerManager {
	return asb.peerStore
}
