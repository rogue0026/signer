package main

import (
	"sort"
	"strconv"
	"strings"
	"sync"
)

// сюда писать код

// пришло очередное значение - запускается конвейер
// и так до тех пор, пока есть значения
// если значений больше нет, надо прервать выполнение функции
// сделать это через канал
// значение отправляет производитель
func ExecutePipeline(jobs ...job) {
	in := make(chan interface{})
	out := make(chan interface{})
	wg := &sync.WaitGroup{}
	wg.Add(1)
	for _, joba := range jobs {
		go func(j job) {
			j(in, out)
		}(joba)
		in = out
		out = make(chan interface{})
	}
	go func() {
		wg.Wait()
	}()
}

// crc32(data)+"~"+crc32(md5(data))
func SingleHash(in chan interface{}, out chan interface{}) {
	for rawVal := range in {
		val, ok := rawVal.(int)
		if ok {
			ch1 := make(chan string)
			ch2 := make(chan string)

			md5Hash := DataSignerMd5(strconv.Itoa(val))

			// запускаем горутину для подсчета crc32(data)
			go func(ch1 chan string) {
				ch1 <- DataSignerCrc32(strconv.Itoa(val))
			}(ch1)

			// запускаем горутину для подсчета crc32 из хеша
			go func(ch2 chan string) {
				ch2 <- DataSignerCrc32(md5Hash)
			}(ch2)
			// получаем результаты из горутин
			out <- (<-ch1 + "~" + <-ch2)
		}
	}
	close(out)
}

// crc32(th+data)
func MultiHash(in chan interface{}, out chan interface{}) {
	type pair struct {
		th   int
		data string
	}

	for rawVal := range in {
		data, ok := rawVal.(string)
		if ok {
			wg := &sync.WaitGroup{}
			wg.Add(6)
			results := make(chan pair)

			// запускаем цикл из 6 горутин для вычисления хеша
			for i := 0; i < 6; i++ {
				go func(th int, singleHashString string) {
					defer wg.Done()
					hash := DataSignerCrc32(strconv.Itoa(th) + singleHashString)
					results <- pair{th, hash}
				}(i, data)
			}
			go func() {
				wg.Wait()
				close(results)
			}()

			res := make([]string, 6)
			for r := range results {
				res[r.th] = r.data
			}
			out <- strings.Join(res, "")
		}
	}
	close(out)
}

func CombineResults(in chan interface{}, out chan interface{}) {
	results := make([]string, 0)
	for rawVal := range in {
		data, ok := rawVal.(string)
		if ok {
			results = append(results, data)
		}
	}
	sort.Slice(results, func(i, j int) bool { return results[i] < results[j] })
	out <- strings.Join(results, "_")
	close(out)
}
