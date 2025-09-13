package ansible

import (
	"testing"

	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/lvyonghuan/mobiles/hainish"
)

// mockPlugin simulates plugin implementation
type mockPlugin struct {
	name        string
	description string
	version     string
	author      string
	license     string
	nodes       map[string]hainish.Node
}

func (m *mockPlugin) Name() string {
	return m.name
}

func (m *mockPlugin) Description() string {
	return m.description
}

func (m *mockPlugin) Version() string {
	return m.version
}

func (m *mockPlugin) Author() string {
	return m.author
}

func (m *mockPlugin) License() string {
	return m.license
}

func (m *mockPlugin) Nodes() map[string]hainish.Node {
	return m.nodes
}

// TestPeerManagerCreation tests peer manager creation
func TestPeerManagerCreation(t *testing.T) {
	ansible := &ImplAnsible{}
	peerManager := &peerManager{
		peers:   make(map[peer.ID]ansiblePeer),
		ansible: ansible,
	}

	if peerManager.peers == nil {
		t.Error("Expected peers map to be initialized")
	}

	if peerManager.ansible != ansible {
		t.Error("Expected ansible reference to be set")
	}
}

// TestAnsiblePeerCreation tests ansible peer creation
func TestAnsiblePeerCreation(t *testing.T) {
	peer := ansiblePeer{
		linkCount: 0,
	}

	if peer.linkCount != 0 {
		t.Errorf("Expected linkCount 0, got %d", peer.linkCount)
	}

	if peer.ctx != nil {
		t.Error("Expected ctx to be nil initially")
	}

	if peer.cancel != nil {
		t.Error("Expected cancel to be nil initially")
	}
}

// TestWorkflowListenerCreation tests workflow listener creation
func TestWorkflowListenerCreation(t *testing.T) {
	ansible := &ImplAnsible{}
	listener := &workflowListener{
		workflowID: 1,
		ansible:    ansible,
	}

	if listener.workflowID != 1 {
		t.Errorf("Expected workflowID 1, got %d", listener.workflowID)
	}

	if listener.ansible != ansible {
		t.Error("Expected ansible reference to be set")
	}
}

// TestMessageStructures tests message structures
func TestMessageStructures(t *testing.T) {
	// Test createNodeMessage
	nodeMsg := createNodeMessage{
		NodeName:   "testNode",
		NodeID:     1,
		WorkflowID: 1,
	}

	if nodeMsg.NodeName != "testNode" {
		t.Errorf("Expected NodeName 'testNode', got '%s'", nodeMsg.NodeName)
	}

	if nodeMsg.NodeID != 1 {
		t.Errorf("Expected NodeID 1, got %d", nodeMsg.NodeID)
	}

	if nodeMsg.WorkflowID != 1 {
		t.Errorf("Expected WorkflowID 1, got %d", nodeMsg.WorkflowID)
	}

	// Test createEdgeMessage
	edgeMsg := createEdgeMessage{
		EdgeID:           1,
		Destination:      "peer123",
		WorkflowID:       1,
		ProducerNodeID:   1,
		ProducerPortName: "output1",
		ConsumerNodeID:   2,
		ConsumerPortName: "input1",
	}

	if edgeMsg.EdgeID != 1 {
		t.Errorf("Expected EdgeID 1, got %d", edgeMsg.EdgeID)
	}

	if edgeMsg.ProducerPortName != "output1" {
		t.Errorf("Expected ProducerPortName 'output1', got '%s'", edgeMsg.ProducerPortName)
	}

	// Test logMessage
	logMsg := logMessage{
		Level:   1,
		Message: "Test log message",
	}

	if logMsg.Level != 1 {
		t.Errorf("Expected Level 1, got %d", logMsg.Level)
	}

	if logMsg.Message != "Test log message" {
		t.Errorf("Expected Message 'Test log message', got '%s'", logMsg.Message)
	}

	// Test setParamMessage
	paramMsg := setParamMessage{
		WorkflowID: 1,
		NodeID:     1,
		ParamName:  "param1",
		ParamValue: "test_value",
	}

	if paramMsg.ParamName != "param1" {
		t.Errorf("Expected ParamName 'param1', got '%s'", paramMsg.ParamName)
	}

	if paramMsg.ParamValue != "test_value" {
		t.Errorf("Expected ParamValue 'test_value', got '%v'", paramMsg.ParamValue)
	}
}

// TestProtocolConstants tests protocol constants
func TestProtocolConstants(t *testing.T) {
	// Test all protocol constants are defined correctly
	expectedProtocols := map[string]string{
		"heartbeat":      "/ansible/heartbeat/1.0.0",
		"identity":       "/ansible/leader/identity/1.0.0",
		"createWorkflow": "/ansible/leader/workflow/create/1.0.0",
		"deleteWorkflow": "/ansible/leader/workflow/delete/1.0.0",
		"createNode":     "/ansible/leader/node/create/1.0.0",
		"deleteNode":     "/ansible/leader/node/delete/1.0.0",
		"setParam":       "/ansible/leader/node/param/1.0.0",
		"createEdge":     "/ansible/leader/edge/create/1.0.0",
		"deleteEdge":     "/ansible/leader/edge/delete/1.0.0",
		"runWorkflow":    "ansible/leader/workflow/run/1.0.0",
		"stopWorkflow":   "ansible/leader/workflow/stop/1.0.0",
		"logUpload":      "ansible/follower/log/1.0.0",
		"resultUpload":   "ansible/follower/result/1.0.0",
		"passingData":    "ansible/follower/data/1.0.0",
	}

	actualProtocols := map[string]string{
		"heartbeat":      heartbeatProtocol,
		"identity":       identityConfirmationProtocol,
		"createWorkflow": createWorkflow,
		"deleteWorkflow": deleteWorkflow,
		"createNode":     createNodeProtocol,
		"deleteNode":     deleteNodeProtocol,
		"setParam":       setParamProtocol,
		"createEdge":     createEdgeProtocol,
		"deleteEdge":     deleteEdgeProtocol,
		"runWorkflow":    runWorkflowProtocol,
		"stopWorkflow":   stopWorkflowProtocol,
		"logUpload":      logUploadProtocol,
		"resultUpload":   resultUploadProtocol,
		"passingData":    passingDataProtocol,
	}

	for name, expected := range expectedProtocols {
		actual := actualProtocols[name]
		if actual != expected {
			t.Errorf("Protocol %s: expected '%s', got '%s'", name, expected, actual)
		}
	}
}

// TestWorkflowListenerInterface tests workflow listener interface
func TestWorkflowListenerInterface(t *testing.T) {
	listener := &workflowListener{}

	// Test initial values
	if listener.workflowID != 0 {
		t.Errorf("Expected initial workflowID 0, got %d", listener.workflowID)
	}

	if listener.resultChan != nil {
		t.Error("Expected initial resultChan to be nil")
	}

	if listener.errChan != nil {
		t.Error("Expected initial errChan to be nil")
	}

	if listener.processChan != nil {
		t.Error("Expected initial processChan to be nil")
	}
}

// TestInitWorkflowListener tests workflow listener initialization
func TestInitWorkflowListener(t *testing.T) {
	ansible := &ImplAnsible{}

	resultChan, errChan, processChan := ansible.initWorkflowListener(1)

	if resultChan == nil {
		t.Error("Expected resultChan to be initialized")
	}

	if errChan == nil {
		t.Error("Expected errChan to be initialized")
	}

	if processChan == nil {
		t.Error("Expected processChan to be initialized")
	}

	// Check channel capacities
	if cap(resultChan) < 1 {
		t.Error("Expected resultChan to have capacity")
	}
}
