package repository

import (
	"context"
	"database/sql"
)

var (
	treeHeightStatsByEstateId = `
        SELECT 
			COALESCE(MAX(height), 0) AS max_height, 
			COALESCE(MIN(height), 0) AS min_height, 
			COALESCE(PERCENTILE_DISC(0.5) WITHIN GROUP (ORDER BY height), 0) AS median
        FROM tree 
        WHERE estate_id = $1`
)

func (r *Repository) GetTreeHeightStats(ctx context.Context, input GetTreeHeightStatsInput) (output GetTreeHeightStatsOutput, err error) {
	err = r.Db.QueryRowContext(ctx, treeHeightStatsByEstateId, input.EstateId).Scan(
		&output.MaxHeight,
		&output.MinHeight,
		&output.MedianHeight,
	)

	if output.MaxHeight == 0 && output.MinHeight == 0 && output.MedianHeight == 0 {
		err = sql.ErrNoRows
		return
	}
	if err != nil {
		return
	}
	return
}
