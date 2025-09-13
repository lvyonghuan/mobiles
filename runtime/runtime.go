package runtime

import (
	"context"

	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/lvyonghuan/Ubik-Util/uerr"
	"github.com/lvyonghuan/mobiles/hainish"
	"github.com/lvyonghuan/mobiles/util"
)

type Runtime struct {
	workflows map[int]*workflow

	nodes map[string]hainish.Node
}

func InitRuntime(nodes map[string]hainish.Node) *Runtime {
	return &Runtime{
		workflows: make(map[int]*workflow),
		nodes:     nodes,
	}
}

type workflow struct {
	runtimeNodes map[int]*runtimeNode
	c            context.Context
	cancel       context.CancelFunc

	edges map[int]edge

	resultChan  chan any
	errChan     chan error
	processChan chan hainish.Edge
}

type edge struct {
	e              hainish.Edge
	fromPort       chan any
	producerNodeID int
}

func (r *Runtime) InitWorkflow(workflowID int) {
	// Initialize a workflow
	// Each workflow has its own context
	c, cancel := context.WithCancel(context.Background())
	runtimeNodes := make(map[int]*runtimeNode)
	r.workflows[workflowID] = &workflow{
		runtimeNodes: runtimeNodes,
		c:            c,
		cancel:       cancel,
		edges:        make(map[int]edge),
	}
}

func (r *Runtime) CreateRuntimeNode(nodeName string, nodeID int, workflowID int) error {
	node, isExist := r.nodes[nodeName]
	if !isExist {
		return uerr.NewError(util.ErrNodeNotFoundInPlugin)
	}
	wf, exist := r.workflows[workflowID]
	if !exist {
		return uerr.NewError(util.ErrWorkflowNotFound)
	}

	// Create a runtime node
	// TODO 这里应该有一个警告判断，当ID已经存在时
	wf.runtimeNodes[nodeID] = &runtimeNode{
		node:        &node,
		outputEdges: make(map[int]edge),
		params:      make(map[string]any),
	}
	return nil
}

func (r *Runtime) CreateEdge(edgeID int, destination peer.ID, workflowID int, producerNodeID int, producerPortName string, consumerNodeID int, consumerPortName string) error {
	wf, exist := r.workflows[workflowID]
	if !exist {
		return uerr.NewError(util.ErrWorkflowNotFound)
	}

	producerNode, exist := wf.runtimeNodes[producerNodeID]
	if !exist {
		return uerr.NewError(util.ErrNodeNotFoundInWorkflow)
	}

	producerPort, exist := (*producerNode.node).Outputs()[producerPortName]
	if !exist {
		return uerr.NewError(util.ErrPortNotFoundInNode)
	}

	e := hainish.NewEdge(destination, workflowID, consumerNodeID, consumerPortName)
	fromPort := producerPort.Chan()
	// Add the edge to the workflow
	wf.edges[edgeID] = edge{
		e:              *e,
		fromPort:       fromPort,
		producerNodeID: producerNodeID,
	}

	// Add the edge to the producer node's output edges
	producerNode.outputEdges[edgeID] = wf.edges[edgeID]

	return nil
}

func (r *Runtime) RunWorkflow(workflowID int, resultChan chan any, errChan chan error, processChan chan hainish.Edge) (context.Context, error) {
	wf, exist := r.workflows[workflowID]
	if !exist {
		return nil, uerr.NewError(util.ErrWorkflowNotFound)
	}

	// Set the result and error channels
	wf.resultChan = resultChan
	wf.errChan = errChan
	wf.processChan = processChan

	// Start all nodes in the workflow
	for _, runtimeNode := range wf.runtimeNodes {
		// Start sending params if any
		err := runtimeNode.startSendParams(wf.c)
		if err != nil {
			return nil, err
		}

		// Listen for results and errors
		go wf.listenResultAndError()

		// Run each node in a separate goroutine
		go wf.runNode(*runtimeNode.node)
	}

	return wf.c, nil
}

func (r *Runtime) StopWorkflow(workflowID int) error {
	wf, exist := r.workflows[workflowID]
	if !exist {
		return uerr.NewError(util.ErrWorkflowNotFound)
	}

	// Stop the workflow by cancelling its context
	wf.cancel()
	return nil
}

func (r *Runtime) DeleteWorkflow(workflowID int) error {
	wf, exist := r.workflows[workflowID]
	if !exist {
		return uerr.NewError(util.ErrWorkflowNotFound)
	}

	// FIXME 其他状态检测
	wf.cancel() // Ensure the workflow is stopped

	delete(r.workflows, workflowID)
	return nil
}

func (r *Runtime) DeleteNode(workflowID, nodeID int) error {
	wf, exist := r.workflows[workflowID]
	if !exist {
		return uerr.NewError(util.ErrWorkflowNotFound)
	}

	wf.cancel() // Ensure the workflow is stopped

	node, exist := wf.runtimeNodes[nodeID]
	if !exist {
		return uerr.NewError(util.ErrNodeNotFoundInWorkflow)
	}

	//Ensure this node don't have any edge
	if len(node.outputEdges) > 0 {
		return uerr.NewError(util.ErrDeletingNodeHasEdges)
	}

	// Delete the node
	delete(wf.runtimeNodes, nodeID)
	return nil
}

func (r *Runtime) DeleteEdge(workflowID, edgeID int) error {
	wf, exist := r.workflows[workflowID]
	if !exist {
		return uerr.NewError(util.ErrWorkflowNotFound)
	}

	edge, exist := wf.edges[edgeID]
	if !exist {
		return uerr.NewError(util.ErrPortNotFoundInNode)
	}

	// Delete the edge from the producer node's output edges
	producerNode, exist := wf.runtimeNodes[edge.producerNodeID]
	if !exist {
		return uerr.NewError(util.ErrNodeNotFoundInWorkflow)
	}
	delete(producerNode.outputEdges, edgeID)

	// Delete the edge from the workflow
	delete(wf.edges, edgeID)
	return nil
}

func (r *Runtime) SetParam(workflowID, nodeID int, portName string, value any) error {
	wf, exist := r.workflows[workflowID]
	if !exist {
		return uerr.NewError(util.ErrWorkflowNotFound)
	}

	node, exist := wf.runtimeNodes[nodeID]
	if !exist {
		return uerr.NewError(util.ErrNodeNotFoundInWorkflow)
	}

	node.params[portName] = value
	return nil
}
