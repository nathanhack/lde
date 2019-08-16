package vec

import (
	"github.com/nathanhack/lde/internal"
	"math/big"
	"strconv"
	"testing"
)

func TestVectors(t *testing.T) {
	tests := []struct {
		input, expected *Vec
		expectedEqual   bool
	}{
		{Zeros(2), NewVec(big.NewInt(0), big.NewInt(0)), true},
		{Ones(2), NewVec(big.NewInt(1), big.NewInt(1)), true},
		{Ones(2).Add(Ones(1)), NewVec(big.NewInt(2), big.NewInt(1)), true},
		{Ones(2).Scalar(big.NewInt(2)), NewVec(big.NewInt(2), big.NewInt(2)), true},
		{Eigen(1), NewVec(big.NewInt(0), big.NewInt(1), big.NewInt(0)), true},
		{Zeros(2).Set(1, internal.One), NewVec(big.NewInt(0), big.NewInt(1)), true},
		{Zeros(2), Ones(2), false},
	}
	for i, test := range tests {
		t.Run(strconv.FormatInt(int64(i), 10), func(t *testing.T) {
			if test.expectedEqual {
				if !test.input.Equals(test.expected) {
					t.Errorf("expected %v but found %v", test.expected, test.input)
				}
			} else {
				if test.input.Equals(test.expected) {
					t.Errorf("expected equallity to fail %v == %v", test.expected, test.input)
				}
			}
		})
	}
}

func TestValues(t *testing.T) {
	tests := []struct {
		input, expected *big.Int
	}{
		{Zeros(2).Get(0), big.NewInt(0)},
		{Ones(2).Get(0), big.NewInt(1)},
		{Ones(2).Dot(Ones(1)), big.NewInt(1)},
		{Ones(2).Dot(Ones(3)), big.NewInt(2)},
	}
	for i, test := range tests {
		t.Run(strconv.FormatInt(int64(i), 10), func(t *testing.T) {
			if test.input.Cmp(test.expected) != 0 {
				t.Errorf("expected %v but found %v", test.expected, test.input)
			}
		})
	}
}

func TestVec_String(t *testing.T) {
	expected := "{1 1}"
	actual := Ones(2).String()
	if actual != expected {
		t.Errorf("expected %v but found %v", expected, actual)
	}
}
