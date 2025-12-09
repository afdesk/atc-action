package provider

import (
	"errors"
)

var (
	ErrUnmarshal = errors.New("unmarshal error")
	ErrGeneral   = errors.New("weird error")
)

type MockContentProvider struct {
	Content string
	Err     error
}

func (mockContentProvider *MockContentProvider) GetContents(path string) (string, error) {
	return mockContentProvider.Content, mockContentProvider.Err
}
