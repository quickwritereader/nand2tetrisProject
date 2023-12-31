package parser

import (
	"reflect"
	"strings"
	"testing"

	"github.com/quickwritereader/vmTranslator/codegen"
)

func TestParser(t *testing.T) {
	mg := strings.NewReader(`
	push constant 3030
	pop pointer 0
	push constant 3040
	pop pointer 1
	push constant 32
	pop this 2
	push constant 46
	pop that 6
	push pointer 0
	push pointer 1
	add
	push this 2
	sub
	push that 6
	add
	// Executes pop and push commands using the static segment.
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

	test_expected := [...]codegen.Command{
		&codegen.CC_PUSH_POP{Pop: false, Segment: "CONSTANT", Index: 3030},
		&codegen.CC_PUSH_POP{Pop: true, Segment: "POINTER", Index: 0},
		&codegen.CC_PUSH_POP{Pop: false, Segment: "CONSTANT", Index: 3040},
		&codegen.CC_PUSH_POP{Pop: true, Segment: "POINTER", Index: 1},
		&codegen.CC_PUSH_POP{Pop: false, Segment: "CONSTANT", Index: 32},
		&codegen.CC_PUSH_POP{Pop: true, Segment: "THIS", Index: 2},
		&codegen.CC_PUSH_POP{Pop: false, Segment: "CONSTANT", Index: 46},
		&codegen.CC_PUSH_POP{Pop: true, Segment: "THAT", Index: 6},
		&codegen.CC_PUSH_POP{Pop: false, Segment: "POINTER", Index: 0},
		&codegen.CC_PUSH_POP{Pop: false, Segment: "POINTER", Index: 1},
		&codegen.CC_ARITHMETIC{Op: "add"},
		&codegen.CC_PUSH_POP{Pop: false, Segment: "THIS", Index: 2},
		&codegen.CC_ARITHMETIC{Op: "sub"},
		&codegen.CC_PUSH_POP{Pop: false, Segment: "THAT", Index: 6},
		&codegen.CC_ARITHMETIC{Op: "add"},
		&codegen.CC_COMMENT{Text: "// Executes pop and push commands using the static segment."},
		&codegen.CC_PUSH_POP{Pop: false, Segment: "CONSTANT", Index: 111},
		&codegen.CC_PUSH_POP{Pop: false, Segment: "CONSTANT", Index: 333},
		&codegen.CC_PUSH_POP{Pop: false, Segment: "CONSTANT", Index: 888},
		&codegen.CC_PUSH_POP{Pop: true, Segment: "STATIC", Index: 8},
		&codegen.CC_PUSH_POP{Pop: true, Segment: "STATIC", Index: 3},
		&codegen.CC_PUSH_POP{Pop: true, Segment: "STATIC", Index: 1},
		&codegen.CC_PUSH_POP{Pop: false, Segment: "STATIC", Index: 3},
		&codegen.CC_PUSH_POP{Pop: false, Segment: "STATIC", Index: 1},
		&codegen.CC_ARITHMETIC{Op: "sub"},
		&codegen.CC_PUSH_POP{Pop: false, Segment: "STATIC", Index: 8},
		&codegen.CC_ARITHMETIC{Op: "add"},
		&codegen.CC_PUSH_POP{Pop: false, Segment: "CONSTANT", Index: 17},
		&codegen.CC_PUSH_POP{Pop: false, Segment: "CONSTANT", Index: 17},
		&codegen.CC_ARITHMETIC{Op: "eq"},
		&codegen.CC_PUSH_POP{Pop: false, Segment: "CONSTANT", Index: 17},
		&codegen.CC_PUSH_POP{Pop: false, Segment: "CONSTANT", Index: 16},
		&codegen.CC_ARITHMETIC{Op: "eq"},
		&codegen.CC_PUSH_POP{Pop: false, Segment: "CONSTANT", Index: 16},
		&codegen.CC_PUSH_POP{Pop: false, Segment: "CONSTANT", Index: 17},
		&codegen.CC_ARITHMETIC{Op: "eq"},
		&codegen.CC_PUSH_POP{Pop: false, Segment: "CONSTANT", Index: 892},
		&codegen.CC_PUSH_POP{Pop: false, Segment: "CONSTANT", Index: 891},
		&codegen.CC_ARITHMETIC{Op: "lt"},
		&codegen.CC_PUSH_POP{Pop: false, Segment: "CONSTANT", Index: 891},
		&codegen.CC_PUSH_POP{Pop: false, Segment: "CONSTANT", Index: 892},
		&codegen.CC_ARITHMETIC{Op: "lt"},
		&codegen.CC_PUSH_POP{Pop: false, Segment: "CONSTANT", Index: 891},
		&codegen.CC_PUSH_POP{Pop: false, Segment: "CONSTANT", Index: 891},
		&codegen.CC_ARITHMETIC{Op: "lt"},
		&codegen.CC_PUSH_POP{Pop: false, Segment: "CONSTANT", Index: 32767},
		&codegen.CC_PUSH_POP{Pop: false, Segment: "CONSTANT", Index: 32766},
		&codegen.CC_ARITHMETIC{Op: "gt"},
		&codegen.CC_PUSH_POP{Pop: false, Segment: "CONSTANT", Index: 32766},
		&codegen.CC_PUSH_POP{Pop: false, Segment: "CONSTANT", Index: 32767},
		&codegen.CC_ARITHMETIC{Op: "gt"},
		&codegen.CC_PUSH_POP{Pop: false, Segment: "CONSTANT", Index: 32766},
		&codegen.CC_PUSH_POP{Pop: false, Segment: "CONSTANT", Index: 32766},
		&codegen.CC_ARITHMETIC{Op: "gt"},
		&codegen.CC_PUSH_POP{Pop: false, Segment: "CONSTANT", Index: 57},
		&codegen.CC_PUSH_POP{Pop: false, Segment: "CONSTANT", Index: 31},
		&codegen.CC_PUSH_POP{Pop: false, Segment: "CONSTANT", Index: 53},
		&codegen.CC_ARITHMETIC{Op: "add"},
		&codegen.CC_PUSH_POP{Pop: false, Segment: "CONSTANT", Index: 112},
		&codegen.CC_ARITHMETIC{Op: "sub"},
		&codegen.CC_ARITHMETIC{Op: "neg"},
		&codegen.CC_ARITHMETIC{Op: "and"},
		&codegen.CC_PUSH_POP{Pop: false, Segment: "CONSTANT", Index: 82},
		&codegen.CC_ARITHMETIC{Op: "or"},
		&codegen.CC_ARITHMETIC{Op: "not"},
		&codegen.CC_EMPTY{},
	}
	pp := NewParser(mg)
exit:
	for _, cc := range test_expected {
		cmd, err := pp.NextCommand()
		if err != nil {
			t.Errorf("Err %s", err.Error())
			break exit
		}

		switch cmd.(type) {
		case *codegen.CC_EMPTY:
			break exit
		}
		if !reflect.DeepEqual(cc, cmd) {
			t.Errorf("%#v != %#v", cmd, cc)
			break exit
		}
	}

}
