package service

import (
	"context"
	"fmt"
	"main/internal/queue"
	"time"
)

type QueueService interface {
	ReadFromQueue(ctx context.Context, topicName string, timeout *int) ([]byte, error)
	WriteToQueue(topicName string, msg []byte) error
}

type Service struct {
	queue          queue.SimpleQueue
	defaultTimeout int
}

func New(q queue.SimpleQueue, timeoutInSec int) (*Service, error) {
	if timeoutInSec < 0 {
		return nil, fmt.Errorf("invalid timeout value")
	}
	return &Service{q, timeoutInSec}, nil
}

func (s *Service) ReadFromQueue(ctx context.Context, topicName string, timeout *int) ([]byte, error) {

	t := time.Duration(s.defaultTimeout) * time.Second

	if timeout != nil {
		t = time.Duration(*timeout) * time.Second
	}

	ctxWithTimeout, _ := context.WithTimeout(ctx, t)

	msg, err := s.queue.Read(ctxWithTimeout, topicName)
	if err != nil {
		return nil, err
	}

	return msg, nil
}

func (s *Service) WriteToQueue(topicName string, msg []byte) error {

	if err := s.queue.Write(topicName, msg); err != nil {
		return err
	}

	return nil
}
