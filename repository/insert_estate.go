package repository

import (
	"context"
)

var (
	insertEstate = "INSERT INTO estate (length, width) VALUES ($1, $2) RETURNING id"
)

func (r *Repository) InsertEstate(ctx context.Context, input CreateEstateInput) (output CreateEstateOutput, err error) {
	err = r.Db.QueryRowContext(ctx, insertEstate, input.Length, input.Width).Scan(&output.Id)
	if err != nil {
		return
	}
	return
}
