package hainish

import "github.com/libp2p/go-libp2p/core/peer"

// Edge No matter who the sender is, only the delivery matters.
type Edge struct {
	Destination      peer.ID // Where to go (the city)
	TargetWorkflowID int     `json:"TargetWorkflowID"` // Which street
	TargetNodeID     int     `json:"TargetNodeID"`     // Which building
	TargetPort       string  `json:"TargetPort"`       // Which door
	Value            any     `json:"Value"`            // What to send
}

func NewEdge(destination peer.ID, targetWorkflowID, targetNodeID int, targetPort string) *Edge {
	// envelope
	return &Edge{
		Destination:      destination,
		TargetWorkflowID: targetWorkflowID,
		TargetNodeID:     targetNodeID,
		TargetPort:       targetPort,
	}
}
