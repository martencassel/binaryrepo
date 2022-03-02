package group

import "github.com/gorilla/mux"

type GroupRouter struct {
}

func NewDockerRouter() *GroupRouter {
	return &GroupRouter{
	}
}

func (router *GroupRouter) RegisterHandlers(r *mux.Router) {
}