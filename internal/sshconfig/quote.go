package sshconfig

import "strings"

// quoteSSH wraps s in double quotes when it contains whitespace or quotes.
func quoteSSH(s string) string {
	if s == "" || !strings.ContainsAny(s, " \t\"'") {
		return s
	}
	return `"` + strings.ReplaceAll(s, `"`, `\"`) + `"`
}

// unquoteSSH removes one layer of double quotes if present.
func unquoteSSH(s string) string {
	s = strings.TrimSpace(s)
	if len(s) >= 2 && s[0] == '"' && s[len(s)-1] == '"' {
		return strings.ReplaceAll(s[1:len(s)-1], `\"`, `"`)
	}
	return s
}

// splitSSHTokens splits on whitespace, respecting double-quoted tokens.
func splitSSHTokens(s string) []string {
	s = strings.TrimSpace(s)
	if s == "" {
		return nil
	}
	var out []string
	for s != "" {
		if s[0] == '"' {
			end := 1
			closed := false
			for end < len(s) {
				if s[end] == '\\' && end+1 < len(s) {
					end += 2
					continue
				}
				if s[end] == '"' {
					out = append(out, unquoteSSH(s[:end+1]))
					s = strings.TrimSpace(s[end+1:])
					closed = true
					break
				}
				end++
			}
			if !closed {
				out = append(out, unquoteSSH(s))
				break
			}
			continue
		}
		i := strings.IndexAny(s, " \t")
		if i < 0 {
			out = append(out, s)
			break
		}
		out = append(out, s[:i])
		s = strings.TrimSpace(s[i:])
	}
	return out
}
