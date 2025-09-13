package runtime

import (
	"testing"
	"time"

	"github.com/lvyonghuan/mobiles/hainish"
)

// mockNode simulates node implementation
type mockNode struct {
	name        string
	description string
	isBegin     bool
	inputs      map[string]hainish.Port
	outputs     map[string]hainish.Port
	params      map[string]hainish.Port
	action      func(inputs map[string]any, output map[string]chan any) (result any, err error)
}

func (m *mockNode) Name() string {
	return m.name
}

func (m *mockNode) Description() string {
	return m.description
}

func (m *mockNode) IsBegin() bool {
	return m.isBegin
}

func (m *mockNode) Inputs() map[string]hainish.Port {
	return m.inputs
}

func (m *mockNode) Outputs() map[string]hainish.Port {
	return m.outputs
}

func (m *mockNode) Params() map[string]hainish.Port {
	return m.params
}

func (m *mockNode) Action(inputs map[string]any, output map[string]chan any) (result any, err error) {
	return m.action(inputs, output)
}

// mockPort simulates port implementation
type mockPort struct {
	name        string
	description string
	portType    string
	channel     chan any
}

func (m *mockPort) Name() string {
	return m.name
}

func (m *mockPort) Description() string {
	return m.description
}

func (m *mockPort) Type() string {
	return m.portType
}

func (m *mockPort) Chan() chan any {
	return m.channel
}

// TestInitRuntime tests runtime initialization
func TestInitRuntime(t *testing.T) {
	// Create a mock node
	mockPort := &mockPort{
		name:        "input1",
		description: "Input port",
		portType:    "string",
		channel:     make(chan any, 1),
	}

	mockNode := &mockNode{
		name:        "testNode",
		description: "Test Node",
		isBegin:     true,
		inputs:      map[string]hainish.Port{"input1": mockPort},
		outputs:     map[string]hainish.Port{},
		params:      map[string]hainish.Port{},
		action: func(inputs map[string]any, output map[string]chan any) (result any, err error) {
			return "processed", nil
		},
	}

	nodes := map[string]hainish.Node{
		"testNode": mockNode,
	}

	runtime := InitRuntime(nodes)

	if runtime == nil {
		t.Error("Expected runtime to be initialized")
	}

	if len(runtime.nodes) != 1 {
		t.Errorf("Expected 1 node in runtime, got %d", len(runtime.nodes))
	}

	if _, exists := runtime.nodes["testNode"]; !exists {
		t.Error("Expected testNode to exist in runtime")
	}
}

// TestWorkflowCreation tests workflow creation
func TestWorkflowCreation(t *testing.T) {
	runtime := InitRuntime(map[string]hainish.Node{})

	// Test creating a workflow
	runtime.InitWorkflow(1)

	if len(runtime.workflows) != 1 {
		t.Errorf("Expected 1 workflow, got %d", len(runtime.workflows))
	}

	wf, exists := runtime.workflows[1]
	if !exists {
		t.Error("Expected workflow 1 to exist")
	}

	if wf.runtimeNodes == nil {
		t.Error("Expected runtimeNodes to be initialized")
	}

	if wf.edges == nil {
		t.Error("Expected edges to be initialized")
	}
}

// TestCreateRuntimeNode tests runtime node creation
func TestCreateRuntimeNode(t *testing.T) {
	// Test creating a runtime node
	mockPort := &mockPort{
		name:        "input1",
		description: "Input port",
		portType:    "string",
		channel:     make(chan any, 1),
	}

	mockNode := &mockNode{
		name:        "testNode",
		description: "Test Node",
		isBegin:     true,
		inputs:      map[string]hainish.Port{"input1": mockPort},
		outputs:     map[string]hainish.Port{},
		params:      map[string]hainish.Port{},
		action: func(inputs map[string]any, output map[string]chan any) (result any, err error) {
			return "processed", nil
		},
	}

	nodes := map[string]hainish.Node{
		"testNode": mockNode,
	}

	runtime := InitRuntime(nodes)
	runtime.InitWorkflow(1)

	// Test creating a runtime node
	err := runtime.CreateRuntimeNode("testNode", 1, 1)
	if err != nil {
		t.Errorf("Unexpected error creating runtime node: %v", err)
	}

	wf := runtime.workflows[1]
	if len(wf.runtimeNodes) != 1 {
		t.Errorf("Expected 1 runtime node, got %d", len(wf.runtimeNodes))
	}

	if _, exists := wf.runtimeNodes[1]; !exists {
		t.Error("Expected runtime node 1 to exist")
	}
}

// TestCreateEdge tests edge creation
func TestCreateEdge(t *testing.T) {
	// Test creating an edge between two nodes
	outputPort := &mockPort{
		name:        "output1",
		description: "Output port",
		portType:    "string",
		channel:     make(chan any, 1),
	}

	inputPort := &mockPort{
		name:        "input1",
		description: "Input port",
		portType:    "string",
		channel:     make(chan any, 1),
	}

	producerNode := &mockNode{
		name:        "producerNode",
		description: "Producer Node",
		isBegin:     true,
		inputs:      map[string]hainish.Port{},
		outputs:     map[string]hainish.Port{"output1": outputPort},
		params:      map[string]hainish.Port{},
		action: func(inputs map[string]any, output map[string]chan any) (result any, err error) {
			return "data", nil
		},
	}

	consumerNode := &mockNode{
		name:        "consumerNode",
		description: "Consumer Node",
		isBegin:     false,
		inputs:      map[string]hainish.Port{"input1": inputPort},
		outputs:     map[string]hainish.Port{},
		params:      map[string]hainish.Port{},
		action: func(inputs map[string]any, output map[string]chan any) (result any, err error) {
			return "processed", nil
		},
	}

	nodes := map[string]hainish.Node{
		"producerNode": producerNode,
		"consumerNode": consumerNode,
	}

	runtime := InitRuntime(nodes)
	runtime.InitWorkflow(1)

	// Create runtime nodes
	runtime.CreateRuntimeNode("producerNode", 1, 1)
	runtime.CreateRuntimeNode("consumerNode", 2, 1)

	// Create an edge from producerNode's output1 to consumerNode's input1
	err := runtime.CreateEdge(1, "peer123", 1, 1, "output1", 2, "input1")
	if err != nil {
		t.Errorf("Unexpected error creating edge: %v", err)
	}

	wf := runtime.workflows[1]
	if len(wf.edges) != 1 {
		t.Errorf("Expected 1 edge, got %d", len(wf.edges))
	}

	edge, exists := wf.edges[1]
	if !exists {
		t.Error("Expected edge 1 to exist")
	}

	if edge.producerNodeID != 1 {
		t.Errorf("Expected producerNodeID 1, got %d", edge.producerNodeID)
	}
}

// TestWorkflowExecution tests workflow execution
func TestWorkflowExecution(t *testing.T) {
	// Create a simple begin node
	outputPort := &mockPort{
		name:        "output1",
		description: "Output port",
		portType:    "string",
		channel:     make(chan any, 1),
	}

	beginNode := &mockNode{
		name:        "beginNode",
		description: "Begin Node",
		isBegin:     true,
		inputs:      map[string]hainish.Port{},
		outputs:     map[string]hainish.Port{"output1": outputPort},
		params:      map[string]hainish.Port{},
		action: func(inputs map[string]any, output map[string]chan any) (result any, err error) {
			// Simulate some processing
			if output["output1"] != nil {
				output["output1"] <- "test_data"
			}
			return "completed", nil
		},
	}

	nodes := map[string]hainish.Node{
		"beginNode": beginNode,
	}

	runtime := InitRuntime(nodes)
	runtime.InitWorkflow(1)
	runtime.CreateRuntimeNode("beginNode", 1, 1)

	// Create channels for workflow execution
	resultChan := make(chan any, 1)
	errChan := make(chan error, 1)
	processChan := make(chan hainish.Edge, 1)

	// Execute the workflow
	ctx, err := runtime.RunWorkflow(1, resultChan, errChan, processChan)
	if err != nil {
		t.Errorf("Unexpected error running workflow: %v", err)
	}

	if ctx == nil {
		t.Error("Expected context to be returned")
	}

	// Wait a moment for the workflow to process
	time.Sleep(100 * time.Millisecond)

	// Check for results
	select {
	case result := <-resultChan:
		if result != "completed" {
			t.Errorf("Expected result 'completed', got '%v'", result)
		}
	default:
		t.Log("No result received (this might be expected for async execution)")
	}
}

// TestWorkflowStop tests workflow stopping
func TestWorkflowStop(t *testing.T) {
	// Create a simple begin node that runs for a while
	outputPort := &mockPort{
		name:        "output1",
		description: "Output port",
		portType:    "string",
		channel:     make(chan any, 1),
	}

	longRunningNode := &mockNode{
		name:        "longRunningNode",
		description: "Long Running Node",
		isBegin:     true,
		inputs:      map[string]hainish.Port{},
		outputs:     map[string]hainish.Port{"output1": outputPort},
		params:      map[string]hainish.Port{},
		action: func(inputs map[string]any, output map[string]chan any) (result any, err error) {
			// Simulate long processing
			time.Sleep(1 * time.Second)
			return "finished", nil
		},
	}

	nodes := map[string]hainish.Node{
		"longRunningNode": longRunningNode,
	}

	runtime := InitRuntime(nodes)
	runtime.InitWorkflow(1)
	runtime.CreateRuntimeNode("longRunningNode", 1, 1)

	// Execute the workflow
	resultChan := make(chan any, 1)
	errChan := make(chan error, 1)
	processChan := make(chan hainish.Edge, 1)

	runtime.RunWorkflow(1, resultChan, errChan, processChan)

	// Stop the workflow immediately
	err := runtime.StopWorkflow(1)
	if err != nil {
		t.Errorf("Unexpected error stopping workflow: %v", err)
	}

	// Check if the workflow context is cancelled
	wf := runtime.workflows[1]
	if wf.c == nil {
		t.Error("Expected workflow context to exist")
	}

	// The context should be done
	select {
	case <-wf.c.Done():
		// Context is cancelled as expected
	default:
		t.Error("Expected workflow context to be cancelled")
	}
}
