//go:build ignore

package main

import (
	. "github.com/mmcloughlin/avo/build"
	. "github.com/mmcloughlin/avo/operand"
)

func main() {
	TEXT("IsSpace2", NOSPLIT, "func(s []byte, bitset *[32]byte) uint8")
	Doc("check that the character is a space or not using a simd register")
	cvs := Mem{Base: Load(Param("bitset"), GP64())}
	// the bitvector
	vector := XMM()

	// load the bitvector form memory
	VMOVDQA(cvs.Offset(0), vector)

	// get the string pointer from the stack
	str := Mem{Base: Load(Param("s").Base(), GP64())}
	// get the string length
	strEnd := Load(Param("s").Len(), GP64())
	x := GP64()

	// mask is an xmm register with value 1
	mask := XMM()

	// is the amount we need to shift
	shiftValue := XMM()
	// zero the mask
	XORPD(mask, mask)

	// load 1 in a GP register
	MOVQ(U32(1), x)
	// load the 1 into the xmm
	MOVQ(x, mask)

	// calculate the end pointer of the string.
	ADDQ(str.Base, strEnd)

	// the start of iteration loop over the string
	Label("blockloop")
	// if i == len(str) break
	CMPQ(str.Base, strEnd)
	JE(LabelRef("end"))

	// copy the character to the mmx register
	MOVD(str.Offset(0), shiftValue)
	// shift the value by x
	VPSLLVD(shiftValue, mask, shiftValue)
	// see if the bit is set
	PTEST(shiftValue, vector)

	// if set jump to the end
	JC(LabelRef("true_end"))

	// otherwise check the next character
	ADDQ(U32(1), str.Base)
	JMP(LabelRef("blockloop"))

	Label("true_end")
	MOVD(U32(1), x)
	Label("end")
	Store(x.As8(), ReturnIndex(0))
	RET()
	Generate()
}
