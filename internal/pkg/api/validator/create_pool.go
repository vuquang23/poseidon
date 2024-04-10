package validator

import "github.com/vuquang23/poseidon/internal/pkg/api/dto"

func ValidateCreatePoolReq(req *dto.CreatePoolReq) error {
	if !ethAddressRegex.MatchString(req.Address) {
		return NewValidationError("address", "invalid")
	}

	if !ethAddressRegex.MatchString(req.Token0) {
		return NewValidationError("token0", "invalid")
	}

	if !ethAddressRegex.MatchString(req.Token1) {
		return NewValidationError("token1", "invalid")
	}

	return nil
}
