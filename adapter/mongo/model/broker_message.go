package model

import "time"

// BrokerMessage is a message that is sent to the broker
type BrokerMessage struct {
	Created          time.Time
	ID               string `bson:"_id"`
	LastErrorMessage string `bson:"last_error_message"`
	Topic            string
	Data             string
	Attempts         int
}
