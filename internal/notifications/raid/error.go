package raid

import "errors"

var (
	ErrorEmptyDirectory = errors.New("error empty directory")
	ErrorToLowViewers   = errors.New("error low count viewers")
)
