package main

import (
	"context"
	"runtime/pprof"
)

func RoutineAnnotator(f func(), labels []string) {
	pprof.Do(context.Background(), pprof.Labels(labels...), func(_ context.Context) {
		go f()
	})
}
