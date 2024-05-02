package api

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"servertest/internal/servertest"
)

var _ ServerInterface = (*Server)(nil)

type Server struct {
	storage    servertest.Storage
	controller servertest.Controller
}

func NewServer(
	storage servertest.Storage,
	controller servertest.Controller,
) *Server {
	return &Server{
		storage:    storage,
		controller: controller,
	}
}

func handleError(w http.ResponseWriter, err error) {
	w.WriteHeader(http.StatusInternalServerError)
	_, werr := w.Write([]byte(err.Error()))
	if werr != nil {
		handleWriteError(werr)
	}
}

func handleWriteError(err error) {
	slog.Error("unable to write response", slog.String("error", err.Error()))
}

func (s *Server) repository() servertest.Repository {
	return servertest.Repository{
		Storage: s.storage,
	}
}

func (s *Server) PostEntity(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req CreateEntityRequest
	slog.Info("kek")

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		handleError(w, err)
		return
	}

	id, err := s.controller.CreateEntity(ctx, s.repository(), servertest.CreateEntityOptions{
		Name:        req.Name,
		Description: req.Description,
		SomeValue:   req.SomeValue,
	})
	if err != nil {
		handleError(w, err)
		return
	}

	res := CreateEntityResponse{
		Id: id,
	}

	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(res)
	if err != nil {
		handleWriteError(err)
		return
	}
}

func (s *Server) GetKek(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(KekResponse{
		Id: "aaaa",
	})
}
