package utils

import "sync"

// Generator создает канал-генератор.
func Generator[T any](doneCh chan struct{}, URLs []T) chan T {
	result := make(chan T)

	go func() {
		defer close(result)

		for _, URL := range URLs {
			select {
			case <-doneCh:
				return
			case result <- URL:
			}
		}
	}()

	return result
}

// FanOut создает выполняет некоторую работу на переданном количестве воркеров.
func FanOut[T any](workersCnt int, work func() chan T) []chan T {
	workers := make([]chan T, workersCnt)

	for i := 0; i < workersCnt; i++ {
		workers[i] = work()
	}

	return workers
}

// FanIn дожидается окончания работы переданных каналов и возвращает результирующий набор итоговых обработанных данных.
func FanIn[T any](doneCh chan struct{}, resultChannels []chan T) chan T {
	finalCh := make(chan T)

	var wg sync.WaitGroup

	for _, ch := range resultChannels {
		chClosure := ch

		wg.Add(1)

		go func() {
			defer wg.Done()

			for res := range chClosure {
				select {
				case <-doneCh:
					return
				case finalCh <- res:
				}
			}
		}()
	}

	go func() {
		wg.Wait()
		close(finalCh)
	}()

	return finalCh
}
