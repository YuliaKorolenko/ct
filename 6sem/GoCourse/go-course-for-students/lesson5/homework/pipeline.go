package executor

import (
	"context"
)

type (
	In  <-chan any
	Out = In
)

type Stage func(in In) (out Out)

func ExecutePipeline(ctx context.Context, in In, stages ...Stage) Out {

	out := in

	for _, stage := range stages {
		forCommunication := make(chan any)
		out = stage(out)
		go func(curOut Out, c chan any) {
			defer close(c)
			for {
				select {
				case <-ctx.Done():
					return
				case v, ok := <-curOut:
					if ok {
						c <- v
					} else {
						return
					}

				}
			}
		}(out, forCommunication)
		out = forCommunication
	}
	return out
}
