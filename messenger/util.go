package messenger

import (
	"fmt"
	"strings"
)

// nstail returns the tail component of a namespace string.
func nstail(namespace string) string {
	if namespace == "" {
		return ""
	}

	doti := strings.LastIndexByte(namespace, '.')
	if doti == -1 {
		return namespace
	}

	return fmt.Sprintf(".%s.json", namespace[doti+1:])
}
