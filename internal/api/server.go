package api

import (
	"encoding/json"
	"errors"
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
	if errors.Is(err, &servertest.ErrNotFound{}) {
		w.WriteHeader(http.StatusNotFound)
	} else {
		w.WriteHeader(http.StatusInternalServerError)
	}
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

func (s *Server) GetLocationRiderId(w http.ResponseWriter, r *http.Request, riderId string, params GetLocationRiderIdParams) {
	ctx := r.Context()

	history, err := s.controller.GetRiderLocationHistory(ctx, s.repository(), servertest.GetRiderLocationHistoryOptions{
		RiderID: riderId,
		Limit:   params.Max,
	})
	if err != nil {
		handleError(w, err)
		return
	}

	res := GetLocationHistoryResponse{
		RiderId: riderId,
		History: make([]LocationEntry, 0, len(history)),
	}

	for _, e := range history {
		res.History = append(res.History, LocationEntry{
			Long: e.Long,
			Lat:  e.Lat,
		})
	}

	w.WriteHeader(http.StatusOK)

	err = json.NewEncoder(w).Encode(res)
	if err != nil {
		handleWriteError(err)
	}
}

func (s *Server) PostLocationRiderIdNow(w http.ResponseWriter, r *http.Request, riderId string) {
	ctx := r.Context()

	var req LocationEntry
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		handleError(w, err)
		return
	}

	err = s.controller.AddRiderLocation(ctx, s.repository(), servertest.LocationEntry{
		RiderID: riderId,
		Long:    req.Long,
		Lat:     req.Lat,
	})
	if err != nil {
		handleError(w, err)
		return
	}

	w.WriteHeader(http.StatusOK)
	return
}
