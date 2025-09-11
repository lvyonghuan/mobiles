package runtime

import (
	"context"

	"github.com/lvyonghuan/mobiles/hainish"
)

func (w *workflow) runNode(node hainish.Node) {
	params := node.Params()
	inputs := node.Inputs()
	outputs := node.Outputs()
	out := make(map[string]chan any)

	// Prepare output channels
	for outputName, port := range outputs {
		out[outputName] = port.Chan()
	}

	// Loop to get inputs and params, then execute the node
	for i := 0; ; i++ {
		// Get inputs value
		in := make(map[string]any)
		for paramName, port := range params {
			select {
			case in[paramName] = <-port.Chan():
			case <-w.c.Done():
				return
			}
		}

		// The beginning node will skip the first epoch's input (if it has any)
		// to start this workflow.
		// Otherwise, the workflow will be blocked forever (if beginning node has any
		// input, the node will wait for it).
		if i != 0 || !node.IsBegin() {
			for inputName, port := range inputs {
				select {
				case in[inputName] = <-port.Chan():
				case <-w.c.Done():
					return
				}
			}
		}

		// Execute the node
		result, err := node.Action(in, out)
		if err != nil {
			w.errChan <- err
		}
		if result != nil {
			w.resultChan <- result
		}
	}
}

func (w *workflow) listenResultAndError() {
	// Listen edges
	for _, e := range w.edges {
		go e.listenEdgeOutput(w.processChan, w.c) // Listen edge output
	}
}

func (e *edge) listenEdgeOutput(processChan chan hainish.Edge, cancelContext context.Context) {
	for {
		select {
		case result := <-e.fromPort:
			edge := e.e
			edge.Value = result
			processChan <- edge
		case <-cancelContext.Done():
			return
		}
	}
}
