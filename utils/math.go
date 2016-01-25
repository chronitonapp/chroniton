package utils

import (
	"fmt"
)

func RoundFloat(x float64) string {
	return fmt.Sprintf("%0.1f", x)
}
