package handler

import (
	"net/http"

	"github.com/ahsoromdoni/drone-patrol/generated"
	"github.com/ahsoromdoni/drone-patrol/repository"
	"github.com/labstack/echo/v4"
)

// This is just a test endpoint to get you started. Please delete this endpoint.
// (GET /hello)
func (s *Server) CreateEstate(ctx echo.Context) error {
	var req generated.EstateRequest

	// Parsing request body
	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, generated.ErrorResponse{Message: "Invalid request payload"})
	}

	// Validasi panjang dan lebar
	if req.Length <= 0 || req.Width <= 0 {
		return ctx.JSON(http.StatusBadRequest, generated.ErrorResponse{Message: "Length and width must be positive"})
	}

	estate, err := s.Repository.InsertEstate(ctx.Request().Context(), repository.CreateEstateInput{Length: req.Length, Width: req.Width})
	if err != nil {
		return err
	}

	resp := generated.EstateResponse{Id: estate.Id}
	return ctx.JSON(http.StatusCreated, resp)
}
