package broker

import "github.com/confluentinc/confluent-kafka-go/kafka"

func bytesToInt(bytes []byte) (value int) {
	for _, bt := range bytes {
		value = value*10 + int(bt-48)
	}
	return
}

func findValueFromHeaders(key string, headers []kafka.Header) (res []byte) {
	for _, header := range headers {
		if header.Key == key {
			return header.Value
		}
	}
	return
}
