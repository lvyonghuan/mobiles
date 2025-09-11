package ansible

import (
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/libp2p/go-libp2p/p2p/discovery/mdns"
	"github.com/lvyonghuan/Ubik-Util/uerr"
)

func (p *peerManager) HandlePeerFound(pi peer.AddrInfo) {
	// Lock the peer store for concurrent access
	p.mu.Lock()
	defer p.mu.Unlock()

	// Create a new ansiblePeer instance
	var asbP ansiblePeer
	asbP.addr = pi
	asbP.linkCount = 0

	// Check if the peer is already in the store
	if _, exists := p.peers[pi.ID]; exists {
		p.peers[pi.ID] = asbP // Update the existing peer info
	} else {
		p.peers[pi.ID] = asbP // Add new peer info
	}
}

func (asb *ImplAnsible) setupDiscovery() error {
	s := mdns.NewMdnsService(asb.h, "p2p-node-discovery", asb.peerStore)
	err := s.Start()
	if err != nil {
		return uerr.NewError(err)
	}
	return nil
}
