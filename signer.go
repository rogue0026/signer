package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"sort"
	"strconv"
	"strings"
)

// сюда писать код
func ExecutePipeline(j ...job) {
	for _, task := range j {
		in := make(chan interface{})
		out := make(chan interface{})
		go task(in, out)
	}
}

func gen(in chan interface{}, out chan interface{}) {
	for i := 0; i < 100; i++ {
		out <- rand.Intn(100)
	}
}

func SingleHash(inputCh chan interface{}, outCh chan interface{}) {

	for {
		rawData, open := <-inputCh
		if open {
			str, asserted := rawData.(string)
			if asserted {
				crc := DataSignerCrc32(str)
				// time.Sleep(time.Millisecond * 10)
				crcWithHash := DataSignerCrc32(DataSignerMd5(str))
				outCh <- (crc + "~" + crcWithHash)
			}
		} else {
			close(outCh)
			break
		}
	}
}

func MultiHash(inputCh chan interface{}, outCh chan interface{}) {
	for {
		rawData, open := <-inputCh
		if open {
			singHash, asserted := rawData.(string)
			if asserted {
				buffer := bytes.NewBuffer(make([]byte, 0))
				for i := 0; i < 6; i++ {
					buffer.Write([]byte(DataSignerCrc32(strconv.Itoa(i) + singHash)))
				}
				outCh <- buffer.String()
			}
		} else {
			close(outCh)
			break
		}
	}
}

func CombineResults(inputCh chan interface{}, outCh chan interface{}) {
	results := make([]string, 100) // делаем размер слайса 100, чтобы не тратить время на пересоздание
	for {
		rawData, open := <-inputCh
		if open {
			if res, asserted := rawData.(string); asserted {
				results = append(results, res)
			}
		} else {
			close(outCh) // подумать над тем, кто будет закрывать выходной канал
			// входной канал для первой функции может быть выходным для второй
			break
		}
	}
	sort.Slice(results, func(i, j int) bool { return results[i] < results[j] })
	fmt.Println(strings.Join(results, "_"))
}
