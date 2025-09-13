package mobiles

import (
	"github.com/lvyonghuan/Ubik-Util/uerr"
	"github.com/lvyonghuan/mobiles/ansible"
	"github.com/lvyonghuan/mobiles/hainish"
	"github.com/lvyonghuan/mobiles/runtime"
	"github.com/lvyonghuan/mobiles/util"
)

type Mobiles interface {
	RegisterPlugin(plugin hainish.Plugin) error // Initialize Mobiles
	SendMessage(peerID string, messageType int, message any) error
}

type ImplMobiles struct {
	Ansible *ansible.Ansible // Communicator for mobiles

	Plugin hainish.Plugin // A mobile for a planet. One mobile for one plugin.
}

func Init(plugin hainish.Plugin) error {
	mobiles := new(ImplMobiles)

	// Register the plugin
	err := mobiles.RegisterPlugin(plugin)
	if err != nil {
		return err
	}

	// Initialize runtime
	r := runtime.InitRuntime(plugin.Nodes())

	// Initialize Ansible for communication
	asb, err := ansible.Init(plugin, r)
	if err != nil {
		return err
	}
	mobiles.Ansible = asb

	return nil
}

func (m *ImplMobiles) RegisterPlugin(plugin hainish.Plugin) error {
	if m.Plugin != nil {
		return uerr.NewError(util.ErrPluginAlreadyRegistered)
	}

	// Check plugin validity
	if plugin.Name() == "" {
		return uerr.NewError(util.ErrPluginNameEmpty)
	}
	if plugin.Nodes() == nil || len(plugin.Nodes()) == 0 {
		return uerr.NewError(util.ErrPluginNodesEmpty)
	}
	err := checkNodesValidity(plugin.Nodes())
	if err != nil {
		return err
	}

	// Register the plugin
	m.Plugin = plugin
	return nil
}

func checkNodesValidity(nodes map[string]hainish.Node) error {
	for _, node := range nodes {
		if node.Name() == "" {
			return uerr.NewError(util.ErrNodeNameEmpty)
		}
	}

	return nil
}

func (m *ImplMobiles) SendMessage(peerID string, messageType int, message any) error {
	// TODO
	return nil
}
