// Package shell builds command lines that are safe to send to a remote sh.
package shell

import (
	"encoding/base64"
	"fmt"
	"strings"
)

// Wrap encodes script as base64 and decodes it into sh on the remote side, so
// no value inside the script is ever parsed by the outer command line.
func Wrap(script string, sudo bool) string {
	sh := "sh"
	if sudo {
		sh = "sudo sh"
	}

	return fmt.Sprintf("printf '%%s' '%s' | base64 -d | %s", base64.StdEncoding.EncodeToString([]byte(script)), sh)
}

// Quote returns value as a single-quoted sh word.
func Quote(value string) string {
	return "'" + strings.ReplaceAll(value, "'", `'\''`) + "'"
}
