package internal

import "math/big"

var One = big.NewInt(1)
var Zero = big.NewInt(0)
var NegOne = new(big.Int).Neg(One)

func Max(a, b uint) uint {
	if a > b {
		return a
	}
	return b
}

func Min(a, b uint) uint {
	if a < b {
		return a
	}
	return b
}
