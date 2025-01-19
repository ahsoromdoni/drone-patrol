// This file contains types that are used in the repository layer.
package repository

type GetTestByIdInput struct {
	Id string
}

type GetTestByIdOutput struct {
	Name string
}

type CreateEstateInput struct {
	Length int
	Width  int
}

type CreateEstateOutput struct {
	Id string
}

type GetEstateByIdInput struct {
	Id string
}

type GetEstateByIdOutput struct {
	Id     string
	Length int
	Width  int
}

type CreateTreeInput struct {
	EstateId string
	XAxis    int
	YAxis    int
	Height   int
}

type CreateTreeOutput struct {
	Id string
}

type CountTreeByEstateIdInput struct {
	EstateId string
}

type CountTreeByEstateIdOutput struct {
	TotalTrees int
}

type GetTreeHeightStatsInput struct {
	EstateId string
}

type GetTreeHeightStatsOutput struct {
	MaxHeight    int
	MinHeight    int
	MedianHeight int
}

type GetTreesByEstateIdInput struct {
	EstateId string
}

type GetTreesByEstateIdOutput struct {
	Id       string
	EstateId string
	XAxis    int
	YAxis    int
	Height   int
}

func (gebio *GetEstateByIdOutput) IsCordinateOutOfBound(x, y int) bool {
	return x > gebio.Length || y > gebio.Width
}
