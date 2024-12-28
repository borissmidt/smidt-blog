package main

import (
	"github.com/bits-and-blooms/bitset"
	"math/rand"
	"strings"
	"testing"
)

type AsciiVector [32]uint8

const mask = 8 - 1

// Set ith bit
func (bv *AsciiVector) Set(i uint8) {
	idx := i & mask
	bit := i << 4
	bv[idx] |= bit
}

func (bv AsciiVector) Get(i uint8) bool {
	idx := i & mask
	bit := i >> 4
	return (bv[idx] & (1 << bit)) != 0
}

//go:generate go run asm/simd_bitset_asm.go -out simd_bitset_asm.go.s
//func IsSpace2(s []byte, bitset *[32]byte) uint8

const alphaNumericCharacters = "1234567890abcdefghijklmnopqrtsuvwABCDEFGHIJKLMOPQRTUVWXYZ"

var randomBytes = [1024]byte{}

var isDigit2 = [256]bool{}

var bs = bitset.New(256)
var v = AsciiVector{}

func init() {
	for _, i := range alphaNumericCharacters {
		isDigit2[i] = true
		bs.Set(uint(i))
		v.Set(uint8(i))
	}

	for i := 0; i < len(randomBytes); i++ {
		randomBytes[i] = byte(rand.Uint32())
		if randomBytes[i] > 127 {
			randomBytes[i] -= 127
		}
	}
}

func TestBitvector(t *testing.T) {
	vect := AsciiVector{}
	vect.Set(10)
	if !vect.Get(10) {
		t.Error("vect 10 is not set")
	}
	if vect.Get(9) {
		t.Error("vect 10 is not set but get return true")
	}

	for _, c := range alphaNumericCharacters {
		if !v.Get(uint8(c)) {
			t.Error("vect alphaNumericCharacter is not set %i", c)
		}
	}
}

// added to avoid the go compiler optimizing away the loop
var result = false

func BenchmarkArrayIndex(b *testing.B) {
	for i := 0; i < b.N; i++ {
		for i := 0; i < len(randomBytes); i++ {
			for i, x := range randomBytes[i:] {
				if strings.IndexByte(alphaNumericCharacters, x) >= 0 {
					result = i != 0
					break
				}
			}
		}
	}
}

func BenchmarkBranch(b *testing.B) {
	for i := 0; i < b.N; i++ {
		for i := 0; i < len(randomBytes); i++ {
			for i, x := range randomBytes[i:] {
				if '0' <= x && x <= '9' || 'a' <= x && x <= 'z' || 'A' <= x && x <= 'Z' {
					result = i != 0
					break
				}
			}
		}
	}
}

func BenchmarkSwitch(b *testing.B) {
	for i := 0; i < b.N; i++ {
		for i := 0; i < len(randomBytes); i++ {
		loop:
			for i, x := range randomBytes[i:] {
				switch x {
				case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
					fallthrough
				case 'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h', 'i', 'j', 'k', 'l', 'm', 'n', 'o', 'p', 'q', 'r', 's', 't', 'u', 'v', 'w', 'y', 'z':
					fallthrough
				case 'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'Y', 'Z':
					result = i != 0
					break loop
				}
			}
		}
	}
}

func BenchmarkAsciiVector(b *testing.B) {
	for i := 0; i < b.N; i++ {
		for i := 0; i < len(randomBytes); i++ {
			for _, x := range randomBytes[i:] {
				if v.Get(x) {
					result = true
					break
				}
			}
		}
	}
}

func BenchmarkWillfLibBitvector(b *testing.B) {
	for i := 0; i < b.N; i++ {
		for i := 0; i < len(randomBytes); i++ {
			for i, x := range randomBytes[i:] {
				if bs.Test(uint(x)) {
					result = i != 0
					break
				}
			}

		}
	}
}

func BenchmarkLookup(b *testing.B) {
	for i := 0; i < b.N; i++ {
		for i := 0; i < len(randomBytes); i++ {
			for i, x := range randomBytes[i:] {
				if isDigit2[x] {
					result = i != 0
					break
				}
			}
		}
	}
}

//func BenchmarkIsSpace(b *testing.B) {
//	vect := [32]uint8(v)
//
//	for i := 0; i < b.N; i++ {
//		for i := 0; i < len(randomBytes); i++ {
//			x := IsSpace2(randomBytes[i:], &vect)
//			if x != 32 {
//				println(string(randomBytes[i : i+31]))
//			}
//
//		}
//	}
//}
