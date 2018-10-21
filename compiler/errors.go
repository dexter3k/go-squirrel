package compiler

import (
	"fmt"
)

var (
	ErrExpectStatementEnd = fmt.Errorf("End of statement expected")
    ErrExpectArgument     = fmt.Errorf("Argument expected after ','")
)
