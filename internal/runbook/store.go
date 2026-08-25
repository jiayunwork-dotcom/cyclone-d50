package runbook

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

const SnapshotVersion = 1

type Snapshot struct {
	Version int     `json:"version"`
	Max     int     `json:"max"`
	Seq     uint64  `json:"seq"`
	Entries []Entry `json:"entries"`
}

func (b *Book) Export() Snapshot {
	b.mu.RLock()
	defer b.mu.RUnlock()
	out := Snapshot{
		Version: SnapshotVersion,
		Max:     b.max,
		Seq:     b.seq,
		Entries: make([]Entry, 0, len(b.items)),
	}
	for _, e := range b.items {
		out.Entries = append(out.Entries, e)
	}
	return out
}

func (b *Book) Import(snapshot Snapshot) error {
	if snapshot.Version != SnapshotVersion {
		return fmt.Errorf("unsupported snapshot version %d", snapshot.Version)
	}
	if snapshot.Max <= 0 {
		return fmt.Errorf("snapshot max must be positive")
	}
	if len(snapshot.Entries) > snapshot.Max {
		return ErrTooMany
	}
	next := make(map[string]Entry, len(snapshot.Entries))
	for _, e := range snapshot.Entries {
		if err := e.Validate(); err != nil {
			return err
		}
		if _, ok := next[e.ID]; ok {
			return ErrExists
		}
		next[e.ID] = e
	}
	b.mu.Lock()
	defer b.mu.Unlock()
	b.items = next
	b.max = snapshot.Max
	b.seq = snapshot.Seq
	return nil
}

func (b *Book) SaveFile(path string) error {
	data, err := json.MarshalIndent(b.Export(), "", "  ")
	if err != nil {
		return err
	}
	dir := filepath.Dir(path)
	tmp, err := os.CreateTemp(dir, ".cyclone-book-*")
	if err != nil {
		return err
	}
	tmpName := tmp.Name()
	if _, err := tmp.Write(data); err != nil {
		_ = tmp.Close()
		_ = os.Remove(tmpName)
		return err
	}
	if err := tmp.Sync(); err != nil {
		_ = tmp.Close()
		_ = os.Remove(tmpName)
		return err
	}
	if err := tmp.Close(); err != nil {
		_ = os.Remove(tmpName)
		return err
	}
	if err := os.Rename(tmpName, path); err != nil {
		_ = os.Remove(tmpName)
		return err
	}
	return nil
}

func LoadFile(path string) (*Book, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var snapshot Snapshot
	if err := json.Unmarshal(data, &snapshot); err != nil {
		return nil, err
	}
	b := NewBook(snapshot.Max)
	if err := b.Import(snapshot); err != nil {
		return nil, err
	}
	return b, nil
}
