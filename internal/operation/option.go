package operation

import "github.com/devbackend/crutch-db/internal/storage"

type Option func(runner *Runner)

// WithStorage adding storage instance to Runner
func WithStorage(st *storage.Storage) Option {
	return func(runner *Runner) {
		runner.st = st
	}
}
