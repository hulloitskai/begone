package loader

import (
	"time"
)

// timeString returns the current time as a kebab-cased strign.
func timeString() string {
	return time.Now().Format("2006-01-02-150405")
}
