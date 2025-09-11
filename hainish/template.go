package hainish

type Plugin interface {
	Name() string
	Description() string
	Version() string
	Author() string
	License() string

	Nodes() map[string]Node //The key is the node name
}

type Node interface {
	Name() string
	Description() string

	IsBegin() bool

	Inputs() map[string]Port  //The key is the port name.
	Outputs() map[string]Port //The key is the port name.
	Params() map[string]Port  //The key is the port name. Ansible use a maker to set.

	Action(inputs map[string]any, output map[string]chan any) (result any, err error)
}

type Port interface {
	Name() string
	Description() string
	Type() string

	Chan() chan any
}
