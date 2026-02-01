package updater

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"
)

type LatestRelease struct {
	Tag     string
	Body    string
	HTMLURL string
}

// CheckUpdate fetches the latest release for `repo` (owner/name) and compares
// it against `currentVersion`. It returns (isNewer, latestRelease, error).
func CheckUpdate(repo, currentVersion string) (bool, LatestRelease, error) {
	client := &http.Client{Timeout: 8 * time.Second}
	url := "https://api.github.com/repos/" + repo + "/releases/latest"
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return false, LatestRelease{}, err
	}
	req.Header.Set("User-Agent", "vexo-updater")

	resp, err := client.Do(req)
	if err != nil {
		return false, LatestRelease{}, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return false, LatestRelease{}, nil
	}

	var data struct {
		TagName string `json:"tag_name"`
		Body    string `json:"body"`
		HTMLURL string `json:"html_url"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return false, LatestRelease{}, err
	}

	latest := data.TagName
	cur := normalize(currentVersion)
	lat := normalize(latest)

	if compareSemver(lat, cur) > 0 {
		return true, LatestRelease{Tag: latest, Body: data.Body, HTMLURL: data.HTMLURL}, nil
	}
	return false, LatestRelease{Tag: latest, Body: data.Body, HTMLURL: data.HTMLURL}, nil
}

func normalize(s string) string {
	if len(s) > 0 && (s[0] == 'v' || s[0] == 'V') {
		return s[1:]
	}
	return s
}

// compareSemver compares two semantic versions (simple numeric comparison).
// returns 1 if a>b, -1 if a<b, 0 if equal.
func compareSemver(a, b string) int {
	if a == b {
		return 0
	}
	as := strings.Split(a, ".")
	bs := strings.Split(b, ".")
	n := len(as)
	if len(bs) > n {
		n = len(bs)
	}
	for i := 0; i < n; i++ {
		var ai, bi int
		if i < len(as) {
			fmtS := strings.TrimSpace(as[i])
			for j := 0; j < len(fmtS); j++ {
				if fmtS[j] < '0' || fmtS[j] > '9' {
					fmtS = fmtS[:j]
					break
				}
			}
			if fmtS != "" {
				ai = atoi(fmtS)
			}
		}
		if i < len(bs) {
			fmtS := strings.TrimSpace(bs[i])
			for j := 0; j < len(fmtS); j++ {
				if fmtS[j] < '0' || fmtS[j] > '9' {
					fmtS = fmtS[:j]
					break
				}
			}
			if fmtS != "" {
				bi = atoi(fmtS)
			}
		}
		if ai > bi {
			return 1
		}
		if ai < bi {
			return -1
		}
	}
	return 0
}

func atoi(s string) int {
	var v int
	for i := 0; i < len(s); i++ {
		c := s[i]
		if c < '0' || c > '9' {
			break
		}
		v = v*10 + int(c-'0')
	}
	return v
}
