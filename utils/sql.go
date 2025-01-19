package utils

import (
	"database/sql"
	"net/http"

	"github.com/ahsoromdoni/drone-patrol/generated"
	"github.com/labstack/echo/v4"
)

func CheckForNotFoundError(ctx echo.Context, err error, message string) error {
	if err == sql.ErrNoRows {
		if message == "" {
			message = err.Error()
		}

		return ctx.JSON(http.StatusNotFound, generated.ErrorResponse{Message: message})
	}
	return err
}
