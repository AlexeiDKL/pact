package handler

import "net/http"

type SchedulerHandler struct{}

func NewSchedulerHandler() *SchedulerHandler {
	return &SchedulerHandler{}
}

func (h *SchedulerHandler) Status(w http.ResponseWriter, r *http.Request) {

}
