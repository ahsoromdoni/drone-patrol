package repository

import (
	"context"
)

func (r *Repository) GetEstateById(ctx context.Context, input GetEstateByIdInput) (output GetEstateByIdOutput, err error) {
	row := r.Db.QueryRowContext(ctx, "SELECT id, length, width FROM estate WHERE id = $1", input.Id)
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
