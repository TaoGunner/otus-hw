package hw06pipelineexecution

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In) (out Out)

func ExecutePipeline(in In, done In, stages ...Stage) Out {
	// Проверка на отсутствие этапов
	if len(stages) == 0 {
		chData := make(Bi)
		defer close(chData)
		return chData
	}

	chOut := in

	for _, stage := range stages {
		chData := make(Bi)

		go func(chData Bi, chOut Out) {
			defer func() {
				close(chData)
				//nolint:revive
				for range chOut {
					// Исправление утечки горутин
				}
			}()

			for {
				select {
				// Реализация остановки пайплайна
				case <-done:
					return
				case v, ok := <-chOut:
					// Канал закрыт - заканчиваем
					if !ok {
						return
					}
					chData <- v
				}
			}
		}(chData, chOut)

		chOut = stage(chData)
	}

	return chOut
}
