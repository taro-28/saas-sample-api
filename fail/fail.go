package fail

import (
	"github.com/morikuni/failure"
)

const (
	BadRequest          failure.StringCode = "BadRequest"
	NotFound            failure.StringCode = "NotFound"
	InternalServerError failure.StringCode = "InternalServerError"
)
