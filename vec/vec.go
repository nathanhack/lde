package vec

import (
	"github.com/nathanhack/lde/internal"
	"math/big"
	"strings"
)

type Vec struct {
	v []*big.Int
}

func NewVec(n ...*big.Int) *Vec {
	return &Vec{v: n}
}

func NewVecInt64(n ...int64) *Vec {
	x := make([]*big.Int, len(n))
	for i := 0; i < len(n); i++ {
		x[i] = big.NewInt(n[i])
	}
	return &Vec{v: x}
}

func Zeros(n uint) *Vec {
	t := make([]*big.Int, n)
	for i := uint(0); i < n; i++ {
		t[i] = internal.Zero
	}
	return &Vec{v: t}
}

func Ones(n uint) *Vec {
	t := make([]*big.Int, n)
	for i := uint(0); i < n; i++ {
		t[i] = internal.One
	}
	return &Vec{v: t}
}

//Eigen creates a zero index Eigenvector
func Eigen(j uint) *Vec {
	return Zeros(j+1).Set(j, internal.One)
}

func (v *Vec) Len() uint {
	return uint(len(v.v))
}

//Set creates a new vector with the ith indexed value to the value passed in
func (v *Vec) Set(i uint, value *big.Int) *Vec {
	t := make([]*big.Int, internal.Max(i+1, v.Len()))
	copy(t, v.v)
	t[i] = new(big.Int).Set(value)
	return &Vec{v: t}
}

//Get returns a copy of the ith index value from the vector
func (v *Vec) Get(i uint) *big.Int {
	if i < uint(len(v.v)) {
		return new(big.Int).Set(v.v[i])
	} else {
		return internal.Zero
	}
}

//Dot returns the dot product of two vectors. The do not have to be the same size.
func (v *Vec) Dot(v1 *Vec) *big.Int {
	t := new(big.Int).Set(internal.Zero)
	for i := uint(0); i < internal.Min(v.Len(), v1.Len()); i++ {
		t = t.Add(t, new(big.Int).Mul(v.v[i], v1.v[i]))
	}
	return t
}

func (v *Vec) Add(vec *Vec) *Vec {
	t := make([]*big.Int, internal.Max(v.Len(), vec.Len()))
	for i := uint(0); i < internal.Max(v.Len(), vec.Len()); i++ {
		t[i] = new(big.Int).Add(v.Get(i), vec.Get(i))
	}
	return &Vec{v: t}
}

func (v *Vec) Sub(vec *Vec) *Vec {
	t := make([]*big.Int, internal.Max(v.Len(), vec.Len()))
	for i := uint(0); i < internal.Max(v.Len(), vec.Len()); i++ {
		t[i] = new(big.Int).Sub(v.Get(i), vec.Get(i))
	}
	return &Vec{v: t}
}

//Scalar returns a new *Vec with the values equal to the elementwise multiplication with value
func (v *Vec) Scalar(value *big.Int) *Vec {
	t := make([]*big.Int, len(v.v))
	for i := 0; i < len(v.v); i++ {
		t[i] = new(big.Int).Mul(v.v[i], value)
	}
	return &Vec{v: t}
}

//Slice returns a copy of this *Vec values from startIndex to endIndex.
// Always returns a new *Vec of len() = endIndex-startIndex. Any index not
// found in v will be set to zero.
func (v *Vec) Slice(start, end uint) *Vec {
	s := internal.Min(start, end)
	e := internal.Max(start, end)
	size := e - s
	t := make([]*big.Int, 0, size)
	for i := s; i < e; i++ {
		t = append(t, v.Get(i))
	}
	return &Vec{v: t}
}

//Equals will compare the values of two *Vec's with or with out the same length. Note if not equal values must be zero to be equal.
func (v *Vec) Equals(vec *Vec) bool {
	for i := uint(0); i < internal.Max(uint(len(v.v)), uint(len(vec.v))); i++ {
		if v.Get(i).Cmp(vec.Get(i)) != 0 {
			return false
		}
	}
	return true
}

func (v *Vec) Cmp(vec *Vec) int {
	//since we almost don't care about how long
	// they are we'll only compare values

	for i := uint(0); i < internal.Max(v.Len(), vec.Len()); i++ {
		x := v.Get(i).Cmp(vec.Get(i))
		if x != 0 {
			return x
		}
	}
	return 0
}

func (v Vec) String() string {
	sb := strings.Builder{}
	sb.WriteString("{")
	for i := 0; i < len(v.v); i++ {
		sb.WriteString(v.v[i].String())
		if i != len(v.v)-1 {
			sb.WriteString(" ")
		}
	}
	sb.WriteString("}")
	return sb.String()
}
