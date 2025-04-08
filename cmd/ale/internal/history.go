package internal

import (
	"os/user"
	"path"
	"strings"

	"github.com/kode4food/ale/pkg/data"
)

func (r *REPL) saveHistory() {
	defer func() { _ = recover() }()
	seq := r.scanBuffer()
	hist := toHistory(seq)
	_ = r.rl.SaveHistory(hist)
}

func getHistoryFile() string {
	if usr, err := user.Current(); err == nil {
		return path.Join(usr.HomeDir, ".ale-history")
	}
	return ""
}

func toHistory(s data.Sequence) string {
	var buf strings.Builder
	for f, r, ok := s.Split(); ok; f, r, ok = r.Split() {
		if buf.Len() > 0 {
			buf.WriteString(" ")
		}
		buf.WriteString(data.ToQuotedString(f))
	}
	return buf.String()
}
