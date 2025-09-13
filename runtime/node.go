package runtime

import (
	"context"

	"github.com/lvyonghuan/Ubik-Util/uerr"
	"github.com/lvyonghuan/mobiles/hainish"
	"github.com/lvyonghuan/mobiles/util"
)

type runtimeNode struct {
	node        *hainish.Node
	outputEdges map[int]edge
	params      map[string]any
}

func (rn *runtimeNode) startSendParams(stopContext context.Context) error {
	paramPorts := (*rn.node).Params()

	for portName, param := range rn.params {
		port, isExist := paramPorts[portName]
		if !isExist {
			return uerr.NewError(util.ErrPortNotExist)
		}

		go sendParam(param, port.Chan(), stopContext)
	}

	return nil
}

func sendParam(param any, port chan any, stopContext context.Context) {
	for {
		select {
		case port <- param: // Send param to port
		case <-stopContext.Done(): // Stop signal received
			return
		}
	}
}
