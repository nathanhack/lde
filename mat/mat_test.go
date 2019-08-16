package mat

import (
	"github.com/nathanhack/lde/vec"
	"math/big"
	"strconv"
	"testing"
)

func TestGeneral(t *testing.T) {
	tests := []struct {
		m1, m2 *Mat
		equal  bool
	}{
		{NewMatRows(vec.Ones(2), vec.Zeros(2)), NewMat(2, 2, big.NewInt(1), big.NewInt(1), big.NewInt(0), big.NewInt(0)), true},
		{NewMatRows(vec.Ones(2), vec.Zeros(2)).T(), NewMatCols(vec.Ones(2), vec.Zeros(2)), true},
		{NewMatRows(vec.Ones(2), vec.Ones(2)).Add(NewMatRows(vec.Ones(2), vec.Zeros(2))), NewMat(2, 2, big.NewInt(2), big.NewInt(2), big.NewInt(1), big.NewInt(1)), true},
		{NewMatRows(vec.Ones(2), vec.Ones(2)).Mul(NewMatRows(vec.Ones(2), vec.Ones(2))), NewMat(2, 2, big.NewInt(2), big.NewInt(2), big.NewInt(2), big.NewInt(2)), true},
	}
	for i, test := range tests {
		t.Run(strconv.FormatInt(int64(i), 10), func(t *testing.T) {
			if test.equal {
				if !test.m1.Equals(test.m2) {
					t.Errorf("expected equal %v and %v", test.m1, test.m2)
				}
			} else {
				if test.m1.Equals(test.m2) {
					t.Errorf("expected not equal %v and %v", test.m1, test.m2)
				}
			}

		})
	}
}
