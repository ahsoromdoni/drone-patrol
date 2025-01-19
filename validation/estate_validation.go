package validation

import (
	"github.com/ahsoromdoni/drone-patrol/generated"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

func ValidateCreateEstate(req generated.EstateRequest) error {
	return validation.ValidateStruct(&req,
		validation.Field(&req.Length,
			validation.Required.Error("Length is required"),
			validation.Min(1).Error("Length must be greater than 0"),
		),
		validation.Field(&req.Width,
			validation.Required.Error("Width is required"),
			validation.Min(1).Error("Width must be greater than 0"),
		),
	)
}
