package handler

import (
	"net/http"

	"github.com/ahsoromdoni/drone-patrol/generated"
	"github.com/ahsoromdoni/drone-patrol/repository"
	"github.com/labstack/echo/v4"
)

// This is just a test endpoint to get you started. Please delete this endpoint.
// (POST /estate/<id>/tree)
func (s *Server) CreateTree(ctx echo.Context, id string) error {
	var req generated.TreeRequest

	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, generated.ErrorResponse{Message: "Invalid request payload"})
	}

	getEstateByIdInput := repository.GetEstateByIdInput{
		Id: id,
	}
	estate, err := s.Repository.GetEstateById(ctx.Request().Context(), getEstateByIdInput)
	if err != nil {
		return err
	}

	if estate.IsCordinateOutOfBound(req.X, req.Y) {
		return ctx.JSON(http.StatusBadRequest, generated.ErrorResponse{Message: "Invalid request payload"})
	}

	createTreeInput := repository.CreateTreeInput{
		EstateId: id,
		XAxis:    req.X,
		YAxis:    req.Y,
		Height:   req.Height,
	}
	tree, err := s.Repository.InsertTree(ctx.Request().Context(), createTreeInput)
	if err != nil {
		return err
	}

	resp := generated.TreeResponse{Id: tree.Id}
	return ctx.JSON(http.StatusCreated, resp)
}

// Unimplement methods
func (s *Server) GetEstateStats(ctx echo.Context, id string) error {
	return nil
}

func (s *Server) GetDronePlan(ctx echo.Context, id string, params generated.GetDronePlanParams) error {
	return nil
}
