package mock

import (
	"context"
	"github.com/stretchr/testify/mock"
)

type ClientMock struct {
	mock.Mock
}

func (m *ClientMock) Do(ctx context.Context, method, url string, body []byte) ([]byte, error) {
	args := m.Called(ctx, method, url, body)
	return args.Get(0).([]byte), args.Error(1)
}
