package handler

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/ahsoromdoni/drone-patrol/repository"
	gomock "github.com/golang/mock/gomock"
	"github.com/labstack/echo/v4"
)

func TestServer_CreateEstate(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockPostgres := repository.NewMockRepositoryInterface(ctrl)
	defer ctrl.Finish()

	tests := []struct {
		name    string
		body    string
		mock    func()
		wantErr bool
	}{
		{
			name: "success case",
			body: `{"length":6,"width":3}`,
			mock: func() {
				mockPostgres.EXPECT().InsertEstate(gomock.Any(), gomock.Any()).Return(repository.CreateEstateOutput{Id: "abcdad-123an"}, nil)
			},
			wantErr: false,
		},
		{
			name: "got error from InsertEstate case",
			body: `{"length":6,"width":3}`,
			mock: func() {
				mockPostgres.EXPECT().InsertEstate(gomock.Any(), gomock.Any()).Return(repository.CreateEstateOutput{}, errors.New("error"))
			},
			wantErr: true,
		},
		{
			name: "got error from validation case",
			body: `{"length":6}`,
			mock: func() {
			},
			wantErr: false,
		},
		{
			name: "got error from ctx.Bind case",
			body: `invalid json`,
			mock: func() {
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()
			s := &Server{Repository: mockPostgres}
			ctx := createMockContext(tt.body, echo.MIMEApplicationJSON)

			if err := s.CreateEstate(ctx); (err != nil) != tt.wantErr {
				t.Errorf("Server.CreateEstate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func createMockContext(body string, contentType string) echo.Context {
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/estate", strings.NewReader(body))
	req.Header.Set(echo.HeaderContentType, contentType)
	rec := httptest.NewRecorder()
	return e.NewContext(req, rec)
}
