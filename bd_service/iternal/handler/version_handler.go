package handler

import (
	"net/http"

	"dkl.ru/pact/bd_service/iternal/basedate"
	"dkl.ru/pact/bd_service/iternal/queue"
)

type VersionHandler struct {
	DB *basedate.Database
	QM *queue.QueueManager
}

func NewVersionHandler(db *basedate.Database, qm *queue.QueueManager) *VersionHandler {
	return &VersionHandler{
		DB: db,
		QM: qm,
	}
}

func (h *VersionHandler) GetVersions(w http.ResponseWriter, r *http.Request) {
	// Получаем список версий по языку из параметра запроса
}

func (h *VersionHandler) GetLastVersion(w http.ResponseWriter, r *http.Request) {

}
