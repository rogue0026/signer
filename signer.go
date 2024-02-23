package main

import (
	"strconv"
	"sync"
)

// сюда писать код

// crc(data) + crc(hash(data))
var SingleHash = func(in, out chan interface{}) {
	wg := sync.WaitGroup{}
	//m := sync.Mutex{}
	for val := range in {
		if n, ok := val.(int); ok {
			data := strconv.Itoa(n)
			hashData := DataSignerMd5(data)

			wg.Add(1)
			go func(normalData string, hash string) {
				defer wg.Done()
				hashCRCCh := make(chan string, 1)
				dataCRCCh := make(chan string, 1)
				go func(hashCRCCh chan string) {
					hashCRC := DataSignerCrc32(hash)
					hashCRCCh <- hashCRC
				}(hashCRCCh)
				go func(dataCRCCh chan string) {
					dataCRC := DataSignerCrc32(normalData)
					dataCRCCh <- dataCRC
				}(dataCRCCh)
				dCRC := <-dataCRCCh
				hCRC := <-hashCRCCh
				out <- dCRC + "~" + hCRC
			}(data, hashData)
		}
	}
	wg.Wait()
	close(out)
	// дождаться пока все отправят данные
}
var MultiHash = func(in, out chan interface{}) {}
var CombineResults = func(in, out chan interface{}) {}
