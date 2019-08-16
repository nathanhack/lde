package lde

import (
	"github.com/nathanhack/lde/mat"
	"github.com/nathanhack/lde/vec"
	"math/big"
	"strconv"
	"testing"
)

func TestA(t *testing.T) {
	tests := []struct {
		m        *mat.Mat
		v        *vec.Vec
		expected *vec.Vec
	}{
		{
			mat.NewMatRows(vec.NewVecInt64(-1, 1, 2, -3),
				vec.NewVecInt64(-1, 3, -2, -1)),
			vec.Eigen(0),
			vec.NewVecInt64(-1, -1),
		},

		{
			mat.NewMatRows(vec.NewVecInt64(-1, 1, 2, -3),
				vec.NewVecInt64(-1, 3, -2, -1)),
			vec.Eigen(2),
			vec.NewVecInt64(2, -2),
		},
	}
	for i, test := range tests {
		t.Run(strconv.FormatInt(int64(i), 10), func(t *testing.T) {
			actual := a(test.m, test.v)

			if !actual.Equals(test.expected) {
				t.Errorf("expected %v but found %v", test.expected, actual)
			}
		})
	}
}

func TestContainment(t *testing.T) {
	tests := []struct {
		check               *vec.Vec
		bases               []*vec.Vec
		expectedContainment bool
	}{
		{vec.NewVecInt64(1, 1, 1, 1), []*vec.Vec{vec.NewVecInt64(0, 1, 1, 1)}, true},
		{vec.NewVecInt64(4, 2, 1, 0), []*vec.Vec{vec.NewVecInt64(0, 1, 1, 1)}, false},
	}
	for i, test := range tests {
		t.Run(strconv.FormatInt(int64(i), 10), func(t *testing.T) {
			a := containedInMinimalSet(test.check, test.bases)
			if test.expectedContainment {
				if a != test.expectedContainment {
					t.Errorf("expected containment for %v in %v", test.check, test.bases)
				}
			} else {
				if a != test.expectedContainment {
					t.Errorf("expected %v not to be contained by %v", test.check, test.bases)
				}
			}
		})
	}
}

func TestHomogeneous(t *testing.T) {
	tests := []struct {
		m        *mat.Mat
		expected []*vec.Vec
	}{
		{
			mat.NewMatRows(vec.NewVecInt64(-1, 1, 2, -3),
				vec.NewVecInt64(-1, 3, -2, -1)),
			[]*vec.Vec{vec.NewVecInt64(0, 1, 1, 1), vec.NewVecInt64(4, 2, 1, 0)},
		},
		{
			mat.NewMatRows(vec.NewVecInt64(6, -9, 2)),
			[]*vec.Vec{vec.NewVecInt64(0, 2, 9), vec.NewVecInt64(1, 2, 6), vec.NewVecInt64(2, 2, 3), vec.NewVecInt64(3, 2, 0)},
		},
	}
	for i, test := range tests {
		t.Run(strconv.FormatInt(int64(i), 10), func(t *testing.T) {
			actual := Homogeneous(test.m)
			if len(actual) != len(test.expected) {
				t.Errorf("expected %v but found %v", test.expected, actual)
				return
			}

			for i, x := range test.expected {
				if actual[i].Cmp(x) != 0 {
					t.Errorf("expected %v but found %v", x, actual[i])
				}
			}
		})
	}
}

func TestNonHomogeneous(t *testing.T) {
	tests := []struct {
		a          *mat.Mat
		b          *vec.Vec
		limits     LimitBy
		expectedM1 []*vec.Vec
		expectedM0 []*vec.Vec
	}{
		{
			mat.NewMatRows(vec.NewVecInt64(3, 9, 5)),
			vec.NewVecInt64(20),
			nil,
			[]*vec.Vec{vec.NewVecInt64(0, 0, 4), vec.NewVecInt64(2, 1, 1), vec.NewVecInt64(5, 0, 1)},
			[]*vec.Vec{},
		},

		{
			mat.NewMatRows(vec.NewVecInt64(6, -9, 2)),
			vec.NewVecInt64(0),
			nil,
			[]*vec.Vec{},
			[]*vec.Vec{vec.NewVecInt64(0, 2, 9), vec.NewVecInt64(1, 2, 6), vec.NewVecInt64(2, 2, 3), vec.NewVecInt64(3, 2, 0)},
		},

		{
			mat.NewMatRows(vec.NewVecInt64(6, -9, 2)),
			vec.NewVecInt64(0),
			NewMaxXLimit(big.NewInt(4)),
			[]*vec.Vec{},
			[]*vec.Vec{vec.NewVecInt64(2, 2, 3), vec.NewVecInt64(3, 2, 0)},
		},
	}
	for i, test := range tests {
		t.Run(strconv.FormatInt(int64(i), 10), func(t *testing.T) {
			var actualM1, actualM0 []*vec.Vec
			if test.limits != nil {
				actualM1, actualM0 = NonHomogeneous(test.a, test.b, test.limits)
			} else {
				actualM1, actualM0 = NonHomogeneous(test.a, test.b)
			}
			if len(actualM1) != len(test.expectedM1) {
				t.Errorf("expected %v but found %v", test.expectedM1, actualM1)
			} else {
				for i, x := range test.expectedM1 {
					if actualM1[i].Cmp(x) != 0 {
						t.Errorf("expected %v but found %v", x, actualM1[i])
					}
				}
			}

			if len(actualM0) != len(test.expectedM0) {
				t.Errorf("expected %v but found %v", test.expectedM0, actualM0)

			} else {
				for i, x := range test.expectedM0 {
					if actualM0[i].Cmp(x) != 0 {
						t.Errorf("expected %v but found %v", x, actualM0[i])
					}
				}
			}
		})
	}
}
