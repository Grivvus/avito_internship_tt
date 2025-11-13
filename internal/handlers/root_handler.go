package handlers

type RootHandler struct {
	*UserHandler
	*TeamHandler
	*PullRequestHandler
}

func NewRootHandler(
	uh *UserHandler,
	th *TeamHandler,
	prh *PullRequestHandler,
) *RootHandler {
	return &RootHandler{
		UserHandler:        uh,
		TeamHandler:        th,
		PullRequestHandler: prh,
	}
}
