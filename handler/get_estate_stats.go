package handler

import (
	"net/http"

	"github.com/ahsoromdoni/drone-patrol/constant"
	"github.com/ahsoromdoni/drone-patrol/generated"
	"github.com/ahsoromdoni/drone-patrol/repository"
	"github.com/ahsoromdoni/drone-patrol/utils"
	"github.com/labstack/echo/v4"
)

// Get statistics for the estate
// (GET /estate/{id}/stats)
func (s *Server) GetEstateStats(ctx echo.Context, id string) error {
	var countTreeByEstateIdInput = repository.CountTreeByEstateIdInput{EstateId: id}
	count, err := s.Repository.CountTreeByEstateId(ctx.Request().Context(), countTreeByEstateIdInput)
	if err != nil {
		return utils.CheckForNotFoundError(ctx, err, constant.EstateNotFound)
	}

	var getTreeHeightStatsInput = repository.GetTreeHeightStatsInput{EstateId: id}
	tree, err := s.Repository.GetTreeHeightStats(ctx.Request().Context(), getTreeHeightStatsInput)
	if err != nil {
		return utils.CheckForNotFoundError(ctx, err, constant.EstateNotFound)
	}

	var resp = generated.EstateStatsResponse{
		Count:  count.TotalTrees,
		Max:    tree.MaxHeight,
		Min:    tree.MinHeight,
		Median: tree.MedianHeight,
	}
	return ctx.JSON(http.StatusOK, resp)
}
