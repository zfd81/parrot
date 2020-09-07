package core

import "github.com/pkg/errors"

var (
	ErrServExists   = errors.New("The service already exists")
	ErrServNotExist = errors.New("Service does not exist")

	ErrDsExists   = errors.New("The datasource already exists")
	ErrDsNotExist = errors.New("DataSource does not exist")

	ErrParamBad      = errors.New("Bad parameter")
	ErrParamNotFound = errors.New("Parameter not found")
)
