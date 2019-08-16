package mat

import (
	"fmt"
	"github.com/nathanhack/lde/internal"
	"github.com/nathanhack/lde/vec"
	"math/big"
	"strings"
)

type Mat struct {
	rows uint
	cols uint
	m    []*big.Int // [rows * col] row based
}

func NewMat(rows, cols uint, v ...*big.Int) *Mat {
	if rows*cols != uint(len(v)) {
		panic("shape mismatch")
	}
	return &Mat{
		rows: rows,
		cols: cols,
		m:    v,
	}
}

func NewMatRows(vectors ...*vec.Vec) *Mat {
	m := &Mat{}

	for r, vec := range vectors {
		for c := uint(0); c < vec.Len(); c++ {
			m = m.Set(uint(r), c, vec.Get(c))
		}
	}
	return m
}

func NewMatCols(vectors ...*vec.Vec) *Mat {
	m := &Mat{}

	for c, vec := range vectors {
		for r := uint(0); r < vec.Len(); r++ {
			//we swap because we're doing s
			m = m.Set(r, uint(c), vec.Get(r))
		}
	}
	return m
}

//Set returns a copy of this Matrix with the element changed. row and col are zero indexed
func (m *Mat) Set(row, col uint, value *big.Int) *Mat {
	rows := internal.Max(m.rows, row+1)
	cols := internal.Max(m.cols, col+1)

	t := make([]*big.Int, rows*cols)
	for r := uint(0); r < rows; r++ {
		for c := uint(0); c < cols; c++ {
			t[cols*r+c] = m.Get(r, c)
		}
	}

	t[cols*row+col] = value
	return &Mat{
		rows: rows,
		cols: cols,
		m:    t,
	}

}

//Get returns the value from the matrix if outside the matrix size returns zero
func (m *Mat) Get(row, col uint) *big.Int {
	if row < m.rows && col < m.cols {
		return new(big.Int).Set(m.m[m.cols*row+col])
	}
	return internal.Zero
}

func (m *Mat) Shape() (rows, cols uint) {
	return m.rows, m.cols
}

//T is the returns a copy of this that is transposed
func (m *Mat) T() *Mat {
	t := &Mat{}
	for r := uint(0); r < m.rows; r++ {
		for c := uint(0); c < m.cols; c++ {
			t = t.Set(c, r, m.Get(r, c))
		}
	}
	return t
}

func (m *Mat) Equals(mat *Mat) bool {
	if m.cols != mat.cols || m.rows != mat.rows {
		return false
	}

	for i := 0; i < len(m.m); i++ {
		if m.m[i].Cmp(mat.m[i]) != 0 {
			return false
		}
	}
	return true
}

func (m Mat) String() string {
	sb := strings.Builder{}
	sb.WriteString("[")
	for r := uint(0); r < m.rows; r++ {
		sb.WriteString("[")
		for c := uint(0); c < m.cols; c++ {
			sb.WriteString(m.m[m.cols*r+c].String())
			if c != m.cols-1 {
				sb.WriteString(" ")
			}
		}
		sb.WriteString("]")
	}
	sb.WriteString("]")
	return sb.String()
}

func (m *Mat) Add(mat *Mat) *Mat {
	if m.cols != mat.cols || m.rows != mat.rows {
		panic("mat's must have the same shape")
	}

	t := make([]*big.Int, len(m.m))

	for i := 0; i < len(m.m); i++ {
		t[i] = new(big.Int).Add(m.m[i], mat.m[i])
	}

	return &Mat{
		rows: m.rows,
		cols: m.cols,
		m:    t,
	}
}

func (m *Mat) GetCol(c uint) *vec.Vec {
	if c >= m.cols {
		return vec.Zeros(m.rows)
	}

	t := make([]*big.Int, m.rows)

	for r := uint(0); r < m.rows; r++ {
		t[r] = m.Get(r, c)
	}

	return vec.NewVec(t...)
}

func (m *Mat) GetCols() []*vec.Vec {
	toReturn := make([]*vec.Vec, m.cols)
	for i := uint(0); i < m.cols; i++ {
		toReturn[i] = m.GetCol(i)
	}
	return toReturn
}

func (m *Mat) GetRow(r uint) *vec.Vec {
	if r >= m.cols {
		return vec.Zeros(m.cols)
	}

	t := make([]*big.Int, m.cols)

	for c := uint(0); c < m.cols; c++ {
		t[c] = m.Get(r, c)
	}

	return vec.NewVec(t...)
}

func (m *Mat) GetRows() []*vec.Vec {
	toReturn := make([]*vec.Vec, m.rows)
	for i := uint(0); i < m.rows; i++ {
		toReturn[i] = m.GetRow(i)
	}
	return toReturn
}

func (m *Mat) Mul(mat *Mat) *Mat {
	if m.cols != mat.rows {
		panic(fmt.Sprintf("shape mismatch %v x %v != %v x %v", m.rows, m.cols, mat.rows, mat.cols))
	}
	t := &Mat{}
	for r := uint(0); r < m.rows; r++ {
		rvec := m.GetRow(r)
		for c := uint(0); c < mat.cols; c++ {
			cvec := mat.GetCol(c)
			t = t.Set(r, c, rvec.Dot(cvec))
		}
	}
	return t
}
