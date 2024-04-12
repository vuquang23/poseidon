package validator

import "github.com/vuquang23/poseidon/internal/pkg/api/dto"

func ValidateGetTxsReq(req *dto.GetTxsReq) error {
	if !ethAddressRegex.MatchString(req.PoolAddress) {
		return NewValidationError("poolAddress", "invalid")
	}

	return nil
}
