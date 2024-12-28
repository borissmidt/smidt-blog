//go:build ignore

package main

import (
	. "github.com/mmcloughlin/avo/build"
	. "github.com/mmcloughlin/avo/operand"
	"github.com/mmcloughlin/avo/reg"
	"github.com/mmcloughlin/avo/x86"
)

func ReadConstToYmm(mask Mem, reg reg.VecVirtual) {
	ptr := Mem{Base: GP64()}
	LEAQ(mask, ptr.Base)
	x86.VMOVDQU8(ptr, reg)
}

func YmmMask(name string, value byte) Mem {
	constLowNibbleMask := GLOBL(name, RODATA|NOPTR)
	for i := range 32 {
		DATA(i, U8(value))
	}
	return constLowNibbleMask
}

func LoadArrayToYmm(name string, r reg.VecVirtual) {
	cvs := Mem{Base: Load(Param(name), GP64())}
	x86.VMOVDQU8(cvs.Offset(0), r)
}

func main() {
	constLowNibleMask := YmmMask("mask_low_nibble", 0x0F)
	constHighNibleMask := YmmMask("mask_higher_nibble", 0xF0)

	TEXT("IsSpace2", NOSPLIT, "func(s []byte, bitset *[32]byte ) uint8")
	Doc("check that the character is a space or not using a simd register")

	// the bitvector
	bitset := YMM()
	LoadArrayToYmm("bitset", bitset)

	lowNibbleMask := YMM()
	ReadConstToYmm(constLowNibleMask, lowNibbleMask)

	highNibbleMask := YMM()
	ReadConstToYmm(constHighNibleMask, highNibbleMask)

	// get the string pointer from the stack
	str := Mem{Base: Load(Param("s").Base(), GP64())}
	// get the string length
	strEnd := Load(Param("s").Len(), GP64())
	leftover := GP64()
	MOVD(strEnd, leftover)
	// round to the 256 bit in length
	ANDQ(U8((256/8)-1), leftover)
	SUBQ(leftover, strEnd)
	// calculate the end pointer of the string.
	ADDQ(str.Base, strEnd)

	// the next 32 characters
	chars := YMM()

	// is the amount we need to shift
	shiftValues := YMM()
	indexes := YMM()

	// the start of iteration loop over the string
	Label("blockloop")

	// if i == len(str) break
	CMPQ(str.Base, strEnd)
	JE(LabelRef("end"))

	// copy the next 32 characters to an ymm register
	x86.VMOVDQU8(str.Offset(0), chars)

	// get the index value (c & 0xF0) >> 4
	VANDNPD(shiftValues, highNibbleMask, chars)
	VPSRLQ(U8(4), shiftValues, shiftValues)
	INT(U8(3))
	// get the shift values c & 0x0F
	VANDNPD(indexes, lowNibbleMask, chars)
	VPSHUFB(chars, indexes, chars)
	x := GP32()
	VPMOVMSKB(chars, x)
	TZCNTL(x, x)
	Store(x.As8(), ReturnIndex(0))
	RET()
	CMPL(x, U8(32))
	Store(x.As8(), ReturnIndex(0))
	RET()
	// if set jump to the end
	JNC(LabelRef("end"))

	// otherwise check the next character
	ADDQ(U32(1), str.Base)
	JMP(LabelRef("blockloop"))

	Label("end")

	Store(x.As8(), ReturnIndex(0))
	RET()
	Generate()
}
