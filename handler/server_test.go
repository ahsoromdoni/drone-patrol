package handler

import (
	"testing"

	"github.com/ahsoromdoni/drone-patrol/repository"
	gomock "github.com/golang/mock/gomock"
)

func TestNewServer(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepository := repository.NewMockRepositoryInterface(ctrl)

	opts := NewServerOptions{
		Repository: mockRepository,
	}

	server := NewServer(opts)

	if server == nil {
		t.Fatal("NewServer() returned nil, expected a valid Server instance")
	}

	if server.Repository != mockRepository {
		t.Errorf("NewServer().Repository = %v, expected %v", server.Repository, mockRepository)
	}
}
