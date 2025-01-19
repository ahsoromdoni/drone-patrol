// This file contains the interfaces for the repository layer.
// The repository layer is responsible for interacting with the database.
// For testing purpose we will generate mock implementations of these
// interfaces using mockgen. See the Makefile for more information.
package repository

import "context"

type RepositoryInterface interface {
	InsertEstate(ctx context.Context, input CreateEstateInput) (output CreateEstateOutput, err error)
	GetEstateById(ctx context.Context, input GetEstateByIdInput) (output GetEstateByIdOutput, err error)
	InsertTree(ctx context.Context, input CreateTreeInput) (output CreateTreeOutput, err error)
	CountTreeByEstateId(ctx context.Context, input CountTreeByEstateIdInput) (output CountTreeByEstateIdOutput, err error)
	GetTreeHeightStats(ctx context.Context, input GetTreeHeightStatsInput) (output GetTreeHeightStatsOutput, err error)
}
