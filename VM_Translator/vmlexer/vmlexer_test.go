package vmlexer

import (
	"fmt"
	"strings"
	"testing"

	"github.com/quickwritereader/vmTranslator/vmtoken"
)

func TestNextToken(t *testing.T) {

	var test_code = `
	// This file is part of www.nand2tetris.org
	// and the book "The Elements of Computing Systems"
	// by Nisan and Schocken, MIT Press.
	// File name: projects/07/MemoryAccess/BasicTest/BasicTest.vm

	// Executes pop and push commands using the virtual memory segments.
	push constant 10
	pop local 0
	push constant 21
	push constant 22
	pop argument 2
	pop argument 1
	push constant 36
	pop this 6
	push constant 42
	push constant 45
	pop that 5
	pop that 2
	push constant 510
	pop temp 6
	push local 0
	push that 5
	add
	push argument 1
	sub
	push this 6
	push this 6
	add
	sub
	push temp 6
	add
	`

	var tokList = [...]vmtoken.Token{
		{Type: vmtoken.COMMENT, Literal: "// This file is part of www.nand2tetris.org"},
		{Type: vmtoken.COMMENT, Literal: "// and the book \"The Elements of Computing Systems\""},
		{Type: vmtoken.COMMENT, Literal: "// by Nisan and Schocken, MIT Press."},
		{Type: vmtoken.COMMENT, Literal: "// File name: projects/07/MemoryAccess/BasicTest/BasicTest.vm"},
		{Type: vmtoken.COMMENT, Literal: "// Executes pop and push commands using the virtual memory segments."},
		{Type: vmtoken.PUSH, Literal: "push"},
		{Type: vmtoken.CONSTANT, Literal: "constant"},
		{Type: vmtoken.NUMBER, Literal: "10"},
		{Type: vmtoken.POP, Literal: "pop"},
		{Type: vmtoken.LOCAL, Literal: "local"},
		{Type: vmtoken.NUMBER, Literal: "0"},
		{Type: vmtoken.PUSH, Literal: "push"},
		{Type: vmtoken.CONSTANT, Literal: "constant"},
		{Type: vmtoken.NUMBER, Literal: "21"},
		{Type: vmtoken.PUSH, Literal: "push"},
		{Type: vmtoken.CONSTANT, Literal: "constant"},
		{Type: vmtoken.NUMBER, Literal: "22"},
		{Type: vmtoken.POP, Literal: "pop"},
		{Type: vmtoken.ARGUMENT, Literal: "argument"},
		{Type: vmtoken.NUMBER, Literal: "2"},
		{Type: vmtoken.POP, Literal: "pop"},
		{Type: vmtoken.ARGUMENT, Literal: "argument"},
		{Type: vmtoken.NUMBER, Literal: "1"},
		{Type: vmtoken.PUSH, Literal: "push"},
		{Type: vmtoken.CONSTANT, Literal: "constant"},
		{Type: vmtoken.NUMBER, Literal: "36"},
		{Type: vmtoken.POP, Literal: "pop"},
		{Type: vmtoken.THIS, Literal: "this"},
		{Type: vmtoken.NUMBER, Literal: "6"},
		{Type: vmtoken.PUSH, Literal: "push"},
		{Type: vmtoken.CONSTANT, Literal: "constant"},
		{Type: vmtoken.NUMBER, Literal: "42"},
		{Type: vmtoken.PUSH, Literal: "push"},
		{Type: vmtoken.CONSTANT, Literal: "constant"},
		{Type: vmtoken.NUMBER, Literal: "45"},
		{Type: vmtoken.POP, Literal: "pop"},
		{Type: vmtoken.THAT, Literal: "that"},
		{Type: vmtoken.NUMBER, Literal: "5"},
		{Type: vmtoken.POP, Literal: "pop"},
		{Type: vmtoken.THAT, Literal: "that"},
		{Type: vmtoken.NUMBER, Literal: "2"},
		{Type: vmtoken.PUSH, Literal: "push"},
		{Type: vmtoken.CONSTANT, Literal: "constant"},
		{Type: vmtoken.NUMBER, Literal: "510"},
		{Type: vmtoken.POP, Literal: "pop"},
		{Type: vmtoken.TEMP, Literal: "temp"},
		{Type: vmtoken.NUMBER, Literal: "6"},
		{Type: vmtoken.PUSH, Literal: "push"},
		{Type: vmtoken.LOCAL, Literal: "local"},
		{Type: vmtoken.NUMBER, Literal: "0"},
		{Type: vmtoken.PUSH, Literal: "push"},
		{Type: vmtoken.THAT, Literal: "that"},
		{Type: vmtoken.NUMBER, Literal: "5"},
		{Type: vmtoken.ARITHMETIC, Literal: "add"},
		{Type: vmtoken.PUSH, Literal: "push"},
		{Type: vmtoken.ARGUMENT, Literal: "argument"},
		{Type: vmtoken.NUMBER, Literal: "1"},
		{Type: vmtoken.ARITHMETIC, Literal: "sub"},
		{Type: vmtoken.PUSH, Literal: "push"},
		{Type: vmtoken.THIS, Literal: "this"},
		{Type: vmtoken.NUMBER, Literal: "6"},
		{Type: vmtoken.PUSH, Literal: "push"},
		{Type: vmtoken.THIS, Literal: "this"},
		{Type: vmtoken.NUMBER, Literal: "6"},
		{Type: vmtoken.ARITHMETIC, Literal: "add"},
		{Type: vmtoken.ARITHMETIC, Literal: "sub"},
		{Type: vmtoken.PUSH, Literal: "push"},
		{Type: vmtoken.TEMP, Literal: "temp"},
		{Type: vmtoken.NUMBER, Literal: "6"},
		{Type: vmtoken.ARITHMETIC, Literal: "add"},
	}
	for _, buffer_size := range [...]int{-1, MIN_VALID_BUFFER + 1, 13, 31, BUFFER_SIZE} {

		myReader1 := strings.NewReader(test_code)
		lexer1 := NewBufferSize(myReader1, buffer_size)
		fmt.Printf("Buffer Size %v\n", buffer_size)
		i := 0
		for {
			tok_1 := lexer1.NextToken()

			if tok_1.Type == vmtoken.EOF || i == len(tokList) {
				break
			}
			tok := tokList[i]
			i += 1

			fmt.Printf("%v::: %v\n", tok_1.Type, tok_1.Literal)
			if tok != tok_1 {
				t.Errorf("%v != %v \n", tok, tok_1)
			}

		}
	}

}

func TestNextToken_2(t *testing.T) {
	mg := strings.NewReader(`// Executes pop and push commands using the static segment.
	push constant 111
	push constant 333
	push constant 888
	pop static 8
	pop static 3
	pop static 1
	push static 3
	push static 1
	sub
	push static 8
	add
	push constant 17
	push constant 17
	eq
	push constant 17
	push constant 16
	eq
	push constant 16
	push constant 17
	eq
	push constant 892
	push constant 891
	lt
	push constant 891
	push constant 892
	lt
	push constant 891
	push constant 891
	lt
	push constant 32767
	push constant 32766
	gt
	push constant 32766
	push constant 32767
	gt
	push constant 32766
	push constant 32766
	gt
	push constant 57
	push constant 31
	push constant 53
	add
	push constant 112
	sub
	neg
	and
	push constant 82
	or
	not
	`)
	toklist := [...]vmtoken.Token{
		{Type: vmtoken.COMMENT, Literal: "// Executes pop and push commands using the static segment."},
		{Type: vmtoken.PUSH, Literal: "push"},
		{Type: vmtoken.CONSTANT, Literal: "constant"},
		{Type: vmtoken.NUMBER, Literal: "111"},
		{Type: vmtoken.PUSH, Literal: "push"},
		{Type: vmtoken.CONSTANT, Literal: "constant"},
		{Type: vmtoken.NUMBER, Literal: "333"},
		{Type: vmtoken.PUSH, Literal: "push"},
		{Type: vmtoken.CONSTANT, Literal: "constant"},
		{Type: vmtoken.NUMBER, Literal: "888"},
		{Type: vmtoken.POP, Literal: "pop"},
		{Type: vmtoken.STATIC, Literal: "static"},
		{Type: vmtoken.NUMBER, Literal: "8"},
		{Type: vmtoken.POP, Literal: "pop"},
		{Type: vmtoken.STATIC, Literal: "static"},
		{Type: vmtoken.NUMBER, Literal: "3"},
		{Type: vmtoken.POP, Literal: "pop"},
		{Type: vmtoken.STATIC, Literal: "static"},
		{Type: vmtoken.NUMBER, Literal: "1"},
		{Type: vmtoken.PUSH, Literal: "push"},
		{Type: vmtoken.STATIC, Literal: "static"},
		{Type: vmtoken.NUMBER, Literal: "3"},
		{Type: vmtoken.PUSH, Literal: "push"},
		{Type: vmtoken.STATIC, Literal: "static"},
		{Type: vmtoken.NUMBER, Literal: "1"},
		{Type: vmtoken.ARITHMETIC, Literal: "sub"},
		{Type: vmtoken.PUSH, Literal: "push"},
		{Type: vmtoken.STATIC, Literal: "static"},
		{Type: vmtoken.NUMBER, Literal: "8"},
		{Type: vmtoken.ARITHMETIC, Literal: "add"},
		{Type: vmtoken.PUSH, Literal: "push"},
		{Type: vmtoken.CONSTANT, Literal: "constant"},
		{Type: vmtoken.NUMBER, Literal: "17"},
		{Type: vmtoken.PUSH, Literal: "push"},
		{Type: vmtoken.CONSTANT, Literal: "constant"},
		{Type: vmtoken.NUMBER, Literal: "17"},
		{Type: vmtoken.ARITHMETIC, Literal: "eq"},
		{Type: vmtoken.PUSH, Literal: "push"},
		{Type: vmtoken.CONSTANT, Literal: "constant"},
		{Type: vmtoken.NUMBER, Literal: "17"},
		{Type: vmtoken.PUSH, Literal: "push"},
		{Type: vmtoken.CONSTANT, Literal: "constant"},
		{Type: vmtoken.NUMBER, Literal: "16"},
		{Type: vmtoken.ARITHMETIC, Literal: "eq"},
		{Type: vmtoken.PUSH, Literal: "push"},
		{Type: vmtoken.CONSTANT, Literal: "constant"},
		{Type: vmtoken.NUMBER, Literal: "16"},
		{Type: vmtoken.PUSH, Literal: "push"},
		{Type: vmtoken.CONSTANT, Literal: "constant"},
		{Type: vmtoken.NUMBER, Literal: "17"},
		{Type: vmtoken.ARITHMETIC, Literal: "eq"},
		{Type: vmtoken.PUSH, Literal: "push"},
		{Type: vmtoken.CONSTANT, Literal: "constant"},
		{Type: vmtoken.NUMBER, Literal: "892"},
		{Type: vmtoken.PUSH, Literal: "push"},
		{Type: vmtoken.CONSTANT, Literal: "constant"},
		{Type: vmtoken.NUMBER, Literal: "891"},
		{Type: vmtoken.ARITHMETIC, Literal: "lt"},
		{Type: vmtoken.PUSH, Literal: "push"},
		{Type: vmtoken.CONSTANT, Literal: "constant"},
		{Type: vmtoken.NUMBER, Literal: "891"},
		{Type: vmtoken.PUSH, Literal: "push"},
		{Type: vmtoken.CONSTANT, Literal: "constant"},
		{Type: vmtoken.NUMBER, Literal: "892"},
		{Type: vmtoken.ARITHMETIC, Literal: "lt"},
		{Type: vmtoken.PUSH, Literal: "push"},
		{Type: vmtoken.CONSTANT, Literal: "constant"},
		{Type: vmtoken.NUMBER, Literal: "891"},
		{Type: vmtoken.PUSH, Literal: "push"},
		{Type: vmtoken.CONSTANT, Literal: "constant"},
		{Type: vmtoken.NUMBER, Literal: "891"},
		{Type: vmtoken.ARITHMETIC, Literal: "lt"},
		{Type: vmtoken.PUSH, Literal: "push"},
		{Type: vmtoken.CONSTANT, Literal: "constant"},
		{Type: vmtoken.NUMBER, Literal: "32767"},
		{Type: vmtoken.PUSH, Literal: "push"},
		{Type: vmtoken.CONSTANT, Literal: "constant"},
		{Type: vmtoken.NUMBER, Literal: "32766"},
		{Type: vmtoken.ARITHMETIC, Literal: "gt"},
		{Type: vmtoken.PUSH, Literal: "push"},
		{Type: vmtoken.CONSTANT, Literal: "constant"},
		{Type: vmtoken.NUMBER, Literal: "32766"},
		{Type: vmtoken.PUSH, Literal: "push"},
		{Type: vmtoken.CONSTANT, Literal: "constant"},
		{Type: vmtoken.NUMBER, Literal: "32767"},
		{Type: vmtoken.ARITHMETIC, Literal: "gt"},
		{Type: vmtoken.PUSH, Literal: "push"},
		{Type: vmtoken.CONSTANT, Literal: "constant"},
		{Type: vmtoken.NUMBER, Literal: "32766"},
		{Type: vmtoken.PUSH, Literal: "push"},
		{Type: vmtoken.CONSTANT, Literal: "constant"},
		{Type: vmtoken.NUMBER, Literal: "32766"},
		{Type: vmtoken.ARITHMETIC, Literal: "gt"},
		{Type: vmtoken.PUSH, Literal: "push"},
		{Type: vmtoken.CONSTANT, Literal: "constant"},
		{Type: vmtoken.NUMBER, Literal: "57"},
		{Type: vmtoken.PUSH, Literal: "push"},
		{Type: vmtoken.CONSTANT, Literal: "constant"},
		{Type: vmtoken.NUMBER, Literal: "31"},
		{Type: vmtoken.PUSH, Literal: "push"},
		{Type: vmtoken.CONSTANT, Literal: "constant"},
		{Type: vmtoken.NUMBER, Literal: "53"},
		{Type: vmtoken.ARITHMETIC, Literal: "add"},
		{Type: vmtoken.PUSH, Literal: "push"},
		{Type: vmtoken.CONSTANT, Literal: "constant"},
		{Type: vmtoken.NUMBER, Literal: "112"},
		{Type: vmtoken.ARITHMETIC, Literal: "sub"},
		{Type: vmtoken.ARITHMETIC, Literal: "neg"},
		{Type: vmtoken.ARITHMETIC, Literal: "and"},
		{Type: vmtoken.PUSH, Literal: "push"},
		{Type: vmtoken.CONSTANT, Literal: "constant"},
		{Type: vmtoken.NUMBER, Literal: "82"},
		{Type: vmtoken.ARITHMETIC, Literal: "or"},
		{Type: vmtoken.ARITHMETIC, Literal: "not"},
		{Type: vmtoken.EOF, Literal: ""},
	}
	tt := New(mg)
	for _, tok := range toklist {
		tok_1 := tt.NextToken()
		fmt.Printf("%v::: %v\n", tok_1.Type, tok_1.Literal)
		if tok != tok_1 {
			t.Errorf("%v != %v \n", tok, tok_1)
		}
	}

}

func TestNextToken_BRANCH(t *testing.T) {
	mg := strings.NewReader(`
	  goto LBL
	  label LBL
	  if-goto LBL
	`)
	toklist := [...]vmtoken.Token{
		{Type: vmtoken.GOTO, Literal: "goto"},
		{Type: vmtoken.IDENTIFIER, Literal: "LBL"},
		{Type: vmtoken.LABEL, Literal: "label"},
		{Type: vmtoken.IDENTIFIER, Literal: "LBL"},
		{Type: vmtoken.IF_GOTO, Literal: "if-goto"},
		{Type: vmtoken.IDENTIFIER, Literal: "LBL"},
		{Type: vmtoken.EOF, Literal: ""},
	}
	tt := New(mg)
	for _, tok := range toklist {
		tok_1 := tt.NextToken()
		fmt.Printf("%#v\n", tok_1)
		if tok != tok_1 {
			t.Errorf("%#v != %#v \n", tok, tok_1)
		}
	}

}

func TestNextToken_Xtra(t *testing.T) {
	mg := strings.NewReader(`
	  goto LBL//okky
	  goto//okky LBL//okky
	`)
	toklist := [...]vmtoken.Token{
		{Type: vmtoken.GOTO, Literal: "goto"},
		{Type: vmtoken.IDENTIFIER, Literal: "LBL"},
		{Type: "COMMENT", Literal: "//okky"},
		{Type: "GOTO", Literal: "goto"},
		{Type: "COMMENT", Literal: "//okky LBL//okky"},
		{Type: "EOF", Literal: ""},
	}
	tt := New(mg)
	for _, tok := range toklist {
		tok_1 := tt.NextToken()
		fmt.Printf("%#v\n", tok_1)
		if tok != tok_1 {
			t.Errorf("%#v != %#v \n", tok, tok_1)
		}
	}

}
