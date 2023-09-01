package jumper

import (
	"fmt"
	"strconv"
)

type Jumper struct {
	Request
	Response
}

type Number float64

func (n Number) String() string {
	return fmt.Sprintf("%.0f", n)
}

func (n Number) Float64() float64 {
	return float64(n)
}

func (n Number) Int64() int64 {
	i, _ := strconv.ParseInt(n.String(), 10, 64)
	return i
}
