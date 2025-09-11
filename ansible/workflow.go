package ansible

import (
	"context"

	"github.com/lvyonghuan/Ubik-Util/ulog"
	"github.com/lvyonghuan/mobiles/hainish"
)

// Handel workflow

type workflowListener struct {
	workflowID int

	resultChan  chan any
	errChan     chan error
	processChan chan hainish.Edge

	stopContext context.Context

	ansible Ansible
}

func (asb *ImplAnsible) initWorkflowListener(workflowID int) (chan any, chan error, chan hainish.Edge) {
	var wl workflowListener
	wl.workflowID = workflowID
	wl.ansible = asb

	wl.resultChan = make(chan any, 1)
	wl.errChan = make(chan error, 1)
	wl.processChan = make(chan hainish.Edge, 1)

	wl.ansible = asb

	return wl.resultChan, wl.errChan, wl.processChan
}

func (asb *ImplAnsible) getWorkflowListener(workflowID int) *workflowListener {
	wl, exist := asb.wfListener[workflowID]
	if !exist {
		return nil
	}
	return &wl
}

func (workflowListener *workflowListener) setStopContext(c context.Context) {
	workflowListener.stopContext = c
}

func (workflowListener *workflowListener) run() {
	// Listen result and error
	for {
		select {
		case processData := <-workflowListener.processChan:
			err := workflowListener.ansible.getPeerManager().sendProcessDataToFollower(processData)
			if err != nil {
				// TODO: 处理错误
			}
		case err := <-workflowListener.errChan:
			er := workflowListener.ansible.getPeerManager().sendLogToLeader(ulog.Error, err.Error())
			if er != nil {
				// TODO: 处理这个错误
			}
		case result := <-workflowListener.resultChan:
			err := workflowListener.ansible.getPeerManager().sendResultToLeader(result)
			if err != nil {
				// TODO: 处理错误
			}
		case <-workflowListener.stopContext.Done():
			return
		}
	}
}
