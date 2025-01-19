package handler

import (
	"net/http"

	"github.com/ahsoromdoni/drone-patrol/constant"
	"github.com/ahsoromdoni/drone-patrol/generated"
	"github.com/ahsoromdoni/drone-patrol/repository"
	"github.com/ahsoromdoni/drone-patrol/utils"
	"github.com/ahsoromdoni/drone-patrol/validation"
	"github.com/labstack/echo/v4"
)

func (s *Server) CreateTree(ctx echo.Context, id string) error {
	var req generated.TreeRequest

	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, generated.ErrorResponse{Message: constant.BadRequest})
	}

	if err := validation.ValidateCreateTree(req); err != nil {
		return ctx.JSON(http.StatusBadRequest, generated.ErrorResponse{Message: err.Error()})
	}

	var getEstateByIdInput = repository.GetEstateByIdInput{Id: id}
	estate, err := s.Repository.GetEstateById(ctx.Request().Context(), getEstateByIdInput)
	if err != nil {
		return utils.CheckForNotFoundError(ctx, err, constant.EstateNotFound)
	}

	if estate.IsCordinateOutOfBound(req.X, req.Y) {
		return ctx.JSON(http.StatusBadRequest, generated.ErrorResponse{Message: constant.BadRequest})
	}

	var createTreeInput = repository.CreateTreeInput{
		EstateId: id,
		XAxis:    req.X,
		YAxis:    req.Y,
		Height:   req.Height,
	}
	tree, err := s.Repository.InsertTree(ctx.Request().Context(), createTreeInput)
	if err != nil {
		return err
	}

	var resp = generated.TreeResponse{Id: tree.Id}
	return ctx.JSON(http.StatusCreated, resp)
}

func (s *Server) GetDronePlan(ctx echo.Context, id string, params generated.GetDronePlanParams) error {
	return nil
}
