package util

import "errors"

var (
	ErrPluginAlreadyRegistered = errors.New("plugin already registered")
	ErrPluginNameEmpty         = errors.New("plugin name cannot be empty")
	ErrPluginNodesEmpty        = errors.New("plugin nodes cannot be empty")

	ErrNodeNameEmpty = errors.New("node name cannot be empty")
	ErrNodeActionNil = errors.New("node action cannot be nil")
)

var (
	ErrNodeNotFoundInPlugin   = errors.New("node not found in plugin")
	ErrWorkflowNotFound       = errors.New("workflow not found")
	ErrNodeNotFoundInWorkflow = errors.New("node not found")
	ErrPortNotFoundInNode     = errors.New("port not found in node")
	ErrDeletingNodeHasEdges   = errors.New("cannot delete node with existing edges")
	ErrPortNotExist           = errors.New("port not exist")
)

var (
	ErrPeerNotExist = errors.New("peer not exist")
)
