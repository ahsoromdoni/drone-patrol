package validation

import (
	"github.com/ahsoromdoni/drone-patrol/generated"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

func ValidateCreateTree(req generated.TreeRequest) error {
	return validation.ValidateStruct(&req,
		validation.Field(&req.X,
			validation.Required.Error("X is required"),
			validation.Min(1).Error("X must be greater than 0"),
		),
		validation.Field(&req.Y,
			validation.Required.Error("Y is required"),
			validation.Min(1).Error("Y must be greater than 0"),
		),
		validation.Field(&req.Height,
			validation.Required.Error("Height is required"),
			validation.Min(1).Error("Height must be at least 1"),
			validation.Max(30).Error("Height must not exceed 30"),
		),
	)
}
