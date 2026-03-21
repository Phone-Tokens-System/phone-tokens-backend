package billing

import "errors"

var (
	ErrNoPackageUnitsLeft = errors.New("no package units left")
	ErrNotEnoughBalance   = errors.New("not enough balance")
)
