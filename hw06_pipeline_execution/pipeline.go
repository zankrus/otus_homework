package hw06pipelineexecution

import "log"

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In) (out Out)

func ExecutePipeline(in In, done In, stages ...Stage) Out {
	out := make(Bi)

	if len(stages) == 0 {
		close(out)
		return out
	}

	go func() {
		defer close(out)

		stream := in

		for _, stage := range stages {
			stream = stage(stream)
		}

		for {
			select {
			case <-done:
				return
			case i, ok := <-stream:
				log.Println(i)
				if !ok {
					return
				}
				out <- i
			}
		}

	}()

	return out
}
