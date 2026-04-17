package prompt

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
)

// Confirmer asks the user for a yes/no confirmation.
type Confirmer struct {
	In  io.Reader
	Out io.Writer
}

// New returns a Confirmer reading from stdin and writing to stdout.
func New() *Confirmer {
	return &Confirmer{In: os.Stdin, Out: os.Stdout}
}

// Ask prints the question and returns true if the user answers "y" or "yes".
// If defaultYes is true, an empty response is treated as yes.
func (c *Confirmer) Ask(question string, defaultYes bool) (bool, error) {
	hint := "y/N"
	if defaultYes {
		hint = "Y/n"
	}
	fmt.Fprintf(c.Out, "%s [%s]: ", question, hint)

	scanner := bufio.NewScanner(c.In)
	if !scanner.Scan() {
		if err := scanner.Err(); err != nil {
			return false, err
		}
		// EOF — treat as default
		return defaultYes, nil
	}

	answer := strings.TrimSpace(strings.ToLower(scanner.Text()))
	switch answer {
	case "":
		return defaultYes, nil
	case "y", "yes":
		return true, nil
	case "n", "no":
		return false, nil
	default:
		return false, fmt.Errorf("unrecognised answer %q", answer)
	}
}
