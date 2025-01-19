package handler

import (
	"net/http"

	"github.com/ahsoromdoni/drone-patrol/generated"
	"github.com/ahsoromdoni/drone-patrol/repository"
	"github.com/ahsoromdoni/drone-patrol/validation"
	"github.com/labstack/echo/v4"
)

func (s *Server) CreateEstate(ctx echo.Context) error {
	var req generated.EstateRequest

	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, generated.ErrorResponse{Message: "Invalid request payload"})
	}

	if err := validation.ValidateCreateEstate(req); err != nil {
		return ctx.JSON(http.StatusBadRequest, generated.ErrorResponse{Message: err.Error()})
	}

	var createEstateInput = repository.CreateEstateInput{Length: req.Length, Width: req.Width}
	estate, err := s.Repository.InsertEstate(ctx.Request().Context(), createEstateInput)
	if err != nil {
		return err
	}

	var resp = generated.EstateResponse{Id: estate.Id}
	return ctx.JSON(http.StatusCreated, resp)
}
