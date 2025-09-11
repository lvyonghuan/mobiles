package ansible

import (
	"github.com/libp2p/go-libp2p/core/network"
)

func (p *peerManager) handelHeartbeat(s network.Stream) {
	defer s.Close()
	remotePeer := s.Conn().RemotePeer()

	// Reset the heartbeat timeout channel to avoid disconnection
	// Get the peer form manager
	p.mu.Lock()
	defer p.mu.Unlock()

	if pr, exists := p.peers[remotePeer]; exists {
		pr.resetHeartbeatTimeout()
	} else {
		// TODO log
		return
	}
}

// The leader's first message to the follower
// So it identity who is the leader, and we can get the leader's addr
// The follower should show the leader its identity (the plugin info) in the response
func (p *peerManager) handelIdentityConfirmation(s network.Stream) {
	defer s.Close()
	remotePeer := s.Conn().RemotePeer()

	// This peer is the leader. It's a rule.
	p.ansible.setLeader(remotePeer)

	// Add a link count for the leader
	err := p.addLink(remotePeer)
	if err != nil {
		//TODO log
		return
	}

	// Send back the identity info of this follower
	err = p.sendPluginInfoToLeader(remotePeer)
	if err != nil {
		//TODO log
		return
	}

	return
}

func (p *peerManager) handelCreateWorkflow(s network.Stream) {
	defer s.Close()

	var workflowID int
	err := readFromStream(s, &workflowID)
	if err != nil {
		//TODO log
		return
	}

	// Create a workflow
	p.ansible.getRuntime().InitWorkflow(workflowID)
}

//TODO deleteWorkflow

func (p *peerManager) createNodeProtocol(s network.Stream) {
	defer s.Close()

	var message createNodeMessage
	err := readFromStream(s, &message)
	if err != nil {
		//TODO log
		return
	}

	// Create a node
	err = p.ansible.getRuntime().CreateRuntimeNode(message.NodeName, message.NodeID, message.WorkflowID)
	if err != nil {
		//TODO log
		return
	}
}

//TODO deleteNodeProtocol

//TODO setParamProtocol

func (p *peerManager) createEdgeProtocol(s network.Stream) {
	defer s.Close()

	var message createEdgeMessage
	err := readFromStream(s, &message)
	if err != nil {
		//TODO log
		return
	}

	// Create an edge
	err = p.ansible.getRuntime().CreateEdge(message.Destination, message.WorkflowID, message.ProducerNodeID, message.ProducerPortName, message.ConsumerNodeID, message.ConsumerPortName)
	if err != nil {
		//TODO log
		return
	}
}

//TODO deleteEdgeProtocol

func (p *peerManager) handelRunWorkflow(s network.Stream) {
	defer s.Close()

	var workflowID int
	err := readFromStream(s, &workflowID)
	if err != nil {
		//TODO log
		return
	}

	// Init listener
	resultChan, errorChan, processChan := p.ansible.initWorkflowListener(workflowID)

	// Run the workflow
	ctx, err := p.ansible.getRuntime().RunWorkflow(workflowID, resultChan, errorChan, processChan)
	if err != nil {
		//TODO log
		return
	}

	// Set the stop context
	wl := p.ansible.getWorkflowListener(workflowID)
	if wl == nil {
		// TODO log
		return
	}
	wl.setStopContext(ctx)

	// Run the listener
	go wl.run()
}
