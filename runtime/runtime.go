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
	runtimeNodes map[int]*hainish.Node
	c            context.Context
	cancel       context.CancelFunc

	edges []edge

	resultChan  chan any
	errChan     chan error
	processChan chan hainish.Edge
}

type edge struct {
	e        hainish.Edge
	fromPort chan any
}

func (r *Runtime) InitWorkflow(workflowID int) {
	// Initialize a workflow
	// Each workflow has its own context
	c, cancel := context.WithCancel(context.Background())
	runtimeNodes := make(map[int]*hainish.Node)
	r.workflows[workflowID] = &workflow{
		runtimeNodes: runtimeNodes,
		c:            c,
		cancel:       cancel,
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
	wf.runtimeNodes[nodeID] = &node
	return nil
}

func (r *Runtime) CreateEdge(destination peer.ID, workflowID int, producerNodeID int, producerPortName string, consumerNodeID int, consumerPortName string) error {
	wf, exist := r.workflows[workflowID]
	if !exist {
		return uerr.NewError(util.ErrWorkflowNotFound)
	}

	producerNode, exist := wf.runtimeNodes[producerNodeID]
	if !exist {
		return uerr.NewError(util.ErrNodeNotFoundInWorkflow)
	}

	producerPort, exist := (*producerNode).Outputs()[producerPortName]
	if !exist {
		return uerr.NewError(util.ErrPortNotFoundInNode)
	}

	e := hainish.NewEdge(destination, workflowID, consumerNodeID, consumerPortName)
	fromPort := producerPort.Chan()
	// Add the edge to the workflow
	wf.edges = append(wf.edges, edge{
		e:        *e,
		fromPort: fromPort,
	})

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
		// Listen for results and errors
		go wf.listenResultAndError()

		// Run each node in a separate goroutine
		go wf.runNode(*runtimeNode)
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
