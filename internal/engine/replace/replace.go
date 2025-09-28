package replace

import (
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

func RenderLines(lines []string, c cfg.Config) ([]string, []st.Entry) {
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
