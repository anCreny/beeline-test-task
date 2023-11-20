package queue

import (
	"context"
	"errors"
)

type SimpleQueue interface {
	Read(ctx context.Context, topicName string) ([]byte, error)
	Write(topicName string, msg []byte) error
}

type Queue struct {
	messagesLimit int
	topicsLimit   int
	topics        map[string]chan []byte
}

func New(messagesLimit int, topicsLimit int) *Queue {
	return &Queue{
		messagesLimit: messagesLimit,
		topicsLimit:   topicsLimit,
		topics:        make(map[string]chan []byte),
	}
}

func (q *Queue) Read(ctx context.Context, topicName string) ([]byte, error) {
	topic, found := q.topics[topicName]
	if !found {
		return nil, ErrorTopicNotFound
	}

	select {
	case msg := <-topic:
		return msg, nil
	case <-ctx.Done():
		if errors.Is(ctx.Err(), context.DeadlineExceeded) {
			return nil, ErrorNoMessage
		}
		return nil, nil
	}
}

func (q *Queue) Write(topicName string, msg []byte) error {
	if _, found := q.topics[topicName]; !found {
		if len(q.topics) == q.topicsLimit {
			return ErrorTopicsOverflow
		}

		q.topics[topicName] = make(chan []byte, q.messagesLimit)
	}

	// catcing messages limit overflowing
	select {
	case q.topics[topicName] <- msg:
	default:
		return ErrorMessagesOverflow
	}

	return nil
}
