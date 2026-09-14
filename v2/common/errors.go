package common

import "fmt"

var (
	ErrUnsupportedFlavorType     = fmt.Errorf("unsupported flavor type")
	ErrUnsupportedFlavorDiskType = fmt.Errorf("unsupported flavor disk type")
	ErrUnsupportedNodeGroupRole  = fmt.Errorf("unsupported node group role")
	ErrPositiveIntegerRequired   = fmt.Errorf("value must be greater than 0")
	ErrFieldRequired             = fmt.Errorf("required field")
	ErrFieldEmptySlice           = fmt.Errorf("must be at least one value")
)
