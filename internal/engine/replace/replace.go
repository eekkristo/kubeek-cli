package replace

import (
	"sort"
	"strings"

	cfg "kubeekcli/internal/config"
	st "kubeekcli/internal/engine/state"
)

// Pure text replacement: placeholders -> values
func RenderContentByConfig(content string, c cfg.Config) (string, bool) {
	out := content
	changed := false
	for ph, val := range c {
		if ph == "" {
			continue
		}
		if strings.Contains(out, ph) {
			out = strings.ReplaceAll(out, ph, val)
			changed = true
		}
	}
	return out, changed
}

// values -> placeholders (fallback reverse)
func ReverseContentByConfig(content string, c cfg.Config) (string, cfg.Config, bool) {
	type kv struct {
		ph, val string
	}
	pairs := make([]kv, 0, len(c))
	for ph, val := range c {
		pairs = append(pairs, kv{ph: ph, val: val})
	}
	sort.Slice(pairs, func(i, j int) bool { return len(pairs[i].val) > len(pairs[j].val) })

	out := content
	used := cfg.Config{}
	changed := false
	for _, p := range pairs {
		if p.val == "" {
			continue
		}
		if strings.Contains(out, p.val) {
			before := strings.Count(out, p.val)
			out = strings.ReplaceAll(out, p.val, p.ph)
			if before > 0 {
				used[p.ph] = p.val
				changed = true
			}
		}
	}
	return out, used, changed
}

// Line-level render that records state entries (for revert).
// It replaces placeholders -> values, tracking per-line changes.
func RenderLinesWithState(lines []string, c cfg.Config) ([]string, []st.Entry) {
	entries := make([]st.Entry, 0, 8)
	out := make([]string, len(lines))
	copy(out, lines)

	for i := range out {
		before := out[i]
		after := before
		reps := make([]st.Replacement, 0, 2)

		for ph, val := range c {
			if ph == "" {
				continue
			}
			if strings.Contains(after, ph) {
				cnt := strings.Count(after, ph)
				after = strings.ReplaceAll(after, ph, val)
				reps = append(reps, st.Replacement{Placeholder: ph, Value: val, Count: cnt})
			}
		}
		if after != before {
			out[i] = after
			entries = append(entries, st.Entry{
				Line:         i,
				Before:       before,
				After:        after,
				Replacements: reps,
			})
		}
	}
	return out, entries
}

// Apply state back: values -> placeholders on tracked lines
func RevertLinesByState(lines []string, entries []st.Entry) ([]string, cfg.Config) {
	out := make([]string, len(lines))
	copy(out, lines)
	regen := cfg.Config{}

	for _, e := range entries {
		if e.Line >= 0 && e.Line < len(out) {
			cur := out[e.Line]
			// if exact match, restore
			if cur == e.After {
				out[e.Line] = e.Before
			} else {
				// best-effort reverse
				for _, r := range e.Replacements {
					cur = strings.ReplaceAll(cur, r.Value, r.Placeholder)
				}
				out[e.Line] = cur
			}
			for _, r := range e.Replacements {
				regen[r.Placeholder] = r.Value
			}
		}
	}
	return out, regen
}
