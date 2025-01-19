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

func TestServer_CreateTree(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockPostgres := repository.NewMockRepositoryInterface(ctrl)
	defer ctrl.Finish()

	getEstateByIdOutput := repository.GetEstateByIdOutput{}
	getEstateByIdOutputOK := repository.GetEstateByIdOutput{
		Id:     "123",
		Length: 10,
		Width:  10,
	}

	createTreeOutput := repository.CreateTreeOutput{}
	createTreeOutputOk := repository.CreateTreeOutput{Id: "abcdad-123an"}

	type args struct {
		id string
	}
	tests := []struct {
		name    string
		body    string
		args    args
		mock    func(args)
		wantErr bool
	}{
		{
			name: "success case",
			body: `{"x":1,"y":3,"height":30}`,
			args: args{
				id: "123",
			},
			mock: func(args) {
				mockPostgres.EXPECT().GetEstateById(gomock.Any(), gomock.Any()).Return(getEstateByIdOutputOK, nil)
				mockPostgres.EXPECT().InsertTree(gomock.Any(), gomock.Any()).Return(createTreeOutputOk, nil)
			},
			wantErr: false,
		},
		{
			name: "got error from GetEstateById case",
			body: `{"x":1,"y":3,"height":30}`,
			args: args{
				id: "123",
			},
			mock: func(args) {
				mockPostgres.EXPECT().GetEstateById(gomock.Any(), gomock.Any()).Return(getEstateByIdOutput, errors.New("error"))
			},
			wantErr: true,
		},
		{
			name: "got error from InsertTree case",
			body: `{"x":1,"y":3,"height":30}`,
			args: args{
				id: "123",
			},
			mock: func(args) {
				mockPostgres.EXPECT().GetEstateById(gomock.Any(), gomock.Any()).Return(getEstateByIdOutputOK, nil)
				mockPostgres.EXPECT().InsertTree(gomock.Any(), gomock.Any()).Return(createTreeOutput, errors.New("error"))
			},
			wantErr: true,
		},
		{
			name: "got error from IsTreeCordinateOutOfBound case",
			body: `{"x":11,"y":11,"height":30}`,
			args: args{
				id: "123",
			},
			mock: func(args) {
				mockPostgres.EXPECT().GetEstateById(gomock.Any(), gomock.Any()).Return(getEstateByIdOutputOK, nil)
			},
			wantErr: false,
		},
		{
			name: "got error from ctx.Bind case",
			body: `invalid json`,
			args: args{
				id: "123",
			},
			mock: func(args) {
			},
			wantErr: false,
		},
		{
			name: "got error from validation case",
			body: `{"x":1,"y":3}`,
			args: args{
				id: "123",
			},
			mock: func(args) {
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock(tt.args)
			s := &Server{Repository: mockPostgres}
			ctx := createMockContextWithPath(tt.body, echo.MIMEApplicationJSON, "/estate/123/tree", "id", tt.args.id)

			if err := s.CreateTree(ctx, tt.args.id); (err != nil) != tt.wantErr {
				t.Errorf("Server.CreateTree() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func createMockContextWithPath(body string, contentType string, path string, paramName string, paramValue string) echo.Context {
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, path, strings.NewReader(body))
	req.Header.Set(echo.HeaderContentType, contentType)
	rec := httptest.NewRecorder()
	ctx := e.NewContext(req, rec)

	ctx.SetParamNames(paramName)
	ctx.SetParamValues(paramValue)
	return ctx
}
