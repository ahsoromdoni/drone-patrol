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

func TestServer_GetEstateStats(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockPostgres := repository.NewMockRepositoryInterface(ctrl)
	defer ctrl.Finish()

	countTreeByEstateIdOutput := repository.CountTreeByEstateIdOutput{}
	countTreeByEstateIdOutputOK := repository.CountTreeByEstateIdOutput{TotalTrees: 10}

	getTreeHeightStatsOutput := repository.GetTreeHeightStatsOutput{}
	getTreeHeightStatsOutputOK := repository.GetTreeHeightStatsOutput{MaxHeight: 20, MinHeight: 5, MedianHeight: 12}

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
			body: "",
			args: args{id: "123"},
			mock: func(args) {
				mockPostgres.EXPECT().CountTreeByEstateId(gomock.Any(), gomock.Any()).Return(countTreeByEstateIdOutputOK, nil)
				mockPostgres.EXPECT().GetTreeHeightStats(gomock.Any(), gomock.Any()).Return(getTreeHeightStatsOutputOK, nil)
			},
			wantErr: false,
		},
		{
			name: "got error from CountTreeByEstateId case",
			body: "",
			args: args{id: "123"},
			mock: func(args) {
				mockPostgres.EXPECT().CountTreeByEstateId(gomock.Any(), gomock.Any()).Return(countTreeByEstateIdOutput, errors.New("error"))
			},
			wantErr: true,
		},
		{
			name: "success case",
			body: "",
			args: args{id: "123"},
			mock: func(args) {
				mockPostgres.EXPECT().CountTreeByEstateId(gomock.Any(), gomock.Any()).Return(countTreeByEstateIdOutputOK, nil)
				mockPostgres.EXPECT().GetTreeHeightStats(gomock.Any(), gomock.Any()).Return(getTreeHeightStatsOutput, errors.New("error"))
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock(tt.args)
			s := &Server{Repository: mockPostgres}

			ctx := createMockContextForGetEstateStats(tt.body, echo.MIMEApplicationJSON, tt.args.id)

			if err := s.GetEstateStats(ctx, tt.args.id); (err != nil) != tt.wantErr {
				t.Errorf("Server.GetEstateStats() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func createMockContextForGetEstateStats(body string, contentType string, id string) echo.Context {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/estate/"+id+"/stats", strings.NewReader(body))
	req.Header.Set(echo.HeaderContentType, contentType)
	rec := httptest.NewRecorder()
	return e.NewContext(req, rec)
}
