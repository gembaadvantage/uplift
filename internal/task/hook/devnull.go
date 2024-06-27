package hook

import "io"

// DevNull simulates the writing to /dev/null within a Linux OS,
// pinched from https://github.com/go-task/task/blob/master/internal/execext/devnull.go
type DevNull struct{}

func (DevNull) Read(_ []byte) (int, error) {
	return 0, io.EOF
}

func (DevNull) Write(p []byte) (int, error) {
	return len(p), nil
}

func (DevNull) Close() error {
	return nil
}
