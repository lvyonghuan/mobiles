package hainish

import (
	"testing"
)

// TestNewPort tests port creation
func TestNewPort(t *testing.T) {
	port := NewPort("input1", "Input port 1", "string")

	if port.Name() != "input1" {
		t.Errorf("Expected port name 'input1', got '%s'", port.Name())
	}

	if port.Description() != "Input port 1" {
		t.Errorf("Expected port description 'Input port 1', got '%s'", port.Description())
	}

	if port.Type() != "string" {
		t.Errorf("Expected port type 'string', got '%s'", port.Type())
	}

	if port.Chan() == nil {
		t.Error("Expected port channel to be initialized")
	}
}

// TestNewNode tests node creation
func TestNewNode(t *testing.T) {
	// Create input ports
	inputPorts := map[string]Port{
		"input1": NewPort("input1", "Input 1", "string"),
	}

	// Create output ports
	outputPorts := map[string]Port{
		"output1": NewPort("output1", "Output 1", "string"),
	}

	// Create param ports
	paramPorts := map[string]Port{
		"param1": NewPort("param1", "Parameter 1", "int"),
	}

	// Simulate action function
	action := func(inputs map[string]any, output map[string]chan any) (result any, err error) {
		// Simple action: echo input1 to output1
		if val, exists := inputs["input1"]; exists {
			if output["output1"] != nil {
				output["output1"] <- val
			}
			return val, nil
		}
		return nil, nil
	}

	node := NewNode("testNode", "Test Node", true, inputPorts, outputPorts, paramPorts, action)

	// Test node basic properties
	if node.Name() != "testNode" {
		t.Errorf("Expected node name 'testNode', got '%s'", node.Name())
	}

	if node.Description() != "Test Node" {
		t.Errorf("Expected node description 'Test Node', got '%s'", node.Description())
	}

	if !node.IsBegin() {
		t.Error("Expected node to be a begin node")
	}

	// Test ports
	if len(node.Inputs()) != 1 {
		t.Errorf("Expected 1 input port, got %d", len(node.Inputs()))
	}

	if len(node.Outputs()) != 1 {
		t.Errorf("Expected 1 output port, got %d", len(node.Outputs()))
	}

	if len(node.Params()) != 1 {
		t.Errorf("Expected 1 param port, got %d", len(node.Params()))
	}

	// Test action execution
	inputs := map[string]any{"input1": "test_value"}
	outputs := make(map[string]chan any)
	outputs["output1"] = make(chan any, 1)

	result, err := node.Action(inputs, outputs)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if result != "test_value" {
		t.Errorf("Expected result 'test_value', got '%v'", result)
	}

	// Check output channel
	select {
	case outputVal := <-outputs["output1"]:
		if outputVal != "test_value" {
			t.Errorf("Expected output 'test_value', got '%v'", outputVal)
		}
	default:
		t.Error("Expected value in output channel")
	}
}

// TestNewPlugin tests plugin creation
func TestNewPlugin(t *testing.T) {
	// Create a simple node to include in the plugin
	inputPorts := map[string]Port{
		"input1": NewPort("input1", "Input 1", "string"),
	}
	outputPorts := map[string]Port{
		"output1": NewPort("output1", "Output 1", "string"),
	}
	action := func(inputs map[string]any, output map[string]chan any) (result any, err error) {
		return "processed", nil
	}

	node := NewNode("testNode", "Test Node", true, inputPorts, outputPorts, nil, action)
	nodes := map[string]Node{
		"testNode": node,
	}

	plugin := NewPlugin("testPlugin", "Test Plugin", "1.0.0", "testAuthor", "MIT", nodes)

	// Test plugin basic properties
	if plugin.Name() != "testPlugin" {
		t.Errorf("Expected plugin name 'testPlugin', got '%s'", plugin.Name())
	}

	if plugin.Description() != "Test Plugin" {
		t.Errorf("Expected plugin description 'Test Plugin', got '%s'", plugin.Description())
	}

	if plugin.Version() != "1.0.0" {
		t.Errorf("Expected plugin version '1.0.0', got '%s'", plugin.Version())
	}

	if plugin.Author() != "testAuthor" {
		t.Errorf("Expected plugin author 'testAuthor', got '%s'", plugin.Author())
	}

	if plugin.License() != "MIT" {
		t.Errorf("Expected plugin license 'MIT', got '%s'", plugin.License())
	}

	// Test nodes
	if len(plugin.Nodes()) != 1 {
		t.Errorf("Expected 1 node, got %d", len(plugin.Nodes()))
	}

	if _, exists := plugin.Nodes()["testNode"]; !exists {
		t.Error("Expected testNode to exist in plugin nodes")
	}
}

// TestEdgeCreation tests edge creation
func TestEdgeCreation(t *testing.T) {
	edge := NewEdge("peer123", 1, 2, "input1")

	if edge.TargetWorkflowID != 1 {
		t.Errorf("Expected TargetWorkflowID 1, got %d", edge.TargetWorkflowID)
	}

	if edge.TargetNodeID != 2 {
		t.Errorf("Expected TargetNodeID 2, got %d", edge.TargetNodeID)
	}

	if edge.TargetPort != "input1" {
		t.Errorf("Expected TargetPort 'input1', got '%s'", edge.TargetPort)
	}

	// Test setting and getting Value
	edge.Value = "test_data"
	if edge.Value != "test_data" {
		t.Errorf("Expected Value 'test_data', got '%v'", edge.Value)
	}
}
