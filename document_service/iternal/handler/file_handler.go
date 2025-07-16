package handler

import "net/http"

type FileHandler struct{}

func NewFileHandler() *FileHandler {
	return &FileHandler{}
}

func (h *FileHandler) SaveFileInBd(w http.ResponseWriter, r *http.Request) {}
