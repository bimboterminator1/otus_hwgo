package hw06pipelineexecution

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In) (out Out)

func ExecutePipeline(in In, done In, stages ...Stage) Out {
	if len(stages) == 0 {
		return in
	}

	proxyChan := func(inChan In, doneChan In) Out {
		out := make(Bi)

		go func() {
			defer func() {
				close(out)
				//nolint:revive
				for range inChan {
				}
			}()
			for {
				select {
				case <-doneChan:
					return
				case v, ok := <-inChan:
					if !ok {
						return
					}
					out <- v
				}
			}
		}()
		return out
	}

	outChan := in
	for i := 0; i < len(stages); i++ {
		outChan = stages[i](proxyChan(outChan, done))
	}

	return outChan
}
