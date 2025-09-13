package ansible

import "github.com/libp2p/go-libp2p/core/peer"

type createNodeMessage struct {
	NodeName   string `json:"NodeName"`
	NodeID     int    `json:"NodeID"`
	WorkflowID int    `json:"WorkflowID"`
}

type createEdgeMessage struct {
	EdgeID           int     `json:"EdgeID"`
	Destination      peer.ID `json:"destination"`
	WorkflowID       int     `json:"WorkflowID"`
	ProducerNodeID   int     `json:"producerNodeID"`
	ProducerPortName string  `json:"producerPortName"`
	ConsumerNodeID   int     `json:"consumerNodeID"`
	ConsumerPortName string  `json:"consumerPortName"`
}

type logMessage struct {
	Level   int    `json:"level"`
	Message string `json:"message"`
}

type deleteNodeMessage struct {
	WorkflowID int `json:"WorkflowID"`
	NodeID     int `json:"NodeID"`
}

type setParamMessage struct {
	WorkflowID int    `json:"WorkflowID"`
	NodeID     int    `json:"NodeID"`
	ParamName  string `json:"ParamName"`
	ParamValue any    `json:"ParamValue"`
}

type deleteEdgeMessage struct {
	WorkflowID int `json:"WorkflowID"`
	EdgeID     int `json:"EdgeID"`
}
