package hainish

type ImplPlugin struct {
	PluginName        string `json:"name"`
	PluginDescription string `json:"description"`
	PluginVersion     string `json:"version"`
	PluginAuthor      string `json:"author"`
	PluginLicense     string `json:"license"`

	NodeMap map[string]Node //The key is the node name
}

func NewPlugin(name, description, version, author, license string, nodes map[string]Node) ImplPlugin {
	return ImplPlugin{
		PluginName:        name,
		PluginDescription: description,
		PluginVersion:     version,
		PluginAuthor:      author,
		PluginLicense:     license,
		NodeMap:           nodes,
	}
}

func (i ImplPlugin) Name() string {
	return i.PluginName
}

func (i ImplPlugin) Description() string {
	return i.PluginDescription
}

func (i ImplPlugin) Version() string {
	return i.PluginVersion
}

func (i ImplPlugin) Author() string {
	return i.PluginAuthor
}

func (i ImplPlugin) License() string {
	return i.PluginLicense
}

func (i ImplPlugin) Nodes() map[string]Node {
	return i.NodeMap
}

type ImplNode struct {
	NodeName        string `json:"name"`
	NodeDescription string `json:"description"`
	IsBeginNode     bool   `json:"is_begin"`

	InputMap   map[string]Port `json:"input"`  //The key is the port name.
	OutputMap  map[string]Port `json:"output"` //The key is the port name.
	ParamMap   map[string]Port `json:"param"`  //The key is the port name. Ansible use a maker to set.
	NodeAction func(inputs map[string]any, output map[string]chan any) (result any, err error)
}

func NewNode(name, description string, isBegin bool, inputs, outputs, params map[string]Port, action func(inputs map[string]any, output map[string]chan any) (result any, err error)) ImplNode {
	return ImplNode{
		NodeName:        name,
		NodeDescription: description,
		IsBeginNode:     isBegin,
		InputMap:        inputs,
		OutputMap:       outputs,
		ParamMap:        params,
		NodeAction:      action,
	}
}

func (i ImplNode) Name() string {
	return i.NodeName
}

func (i ImplNode) Description() string {
	return i.NodeDescription
}

func (i ImplNode) IsBegin() bool {
	return i.IsBeginNode
}

func (i ImplNode) Inputs() map[string]Port {
	return i.InputMap
}

func (i ImplNode) Outputs() map[string]Port {
	return i.OutputMap
}

func (i ImplNode) Params() map[string]Port {
	return i.ParamMap
}

func (i ImplNode) Action(inputs map[string]any, output map[string]chan any) (result any, err error) {
	return i.NodeAction(inputs, output)
}

type ImplPort struct {
	PortName        string `json:"name"`
	PortDescription string `json:"description"`
	PortType        string `json:"type"`
	PortChan        chan any
}

func NewPort(name, description, portType string) ImplPort {
	return ImplPort{
		PortName:        name,
		PortDescription: description,
		PortType:        portType,
		PortChan:        make(chan any, 1),
	}
}

func (i ImplPort) Name() string {
	return i.PortName
}

func (i ImplPort) Description() string {
	return i.PortDescription
}

func (i ImplPort) Type() string {
	return i.PortType
}

func (i ImplPort) Chan() chan any {
	return i.PortChan
}
