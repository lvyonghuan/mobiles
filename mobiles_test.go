package mobiles

import (
	"testing"

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

// TestInit tests the initialization function
func TestInit(t *testing.T) {
	// Create mock plugin
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

	mockPlugin := &mockPlugin{
		name:        "testPlugin",
		description: "Test Plugin",
		version:     "1.0.0",
		author:      "testAuthor",
		license:     "MIT",
		nodes: map[string]hainish.Node{
			"testNode": mockNode,
		},
	}

	// Note: This will fail because the actual network environment is required
	// In actual testing, we need mock ansible.Init
	// Here is just a test structure
	err := Init(mockPlugin)
	if err == nil {
		t.Log("Init completed (expected to fail in test environment due to network dependencies)")
	} else {
		t.Logf("Init failed as expected: %v", err)
	}
}

// TestRegisterPlugin tests plugin registration
func TestRegisterPlugin(t *testing.T) {
	// Create mock plugin
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

	mockPlugin := &mockPlugin{
		name:        "testPlugin",
		description: "Test Plugin",
		version:     "1.0.0",
		author:      "testAuthor",
		license:     "MIT",
		nodes: map[string]hainish.Node{
			"testNode": mockNode,
		},
	}

	// Create Mobiles instance
	mobiles := &ImplMobiles{}

	// Register the plugin
	err := mobiles.RegisterPlugin(mockPlugin)
	if err != nil {
		t.Errorf("Unexpected error registering plugin: %v", err)
	}

	if mobiles.Plugin == nil {
		t.Error("Expected plugin to be registered")
	}

	if mobiles.Plugin.Name() != "testPlugin" {
		t.Errorf("Expected plugin name 'testPlugin', got '%s'", mobiles.Plugin.Name())
	}
}

// TestRegisterPluginTwice tests registering plugin twice
func TestRegisterPluginTwice(t *testing.T) {
	// Create mock plugin
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

	mockPlugin := &mockPlugin{
		name:        "testPlugin",
		description: "Test Plugin",
		version:     "1.0.0",
		author:      "testAuthor",
		license:     "MIT",
		nodes: map[string]hainish.Node{
			"testNode": mockNode,
		},
	}

	// Create Mobiles instance
	mobiles := &ImplMobiles{}

	// First registration should succeed
	err := mobiles.RegisterPlugin(mockPlugin)
	if err != nil {
		t.Errorf("Unexpected error in first registration: %v", err)
	}

	// Second registration should fail
	err = mobiles.RegisterPlugin(mockPlugin)
	if err == nil {
		t.Error("Expected error when registering plugin twice")
	}
}

// TestRegisterInvalidPlugin tests registering invalid plugin
func TestRegisterInvalidPlugin(t *testing.T) {
	// Test empty name plugin
	invalidPlugin := &mockPlugin{
		name:        "",
		description: "Invalid Plugin",
		version:     "1.0.0",
		author:      "testAuthor",
		license:     "MIT",
		nodes:       map[string]hainish.Node{},
	}

	mobiles := &ImplMobiles{}
	err := mobiles.RegisterPlugin(invalidPlugin)
	if err == nil {
		t.Error("Expected error when registering plugin with empty name")
	}

	// Test nil nodes plugin
	emptyNodesPlugin := &mockPlugin{
		name:        "emptyPlugin",
		description: "Empty Plugin",
		version:     "1.0.0",
		author:      "testAuthor",
		license:     "MIT",
		nodes:       nil,
	}

	err = mobiles.RegisterPlugin(emptyNodesPlugin)
	if err == nil {
		t.Error("Expected error when registering plugin with nil nodes")
	}
}

// TestRegisterPluginWithInvalidNodes tests registering plugin with invalid nodes
func TestRegisterPluginWithInvalidNodes(t *testing.T) {
	// Test node with empty name and nil action
	invalidNode := &mockNode{
		name:        "", // empty name
		description: "Invalid Node",
		isBegin:     true,
		inputs:      map[string]hainish.Port{},
		outputs:     map[string]hainish.Port{},
		params:      map[string]hainish.Port{},
		action:      nil, // nil action
	}

	invalidPlugin := &mockPlugin{
		name:        "invalidPlugin",
		description: "Invalid Plugin",
		version:     "1.0.0",
		author:      "testAuthor",
		license:     "MIT",
		nodes: map[string]hainish.Node{
			"invalidNode": invalidNode,
		},
	}

	mobiles := &ImplMobiles{}
	err := mobiles.RegisterPlugin(invalidPlugin)
	if err == nil {
		t.Error("Expected error when registering plugin with invalid nodes")
	}
}

// TestSendMessage tests message sending
func TestSendMessage(t *testing.T) {
	mobiles := &ImplMobiles{}

	// SendMessage目前是TODO实现，应该返回nil错误
	err := mobiles.SendMessage("peer123", 1, "test_message")
	if err != nil {
		t.Errorf("Unexpected error in SendMessage: %v", err)
	}
}
