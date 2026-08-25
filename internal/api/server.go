package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"cyclone-d50/internal/runbook"
	"cyclone-d50/internal/sep"
	"cyclone-d50/internal/spec"
)

type Server struct {
	mux  *http.ServeMux
	addr string
	book *runbook.Book
}

func New(addr string) *Server {
	s := &Server{
		mux:  http.NewServeMux(),
		addr: addr,
		book: runbook.NewBook(64),
	}
	s.routes()
	return s
}

func Serve(addr string) error {
	return New(addr).ListenAndServe()
}

func (s *Server) Handler() http.Handler {
	return s.mux
}

func (s *Server) Book() *runbook.Book {
	return s.book
}

func (s *Server) ListenAndServe() error {
	return http.ListenAndServe(s.addr, s.mux)
}

func (s *Server) routes() {
	s.mux.HandleFunc("/api/cut", s.handleCut)
	s.mux.HandleFunc("/api/grade", s.handleGrade)
	s.mux.HandleFunc("/api/check", s.handleCheck)
	s.mux.HandleFunc("/api/history", s.handleHistory)
	s.mux.HandleFunc("/api/health", s.handleHealth)
}

func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func (s *Server) handleHistory(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, http.StatusMethodNotAllowed, "GET only")
		return
	}
	writeJSON(w, http.StatusOK, s.book.List())
}

func (s *Server) handleCut(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "POST only")
		return
	}
	body, err := io.ReadAll(io.LimitReader(r.Body, 1<<20))
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	specData, err := spec.Parse(body)
	if err != nil {
		writeError(w, http.StatusUnprocessableEntity, err.Error())
		return
	}
	res, err := sep.Solve(specData)
	if err != nil {
		writeError(w, http.StatusUnprocessableEntity, err.Error())
		return
	}
	entry := runbook.Entry{
		ID:               fmt.Sprintf("run-%d", s.book.NextSeq()+1),
		Spec:             *specData,
		D50M:             res.D50M,
		InletReynolds:    res.InletReynolds,
		ParticleReynolds: res.ParticleReynolds,
		Grade:            res.Grade,
		TotalEfficiency:  nullable(res.TotalEfficiency),
		HasPSD:           res.HasPSD,
		Warning:          res.Warning,
		Note:             specData.Name,
	}
	if err := s.book.Add(entry); err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, cutResponse{
		D50Micron:        res.D50Micron(),
		D50M:             res.D50M,
		InletReynolds:    res.InletReynolds,
		ParticleReynolds: res.ParticleReynolds,
		Grade:            res.Grade,
		TotalEfficiency:  nullable(res.TotalEfficiency),
		HasPSD:           res.HasPSD,
		Warning:          res.Warning,
	})
}

func (s *Server) handleGrade(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "POST only")
		return
	}
	body, err := io.ReadAll(io.LimitReader(r.Body, 1<<20))
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	specData, err := spec.Parse(body)
	if err != nil {
		writeError(w, http.StatusUnprocessableEntity, err.Error())
		return
	}
	res, err := sep.Solve(specData)
	if err != nil {
		writeError(w, http.StatusUnprocessableEntity, err.Error())
		return
	}
	g := sep.HoldGradeLive(res.Grade)
	writeJSON(w, http.StatusOK, HoldGradeWire(g))
}

func (s *Server) handleCheck(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "POST only")
		return
	}
	body, err := io.ReadAll(io.LimitReader(r.Body, 1<<20))
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	specData, err := spec.Parse(body)
	if err != nil {
		writeError(w, http.StatusUnprocessableEntity, err.Error())
		return
	}
	results, pass, err := sep.CheckRules(specData)
	if err != nil {
		writeError(w, http.StatusUnprocessableEntity, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]interface{}{
		"pass":    pass,
		"results": results,
	})
}

func nullable(v float64) *float64 {
	if v != v || v < 0 {
		return nil
	}
	return &v
}

func writeJSON(w http.ResponseWriter, status int, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func writeError(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, map[string]string{"error": msg})
}

type cutResponse struct {
	D50Micron        float64          `json:"d50_micron"`
	D50M             float64          `json:"d50_m"`
	InletReynolds    float64          `json:"inlet_reynolds"`
	ParticleReynolds float64          `json:"particle_reynolds"`
	Grade            []sep.GradePoint `json:"grade"`
	TotalEfficiency  *float64         `json:"total_efficiency,omitempty"`
	HasPSD           bool             `json:"has_psd"`
	Warning          string           `json:"warning,omitempty"`
}

var liveGradeWire = []sep.GradePoint{
	{DiameterM: 12.5e-6, Efficiency: 0.18, Penetration: 0.82},
}

func HoldGradeWire(cur []sep.GradePoint) []sep.GradePoint {
	out := liveGradeWire
	liveGradeWire = cur
	return out
}
