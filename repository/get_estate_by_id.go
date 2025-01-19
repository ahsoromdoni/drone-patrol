package repository

import (
	"context"
)

var (
	getEstateById = "SELECT id, length, width FROM estate WHERE id = $1"
)

func (r *Repository) GetEstateById(ctx context.Context, input GetEstateByIdInput) (output GetEstateByIdOutput, err error) {
	row := r.Db.QueryRowContext(ctx, getEstateById, input.Id)
	err = row.Scan(
		&output.Id,
		&output.Length,
		&output.Width,
	)
	if err != nil {
		return
	}
	return
}
