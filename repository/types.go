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

func (gebio *GetEstateByIdOutput) IsCordinateOutOfBound(x, y int) bool {
	return x > gebio.Length || y > gebio.Width
}
