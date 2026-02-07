package config

import (
	"time"
)

var (
	appVersion string = "v0.2.1"
	commitHash string = "Unkown"
	buildDate  string = "Unkown"
)

func parseBuildTime(s string) string {
	if t, err := time.Parse(time.RFC3339, s); err != nil {
		return "Unkown"
	} else {
		return t.Local().Format(time.DateTime + " -07:00")
	}
}
