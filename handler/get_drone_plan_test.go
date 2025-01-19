package handler

import (
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"

	"github.com/ahsoromdoni/drone-patrol/generated"
	"github.com/ahsoromdoni/drone-patrol/repository"
	gomock "github.com/golang/mock/gomock"
	"github.com/labstack/echo/v4"
)

func TestServer_GetDronePlan(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockPostgres := repository.NewMockRepositoryInterface(ctrl)
	defer ctrl.Finish()

	getEstateByIdOutputOK := repository.GetEstateByIdOutput{Id: "123", Length: 5, Width: 3}

	tree1 := repository.GetTreesByEstateIdOutput{Id: "abc", EstateId: "123", XAxis: 2, YAxis: 1, Height: 5}
	tree2 := repository.GetTreesByEstateIdOutput{Id: "abc", EstateId: "123", XAxis: 5, YAxis: 2, Height: 3}
	tree3 := repository.GetTreesByEstateIdOutput{Id: "abc", EstateId: "123", XAxis: 1, YAxis: 3, Height: 6}
	tree4 := repository.GetTreesByEstateIdOutput{Id: "abc", EstateId: "123", XAxis: 2, YAxis: 4, Height: 15}
	getTreesByEstateIdOutputOk := []repository.GetTreesByEstateIdOutput{tree1, tree2, tree3, tree4}

	maxDistance := 180

	type args struct {
		id     string
		params generated.GetDronePlanParams
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
			args: args{
				id:     "123",
				params: generated.GetDronePlanParams{MaxDistance: nil},
			},
			mock: func(args) {
				mockPostgres.EXPECT().GetEstateById(gomock.Any(), gomock.Any()).Return(getEstateByIdOutputOK, nil)
				mockPostgres.EXPECT().GetTreesByEstateId(gomock.Any(), gomock.Any()).Return(getTreesByEstateIdOutputOk, nil)
			},
			wantErr: false,
		},
		{
			name: "success case with params",
			body: "",
			args: args{
				id:     "123",
				params: generated.GetDronePlanParams{MaxDistance: &maxDistance},
			},
			mock: func(args) {
				mockPostgres.EXPECT().GetEstateById(gomock.Any(), gomock.Any()).Return(getEstateByIdOutputOK, nil)
				mockPostgres.EXPECT().GetTreesByEstateId(gomock.Any(), gomock.Any()).Return(getTreesByEstateIdOutputOk, nil)
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock(tt.args)
			s := &Server{Repository: mockPostgres}
			ctx := createMockContextForGetDronePlan(tt.body, echo.MIMEApplicationJSON, tt.args.id, tt.args.params)

			if err := s.GetDronePlan(ctx, tt.args.id, tt.args.params); (err != nil) != tt.wantErr {
				t.Errorf("Server.GetDronePlan() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func createMockContextForGetDronePlan(body string, contentType string, id string, params generated.GetDronePlanParams) echo.Context {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/estate/"+id+"/droneplan", strings.NewReader(body))
	req.Header.Set(echo.HeaderContentType, contentType)

	queryString := ""
	if params.MaxDistance != nil {
		queryString = "maxDistance=" + strconv.Itoa(*params.MaxDistance)
	}

	if queryString != "" {
		req.URL.RawQuery = queryString
	}

	rec := httptest.NewRecorder()
	return e.NewContext(req, rec)
}
