package operation

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/devbackend/crutch-db/internal/storage"
)

type Runner struct {
	st *storage.Storage
}

// NewRunner return instance of Runner
func NewRunner(opts ...Option) *Runner {
	runner := new(Runner)

	for _, opt := range opts {
		opt(runner)
	}

	return runner
}

// Run for running operation over key value storage
func (r *Runner) Run(op *Operation) (string, error) {
	switch op.Type {
	case Get:
		stVal, ok := r.st.Get(op.Key)
		if !ok {
			return "", errors.New("not found")
		}

		switch val := stVal.(type) {
		case string:
			return val, nil
		case int:
			return strconv.Itoa(val), nil
		case float64:
			return strconv.FormatFloat(val, 'f', 4, 64), nil
		default:
			return "", errors.New("value parse error")
		}
	case Set:
		return "OK", r.st.Set(op.Key, op.Value)
	case Delete:
		return "OK", r.st.Delete(op.Key)
	case Keys:
		return fmt.Sprintf("[%s]", strings.Join(r.st.Keys(), ",")), nil
	}

	return "", errors.New("unknown operation for run")
}
