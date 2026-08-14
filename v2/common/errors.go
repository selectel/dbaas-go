package common

import "fmt"

var (
	ErrUnsupportedFlavorType    = fmt.Errorf("unsupported flavor type")
	ErrUnsupportedNodeGroupRole = fmt.Errorf("unsupported node group role")
)
