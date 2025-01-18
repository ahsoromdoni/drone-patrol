package repository

import (
	"context"
)

var (
	insertTree = "INSERT INTO tree (estate_id, x_axis, y_axis, height) VALUES ($1, $2, $3, $4) RETURNING id"
)

func (r *Repository) InsertTree(ctx context.Context, input CreateTreeInput) (output CreateTreeOutput, err error) {
	err = r.Db.QueryRowContext(ctx, insertTree, input.EstateId, input.XAxis, input.YAxis, input.Height).Scan(&output.Id)
	if err != nil {
		return
	}
	return
}
