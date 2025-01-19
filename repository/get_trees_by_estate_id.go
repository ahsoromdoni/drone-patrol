package repository

import (
	"context"
)

var (
	getTreesByEstateId = "SELECT id, estate_id, x_axis, y_axis, height FROM tree WHERE estate_id = $1"
)

func (r *Repository) GetTreesByEstateId(ctx context.Context, input GetTreesByEstateIdInput) (output []GetTreesByEstateIdOutput, err error) {
	rows, err := r.Db.QueryContext(ctx, getTreesByEstateId, input.EstateId)
	if err != nil {
		return
	}
	defer rows.Close()

	for rows.Next() {
		var tree GetTreesByEstateIdOutput
		err = rows.Scan(
			&tree.Id,
			&tree.EstateId,
			&tree.XAxis,
			&tree.YAxis,
			&tree.Height,
		)
		if err != nil {
			return
		}
		output = append(output, tree)
	}

	return
}
