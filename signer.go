package main

import (
	"strconv"
	"strings"
)

func ExecutePipeline(jobs ...job) {
	in := make(chan interface{})
	out := make(chan interface{})
	for _, f := range jobs {
		go f(in, out)
		in = out
		out = make(chan interface{})
	}

}

// crc32(data)+"~"+crc32(md5(data))
func SingleHash(inputData chan interface{}, out chan interface{}) {

	hashCRCCh := make(chan string)
	dataCRCCh := make(chan string)
	for rawData := range inputData {
		if num, ok := rawData.(int); ok {

			strNum := strconv.Itoa(num)
			hash := DataSignerMd5(strNum)

			go func(hashCRCCh chan string) {
				hashCRC := DataSignerCrc32(hash)
				hashCRCCh <- hashCRC
			}(hashCRCCh)

			go func(dataCRCCh chan string) {
				dataCRC := DataSignerCrc32(strNum)
				dataCRCCh <- dataCRC
			}(dataCRCCh)

			out <- (<-dataCRCCh + "~" + <-hashCRCCh)
		}
	}
	close(out)
}

// crc32(th+data)
func MultiHash(in chan interface{}, out chan interface{}) {
	for rawData := range in {
		if hashVal, ok := rawData.(string); ok {
			chans := make([]chan string, 6)
			for i := 0; i < len(chans); i++ {
				go func(th int) {
					chans[th] <- DataSignerCrc32(strconv.Itoa(th) + hashVal)
				}(i)
			}
			result := make([]string, 6)
			for i := 0; i < len(chans); i++ {
				result[i] = <-chans[i]
			}
			out <- strings.Join(result, "")
		}
	}
	close(out)
}

func CombineResults(in chan interface{}, out chan interface{}) {}
