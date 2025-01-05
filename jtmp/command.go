package jtmp

import "bufio"

func WriteCommand(w *bufio.Writer, cmd string, root bool) {
	// Prefix all template commands with '# ' or '$ ' to indicate root and unprivileged execution respectively
	if root {
		w.WriteString("sudo ")
	}

	// Write the command itself followed by a newline
	w.WriteString(cmd)
	w.WriteRune('\n')
}
