package model

import "time"

type BrokerMessage struct {
	ID               string `bson:"_id"`
	LastErrorMessage string `bson:"last_error_message"`
	Topic            string
	Created          time.Time
	Data             []byte
}
