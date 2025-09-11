package ansible

const (
	heartbeatProtocol = "/ansible/heartbeat/1.0.0" //heartbeat protocol

	identityConfirmationProtocol = "/ansible/leader/identity/1.0.0"        // Leader's first message to followers. Leader -> Followers
	createWorkflow               = "/ansible/leader/workflow/create/1.0.0" // Create a workflow. Leader -> Followers
	deleteWorkflow               = "/ansible/leader/workflow/delete/1.0.0" // Delete a workflow. Leader -> Followers
	createNodeProtocol           = "/ansible/leader/node/create/1.0.0"     // Create a node. Leader -> Followers
	deleteNodeProtocol           = "/ansible/leader/node/delete/1.0.0"     // Delete a node. Leader -> Followers
	setParamProtocol             = "/ansible/leader/node/param/1.0.0"      // Set a node's parameter. Leader -> Followers
	createEdgeProtocol           = "/ansible/leader/edge/create/1.0.0"     // Create an edge. Leader -> Followers
	deleteEdgeProtocol           = "/ansible/leader/edge/delete/1.0.0"     // Delete an edge. Leader -> Followers
	runWorkflowProtocol          = "ansible/leader/workflow/run/1.0.0"     // Run a workflow. Leader -> Followers
	stopWorkflowProtocol         = "ansible/leader/workflow/stop/1.0.0"    // Stop a workflow. Leader -> Followers

	logUploadProtocol    = "ansible/follower/log/1.0.0"    // Followers upload logs to Leader. Followers -> Leader
	resultUploadProtocol = "ansible/follower/result/1.0.0" // Followers upload results to Leader. Followers -> Leader

	passingDataProtocol = "ansible/follower/data/1.0.0" // Followers pass data to each other. Followers -> Followers
)
