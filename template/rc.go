package template

import "bufio"

func WriteRc(w *bufio.Writer, key, value string) {
	w.WriteString(key)
	w.WriteRune('=')
	w.WriteRune('"')
	w.WriteString(value)
	w.WriteRune('"')

	w.WriteRune('\n')
}
