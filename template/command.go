package template

import "bufio"

func WriteCommand(w *bufio.Writer, cmd string, root bool) {
	// Prefix all template commands with '# ' or '$ ' to indicate root and unprivileged execution respectively
	if root {
		w.WriteString("# ")
	} else {
		w.WriteString("$ ")
	}

	// Write the command itself followed by a newline
	w.WriteString(cmd)
	w.WriteRune('\n')
}
