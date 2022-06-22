package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"strings"
	"sync"
	"time"

	"github.com/Shopify/sarama"
)

func main() {
	buf, err := ioutil.ReadFile("test.json")
	if err != nil {
		panic(fmt.Sprintf("json-file read failed %v", err))
	}
	buf = bytes.TrimRight(buf, "\n")
	brokerAddrs := strings.Split("localhost:9092;", ";")
	producer, err := sarama.NewAsyncProducer(brokerAddrs, sarama.NewConfig())
	if err != nil {
		panic(fmt.Sprintf("kafka producer init failed:%v", err))
	}
	go func() {
		err := <-producer.Errors()
		if err != nil {
			panic(err)
		}
	}()
	pwg := &sync.WaitGroup{}
	partions := 1
	for i := 0; i < partions; i++ {
		pwg.Add(1)
		go func(partition int) {
			defer pwg.Done()
			println("start push ", partition)
			for j := 0; j < 100_0000; j++ {
				producer.Input() <- &sarama.ProducerMessage{
					Topic:     "test",
					Value:     sarama.ByteEncoder(buf),
					Partition: int32(partition),
				}
			}
		}(i)
	}
	time.Sleep(2 * time.Second)
	pwg.Wait()
}
