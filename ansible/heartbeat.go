package ansible

import (
	"time"
)

func (p *peerManager) startSendHeartbeat(asbPeer ansiblePeer) {
	// Start a ticker to send heartbeat messages every 30 seconds
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			err := p.sendHeartbeat(asbPeer.addr.ID)
			if err != nil {
				// TODO: Handle send heartbeat error
				continue
			}
		case <-asbPeer.ctx.Done():
			// Stop the heartbeat goroutine if the context is cancelled
			return
		}
	}
}

func (p *peerManager) startListenHeartbeat(asbPeer ansiblePeer) error {
	// Start a ticker to listen for heartbeat messages every 90 seconds
	ticker := time.NewTicker(90 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			// If no heartbeat received in the last 90 seconds, consider the peer offline
			// TODO: Handle peer offline logic
		case <-p.peers[asbPeer.addr.ID].resetHeartbeatTimeoutChan:
			// Reset the heartbeat timeout if a heartbeat message is received
			ticker.Reset(90 * time.Second)
		case <-asbPeer.ctx.Done():
			// Stop the listening goroutine if the context is cancelled
			return nil
		}
	}
}

func (asbPeer *ansiblePeer) resetHeartbeatTimeout() {
	//Check the channel's status
	if asbPeer.resetHeartbeatTimeoutChan == nil {
		return
	}

	// Reset the heartbeat timeout ticker
	asbPeer.resetHeartbeatTimeoutChan <- struct{}{}
	return
}
