// +build ignore
//go:generate go run asm.go -out add_in_place.s -stubs stub.go

package main

import (
	. "github.com/mmcloughlin/avo/build"
	. "github.com/mmcloughlin/avo/operand"
	. "github.com/mmcloughlin/avo/reg"
)

const unroll = 2
const blockItems = 4 * unroll
const blockSize = 8 * blockItems

func main() {
	TEXT("AddInPlace", NOSPLIT, "func(x, y []float64)")
	Doc("Computes y = x + y")

	x := Mem{Base: Load(Param("x").Base(), GP64())}
	y := Mem{Base: Load(Param("y").Base(), GP64())}
	n := Load(Param("y").Len(), GP64())

	xVectors := make([]VecVirtual, unroll)
	yVectors := make([]VecVirtual, unroll)
	for i := range xVectors {
		xVectors[i] = YMM()
		yVectors[i] = YMM()
	}

	Label("block_loop")

	CMPQ(n, U8(blockItems))
	JL(LabelRef("tail_loop"))

	for i, xVec := range xVectors {
		VMOVUPD(x.Offset(i*32), xVec)
	}

	for i, yVec := range yVectors {
		VMOVUPD(y.Offset(i*32), yVec)
	}

	for i, xVec := range xVectors {
		yVec := yVectors[i]
		VADDPD(yVec, xVec, yVec)
	}

	for i, yVec := range yVectors {
		VMOVUPD(yVec, y.Offset(i*32))
	}

	ADDQ(U8(blockSize), x.Base)
	ADDQ(U8(blockSize), y.Base)
	SUBQ(U8(blockItems), n)
	JMP(LabelRef("block_loop"))

	Label("tail_loop")
	xVec := XMM()
	yVec := XMM()

	CMPQ(n, U8(0))
	JE(LabelRef("ret"))

	MOVSD(x, xVec)
	MOVSD(y, yVec)
	ADDSD(xVec, yVec)
	MOVSD(yVec, y)

	ADDQ(U8(8), x.Base)
	ADDQ(U8(8), y.Base)

	DECQ(n)
	JMP(LabelRef("tail_loop"))

	Label("ret")
	RET()

	Generate()
}
