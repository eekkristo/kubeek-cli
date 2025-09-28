package state

import (
	"encoding/json"
	"os"
)

type Replacement struct {
	Placeholder string `json:"placeholder"`
	Value       string `json:"value"`
	Count       int    `json:"count,omitempty"`
}

type Entry struct {
	Line         int           `json:"line"`
	Before       string        `json:"before"`
	After        string        `json:"after"`
	Replacements []Replacement `json:"replacements"`
}

type State map[string][]Entry // key = file path (we use RELATIVE for generated trees)

func Load(path string) (State, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return State{}, nil
		}
		return nil, err
	}
	var st State
	if err := json.Unmarshal(b, &st); err != nil {
		return nil, err
	}
	return st, nil
}

func Save(path string, st State) error {
	b, err := json.MarshalIndent(st, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, b, 0644)
}

func MergeEntries(existing, additions []Entry) []Entry {
	byLine := map[int]Entry{}
	for _, e := range existing {
		byLine[e.Line] = e
	}
	for _, a := range additions {
		byLine[a.Line] = a
	}
	out := make([]Entry, 0, len(byLine))
	for _, v := range byLine {
		out = append(out, v)
	}
	// stable sort by line
	for i := 0; i < len(out)-1; i++ {
		for j := i + 1; j < len(out); j++ {
			if out[j].Line < out[i].Line {
				out[i], out[j] = out[j], out[i]
			}
		}
	}
	return out
}
