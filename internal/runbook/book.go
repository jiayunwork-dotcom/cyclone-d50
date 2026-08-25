package runbook

import (
	"errors"
	"sort"
	"strings"
	"sync"

	"cyclone-d50/internal/sep"
	"cyclone-d50/internal/spec"
)

const DefaultMax = 64

var (
	ErrNotFound   = errors.New("run not found")
	ErrExists     = errors.New("run id already exists")
	ErrTooMany    = errors.New("run book is full")
	ErrInvalidRun = errors.New("invalid run")
	ErrFrozen     = errors.New("run is frozen")
	ErrEmptyID    = errors.New("run id is empty")
)

type Entry struct {
	ID               string           `json:"id"`
	Spec             spec.Spec        `json:"spec"`
	D50M             float64          `json:"d50_m"`
	InletReynolds    float64          `json:"inlet_reynolds"`
	ParticleReynolds float64          `json:"particle_reynolds"`
	Grade            []sep.GradePoint `json:"grade"`
	TotalEfficiency  *float64         `json:"total_efficiency,omitempty"`
	HasPSD           bool             `json:"has_psd"`
	Warning          string           `json:"warning,omitempty"`
	Note             string           `json:"note"`
	Seq              uint64           `json:"seq"`
	Frozen           bool             `json:"frozen"`
}

type Book struct {
	mu    sync.RWMutex
	items map[string]Entry
	seq   uint64
	max   int
}

func NewBook(max int) *Book {
	if max <= 0 {
		max = DefaultMax
	}
	return &Book{
		items: make(map[string]Entry),
		max:   max,
	}
}

func (b *Book) Add(e Entry) error {
	if err := e.Validate(); err != nil {
		return err
	}
	b.mu.Lock()
	defer b.mu.Unlock()
	if _, ok := b.items[e.ID]; ok {
		return ErrExists
	}
	if len(b.items) >= b.max {
		return ErrTooMany
	}
	b.seq++
	e.Seq = b.seq
	b.items[e.ID] = e
	return nil
}

func (b *Book) Get(id string) (Entry, bool) {
	b.mu.RLock()
	defer b.mu.RUnlock()
	e, ok := b.items[id]
	return e, ok
}

func (b *Book) Remove(id string) bool {
	b.mu.Lock()
	defer b.mu.Unlock()
	if _, ok := b.items[id]; !ok {
		return false
	}
	delete(b.items, id)
	return true
}

func (b *Book) Len() int {
	b.mu.RLock()
	defer b.mu.RUnlock()
	return len(b.items)
}

func (b *Book) Max() int {
	b.mu.RLock()
	defer b.mu.RUnlock()
	return b.max
}

func (b *Book) List() []Entry {
	b.mu.RLock()
	defer b.mu.RUnlock()
	out := make([]Entry, 0, len(b.items))
	for _, e := range b.items {
		out = append(out, e)
	}
	sort.Slice(out, func(i, j int) bool {
		return out[i].Seq < out[j].Seq
	})
	return out
}

func (b *Book) Rename(oldID, newID string) error {
	b.mu.Lock()
	defer b.mu.Unlock()
	e, ok := b.items[oldID]
	if !ok {
		return ErrNotFound
	}
	if strings.TrimSpace(newID) == "" {
		return ErrEmptyID
	}
	if newID != oldID {
		if _, exists := b.items[newID]; exists {
			return ErrExists
		}
		delete(b.items, oldID)
		e.ID = newID
	}
	b.items[newID] = e
	return nil
}

func (b *Book) Freeze(id string) error {
	b.mu.Lock()
	defer b.mu.Unlock()
	e, ok := b.items[id]
	if !ok {
		return ErrNotFound
	}
	e.Frozen = true
	b.items[id] = e
	return nil
}

func (b *Book) SetNote(id, note string) error {
	b.mu.Lock()
	defer b.mu.Unlock()
	e, ok := b.items[id]
	if !ok {
		return ErrNotFound
	}
	if e.Frozen {
		return ErrFrozen
	}
	e.Note = note
	b.seq++
	e.Seq = b.seq
	b.items[id] = e
	return nil
}

func (b *Book) NextSeq() uint64 {
	b.mu.RLock()
	defer b.mu.RUnlock()
	return b.seq
}

func (e Entry) Validate() error {
	if strings.TrimSpace(e.ID) == "" {
		return ErrEmptyID
	}
	if e.Spec.CylinderDiameterM <= 0 || e.Spec.InletVelocityMPS <= 0 {
		return ErrInvalidRun
	}
	if e.D50M <= 0 || e.InletReynolds <= 0 {
		return ErrInvalidRun
	}
	return nil
}
