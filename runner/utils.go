package main

import (
	"fmt"
	"io"
)

// prefixer implements the Writer interface,
// so that it can be used to prefix some text
// before something is sent to stdout and stderr
type prefixer struct {
	prefix string
	w      io.Writer
}

func (w *prefixer) Write(p []byte) (int, error) {
	_, err := fmt.Fprint(w.w, w.prefix)
	if err != nil {
		return 0, err
	}
	return w.w.Write(p)
}
