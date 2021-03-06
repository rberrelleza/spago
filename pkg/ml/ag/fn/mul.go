// Copyright 2019 spaGO Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package fn

import (
	"github.com/nlpodyssey/spago/pkg/mat"
	"sync"
)

var _ Function = &Mul{}

// Mul is an operator to perform matrix-vector multiplication.
type Mul struct {
	x1 Operand // matrix
	x2 Operand // vector
}

// NewMul returns a new Mul Function.
func NewMul(x1, x2 Operand) *Mul {
	return &Mul{x1: x1, x2: x2}
}

// Forward computes the output of the function.
func (r *Mul) Forward() mat.Matrix {
	if r.x1.Value().Columns() != r.x2.Value().Rows() {
		panic("fn: matrices with not compatible size")
	}
	return r.x1.Value().Mul(r.x2.Value())
}

// TODO: backward of sparse gradients
func (r *Mul) Backward(gy mat.Matrix) {
	if !(r.x1.Value().Rows() == gy.Rows() && r.x2.Value().Columns() == gy.Columns()) {
		panic("fn: matrices with not compatible size")
	}
	var wg sync.WaitGroup
	if r.x1.RequiresGrad() {
		wg.Add(1)
		go func() {
			defer wg.Done()
			x2t := r.x2.Value().T()
			defer mat.ReleaseDense(x2t.(*mat.Dense))
			gx := gy.Mul(x2t)
			defer mat.ReleaseDense(gx.(*mat.Dense))
			r.x1.PropagateGrad(gx)
		}()
	}
	if r.x2.RequiresGrad() {
		wg.Add(1)
		go func() {
			defer wg.Done()
			//r.x2.PropagateGrad(gy.T().Mul(r.x1).T()) // alternative method
			if gy.Columns() == 1 {
				gx := r.x1.Value().(*mat.Dense).MulT(gy)
				defer mat.ReleaseDense(gx.(*mat.Dense))
				r.x2.PropagateGrad(gx)
			} else {
				x1t := r.x1.Value().T()
				defer mat.ReleaseDense(x1t.(*mat.Dense))
				gx := x1t.Mul(gy)
				defer mat.ReleaseDense(gx.(*mat.Dense))
				r.x2.PropagateGrad(gx)
			}
		}()
	}
	wg.Wait()
}
