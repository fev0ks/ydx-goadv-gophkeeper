package services

import (
	"go.uber.org/zap"

	"ydx-goadv-gophkeeper/pkg/logger"
)

type FileProcessor interface {
	Receive(stream func()) ([]byte, error)
	Send(stream func([]byte)) error
}

type fileProcessor struct {
	log       *zap.SugaredLogger
	chunkSize int
}

func NewFileProcessor() FileProcessor {
	return &fileProcessor{
		log:       logger.NewLogger("file_processor"),
		chunkSize: 4096,
	}
}

func (f fileProcessor) Receive(stream func()) ([]byte, error) {
	return nil, nil
}

func (f fileProcessor) Send(stream func([]byte)) error {
	return nil
}
