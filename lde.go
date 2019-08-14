package LDE

import (
	"fmt"
	"github.com/nathanhack/LDE/internal"
	"github.com/nathanhack/LDE/mat"
	"github.com/nathanhack/LDE/vec"
	"math/big"
	"sort"
)

type fvec struct {
	v, f *vec.Vec
}

type LimitBy interface {
	Stop(current *vec.Vec) bool
}

type maxX struct {
	v *big.Int
}

func NewMaxXLimit(i *big.Int) LimitBy {
	return &maxX{v: i}
}

func (m *maxX) Stop(current *vec.Vec) bool {
	for i := uint(0); i < current.Len(); i++ {
		if current.Get(i).Cmp(m.v) >= 0 {
			return true
		}
	}
	return false
}

//Solves A𝑥 = 0, returns the minimal bases. Each basis can be added in linear
// combination with other bases to construct new solutions.
func Homogeneous(A *mat.Mat, limits ...LimitBy) []*vec.Vec {
	//first we create a set of basis vectors
	𝓟 := make([]*fvec, 0)
	_, cols := A.Shape()
	eigens := make(map[uint]*vec.Vec, 0)
	for i := uint(0); i < cols; i++ {
		v := vec.Eigen(i).Slice(0, cols)
		𝓟 = append(𝓟, &fvec{
			v: v,
			f: vec.Ones(cols).Sub(vec.Ones(i + 1)),
		})
		eigens[i] = v
	}
	return homogeneous(A, 𝓟, eigens, limits...)
}

func homogeneous(A *mat.Mat, 𝓟 []*fvec, eigens map[uint]*vec.Vec, limits ...LimitBy) []*vec.Vec {
	𝓑 := make([]*vec.Vec, 0)
	𝓑Map := make(map[string]bool)
	_, cols := A.Shape()
	zeroVec := vec.Zeros(cols)
	for len(𝓟) > 0 {
		//fist we 𝓑 := 𝓑 ⋃ {𝑥 ∈ 𝓟 | a(𝑥) = 0}
		// which means we add to 𝓑 any non-dup 𝑥 from 𝓟 that solves the equation

		𝓟Not𝓑 := make([]*fvec, 0)
		for _, v := range 𝓟 {
			if a(A, v.v).Equals(zeroVec) {
				//we only add it if it's new
				s := v.v.String()
				if _, has := 𝓑Map[s]; !has {
					𝓑 = append(𝓑, v.v)
					𝓑Map[s] = true
				}
			} else {
				//we put it in 𝓟Not𝓑 for later use
				𝓟Not𝓑 = append(𝓟Not𝓑, v)
			}
		}

		//next we make 𝓠, 𝓠 := { 𝑥 ∈ 𝓟 \ 𝓑 | ∀s ∈ 𝓑, 𝑥 not(⨠)s}
		// which mean we make 𝓠 with everything in 𝓟 not in 𝓑 that is not contained in 𝓑

		𝓠 := make([]*fvec, 0)
		for _, v := range 𝓟Not𝓑 {
			if !containedInMinimalSet(v.v, 𝓑) {
				𝓠 = append(𝓠, v)
			}
		}

		// last we make 𝓟 for the next round, 𝓟 := {𝑥 + e𝑖| 𝑥 ∈ 𝓠, a(𝑥)⋅a(e𝑖) < 0}
		// which means we take all the values in 𝓠 and make new values for 𝓟
		// take each vec in 𝓠 and each eigen vector and if a(𝑥)⋅a(e𝑖) < 0 is true
		// then add 𝑥 and e𝑖 and put that in 𝓟

		𝓟 = make([]*fvec, 0)
		for _, 𝑥 := range 𝓠 {
			a𝑥 := a(A, 𝑥.v)
			frozen := vec.Zeros(cols)
		eigenLoop:
			for _, e𝑖 := range eigens {
				//if not frozen
				if 𝑥.f.Dot(e𝑖).Cmp(internal.Zero) == 0 {
					if a𝑥.Dot(a(A, e𝑖)).Cmp(internal.Zero) < 0 {
						nv := 𝑥.v.Add(e𝑖)
						for _, l := range limits {
							if l.Stop(nv) {
								continue eigenLoop
							}
						}
						𝓟 = append(𝓟, &fvec{
							v: nv,
							f: 𝑥.f.Add(frozen),
						})
						frozen = frozen.Add(e𝑖)
					}
				}
			}
		}
	}

	//we'll sort 𝓑 before returning it

	sort.Slice(𝓑, func(i, j int) bool {
		return 𝓑[i].Cmp(𝓑[j]) < 0
	})

	return 𝓑
}

//NonHomogeneous solve the A𝑥 = b equation. Returns the set of specific solutions (M1) and the homogeneous bases (M0).
// all solutions can be made by taking one from M1 and adding any number the bases from M0 (aka M1+ M0 + M0+...)
func NonHomogeneous(A *mat.Mat, b *vec.Vec, limits ...LimitBy) (M1 []*vec.Vec, M0 []*vec.Vec) {
	// for this case we create a new matrix
	// with the 0th index column set to b
	// then with that as A' solve A'𝑥 = 0
	// then using only the bases with x0 value
	// equal to 1 or 0 so we freeze them
	_, cols := A.Shape()
	c := make([]*vec.Vec, 0, cols+1)
	//make the b vector the 0th columns
	c = append(c, b.Scalar(internal.NegOne))
	//add the rest
	c = append(c, A.GetCols()...)

	//then solve A𝑥 =0
	newA := mat.NewMatCols(c...)
	_, newCols := newA.Shape()

	//next we create a set of basis vectors
	𝓟 := make([]*fvec, 0)
	eigens := make(map[uint]*vec.Vec, 0)
	bIsZeroVec := b.Cmp(vec.Zeros(b.Len())) == 0
	for i := uint(1); i < newCols; i++ {
		// we freeze the first column
		// with the value set to zero and one
		v := vec.Eigen(i).Slice(0, newCols)
		f := vec.Ones(newCols).Sub(vec.Ones(i + 1))
		//first the zero one
		𝓟 = append(𝓟, &fvec{
			v: v,
			f: f,
		})
		//next the with one if
		if !bIsZeroVec {
			𝓟 = append(𝓟, &fvec{
				v: v.Set(0, internal.One),
				f: f,
			})
		}
		eigens[i] = v
	}

	𝓑 := homogeneous(newA, 𝓟, eigens, limits...)
	//𝓑 will contain both the specific and homogeneous values
	//we'll separate them
	M0 = make([]*vec.Vec, 0)
	M1 = make([]*vec.Vec, 0)

	for _, b := range 𝓑 {
		a := b.Slice(1, b.Len())
		if b.Get(0).Cmp(internal.Zero) == 0 {
			M0 = append(M0, a)
		} else {
			M1 = append(M1, a)
		}
	}
	return M1, M0
}

func a(m *mat.Mat, vec *vec.Vec) *vec.Vec {
	_, cols := m.Shape()
	if vec.Len() > cols {
		panic(fmt.Sprintf("vec must be equal the number of cols in the matrix, expected %v but found %v", cols, vec.Len()))
	}
	m1 := m.Mul(mat.NewMatCols(vec.Slice(0, cols)))
	return m1.GetCol(0)
}

func containedInMinimalSet(check *vec.Vec, currentBases []*vec.Vec) bool {
	// to check if it's contained in the minimal bases
	// we check against all current bases

	// we verify for all basis in currentBases
	// that it's not contained in any of them

	//to be contained it must satisfy the follow:
	// (a1,a2,...,aq) ⨠ (b1,b2,...,bq)  ∀𝑖 ∈ [1..q] | a𝑖 ≥ b𝑖 and ∃𝑖 ∈ [1..q] | a𝑖 ≠ b𝑖
mainLoop:
	for _, basis := range currentBases {
		for i := uint(0); i < internal.Max(check.Len(), basis.Len()); i++ {
			//so first we check ∀𝑖 ∈ [1..q] | a𝑖 ≥ b𝑖
			//we care about the opposite so < instead of ≥
			if check.Get(i).Cmp(basis.Get(i)) < 0 {
				continue mainLoop
			}
		}
		//if here we have met the first part of the containment
		// now we check for ∃𝑖 ∈ [1..q] | a𝑖 ≠ b𝑖
		for i := uint(0); i < internal.Max(check.Len(), basis.Len()); i++ {
			//so first we check ∀𝑖 ∈ [1..q] | a𝑖 ≥ b𝑖
			//we care about the opposite so < instead of ≥
			if check.Get(i).Cmp(basis.Get(i)) > 0 {
				return true
			}
		}
	}
	return false
}
