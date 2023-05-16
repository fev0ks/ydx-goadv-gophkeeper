package services

import (
	"bufio"
	"io"
	"os"
	"time"

	"go.uber.org/zap"

	"ydx-goadv-gophkeeper/internal/logger"
)

const (
	bufferSize = 4096
)

type FileService interface {
	ReadFile(path string, errCh chan error) (chan []byte, os.FileInfo, error)
	SaveFile(path string, chunks chan []byte) (chan error, error)
}

type fileService struct {
	log *zap.SugaredLogger
}

func NewFileService() FileService {
	return &fileService{log: logger.NewLogger("file-srv")}
}

func (fm *fileService) ReadFile(path string, errCh chan error) (chan []byte, os.FileInfo, error) {
	buf := make(chan []byte)
	file, err := os.Open(path)
	if err != nil {
		return buf, nil, err
	}

	stat, err := file.Stat()
	if err != nil {
		return nil, nil, err
	}

	go func() {
		defer file.Close()
		reader := bufio.NewReader(file)
		buffer := make([]byte, bufferSize)
		n := 0
		for {
			n, err = reader.Read(buffer)
			if err == io.EOF || n == 0 {
				close(buf)
				return
			}
			if err != nil {
				close(buf)
				return
			}

			select {
			case buf <- buffer[:n]:
			case _ = <-errCh:
				close(buf)
				return
			case <-time.After(1 * time.Minute):
				fm.log.Errorf("failed to read file: channel send timeout")
				close(buf)
				return
			}
		}
	}()

	return buf, stat, nil
}

func (fm *fileService) SaveFile(path string, chunks chan []byte) (chan error, error) {
	errCh := make(chan error)
	file, err := os.Create(path)
	if err != nil {
		return nil, err
	}
	go func() {
		writer := bufio.NewWriter(file)
		defer file.Close()
		defer writer.Flush()
		for {
			if bytes, ok := <-chunks; ok {
				_, err = writer.Write(bytes)
				if err != nil {
					fm.log.Errorf("failed to save file: %v", err)
					errCh <- err
					return
				}
				continue
			}
			break
		}
	}()
	return errCh, nil
}
