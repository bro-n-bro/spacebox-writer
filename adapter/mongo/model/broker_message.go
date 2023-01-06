package model

import "time"

type BrokerMessage struct {
	Created          time.Time
	ID               string `bson:"_id"`
	LastErrorMessage string `bson:"last_error_message"`
	Topic            string
	Data             string
	Attempts         int
}
