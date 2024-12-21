Comparing bitsets with array lookups for character comparing.

To run the benchmark:
```bash
go test -bench=.
```

to view the assembly you can run:

```bash
go test -gcflags=-S is_alpha_numeric_test.go 2> asm.s
```

The assembly of the bitvector:

```asm
command-line-arguments.AsciiVector.Get STEXT nosplit size=21 args=0x28 locals=0x0 funcid=0x0 align=0x0
	0x0000 00000 (is_alpha_numeric_test.go:20)	TEXT	command-line-arguments.AsciiVector.Get(SB), NOSPLIT|NOFRAME|ABIInternal, $0-40
	0x0000 00000 (is_alpha_numeric_test.go:20)	FUNCDATA	$0, gclocals·g2BeySu+wFnoycgXfElmcg==(SB)
	0x0000 00000 (is_alpha_numeric_test.go:20)	FUNCDATA	$1, gclocals·g2BeySu+wFnoycgXfElmcg==(SB)
	0x0000 00000 (is_alpha_numeric_test.go:20)	FUNCDATA	$5, command-line-arguments.AsciiVector.Get.arginfo1(SB)
	0x0000 00000 (is_alpha_numeric_test.go:20)	FUNCDATA	$6, command-line-arguments.AsciiVector.Get.argliveinfo(SB)
	0x0000 00000 (is_alpha_numeric_test.go:20)	PCDATA	$3, $1
	0x0000 00000 (is_alpha_numeric_test.go:21)	MOVL	AX, CX
	0x0002 00002 (is_alpha_numeric_test.go:21)	SHRB	$6, AL
	0x0005 00005 (is_alpha_numeric_test.go:23)	MOVBLZX	AL, DX
	0x0008 00008 (is_alpha_numeric_test.go:23)	MOVQ	command-line-arguments.bv+8(SP)(DX*8), DX
	0x000d 00013 (is_alpha_numeric_test.go:23)	BTQ	CX, DX
	0x0011 00017 (is_alpha_numeric_test.go:23)	SETCS	AL
	0x0014 00020 (is_alpha_numeric_test.go:23)	RET
```