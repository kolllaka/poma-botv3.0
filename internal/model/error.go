package model

import (
	"errors"
)

var (
	ErrBadLink = errors.New("bad link")

	ErrorEmptyDirectory                = errors.New("error empty directory")
	ErrorRaidToLowViewers              = errors.New("error on raid low count viewers")
	ErrorSubscribeNotAllowedConditions = errors.New("error subscribe not allowed conditions")
)
