package execx

import "strings"

var promptMarkers = []struct {
	marker string
	id     string
}{
	{"Apple ID: ", "username"},
	{"Password: ", "password"},
	{"Code: ", "twofa"},
	{"Key ID: ", "keyid"},
	{"Issuer ID: ", "issuer"},
	{"Key path: ", "keypath"},
	{"Select a team", "team"},
	{"Continue? (yes/no): ", "confirm"},
	{"(yes/no): ", "confirm"},
	{"Choice (", "choice"},
}

// promptFor returns the prompt id and text for one partial output line.
func promptFor(text string) (string, string) {
	line := text
	if at := strings.LastIndexAny(line, "\r\n"); at >= 0 {
		line = line[at+1:]
	}
	for _, item := range promptMarkers {
		if strings.Contains(line, item.marker) {
			return item.id, strings.TrimSpace(line)
		}
	}
	return "", ""
}
