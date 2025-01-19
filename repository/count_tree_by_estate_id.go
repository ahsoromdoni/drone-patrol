package repository

import (
	"context"
	"database/sql"
)

var (
	countTreeByEstateId = "SELECT COUNT(id) AS total_trees FROM tree WHERE estate_id = $1"
)

func (r *Repository) CountTreeByEstateId(ctx context.Context, input CountTreeByEstateIdInput) (output CountTreeByEstateIdOutput, err error) {
	err = r.Db.QueryRowContext(ctx, countTreeByEstateId, input.EstateId).Scan(&output.TotalTrees)

	if output.TotalTrees == 0 {
		err = sql.ErrNoRows
		return
	}
	if err != nil {
		return
	}
	return
}
