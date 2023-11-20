package queue

import "fmt"

var (
	ErrorTopicNotFound    error = fmt.Errorf("topic not found")
	ErrorNoMessage        error = fmt.Errorf("no message in the topic")
	ErrorMessagesOverflow error = fmt.Errorf("messages limit was reached")
	ErrorTopicsOverflow   error = fmt.Errorf("topics limit was reached yet")
)
