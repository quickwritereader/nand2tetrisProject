package codegen

import (
	"strings"
	"testing"
)

func TestCodeGenPopStatic(t *testing.T) {
	var sb strings.Builder
	cg := NewCodeGenFile(&sb, "SSS")
	cg.Process(&CC_PUSH_POP{true, STATIC, 15})
	str := strings.TrimSpace(sb.String())

	expected := "@SP\nAM=M-1\nD=M\n@SSS.15\nM=D\n"

	compareIgnoringCommentsAndSpaces(t, expected, str)
}

func TestInitCode(t *testing.T) {
	var sb strings.Builder
	cg := NewCodeGenFile(&sb, "SSS")
	cg.WriteInit()

	expected := `
	@256
	D=A
	@SP
	M=D
	///// BEGIN_SECTION: GLOBAL CALL CODE /////
	@$END_GLOBAL_CALL_CODE
	0;JMP
	($GLOBAL_CALL_CODE)
	@SP
	A=M
	M=D
	@SP
	M=M+1
	@LCL
	D=M
	@SP
	A=M
	M=D
	@SP
	M=M+1
	@ARG
	D=M
	@SP
	A=M
	M=D
	@SP
	M=M+1
	@THIS
	D=M
	@SP
	A=M
	M=D
	@SP
	M=M+1
	@THAT
	D=M
	@SP
	A=M
	M=D
	@SP
	M=M+1
	@SP
	D=M
	@R13
	D=D-M
	@5
	D=D-A
	@ARG
	M=D
	@SP
	D=M
	@LCL
	M=D
	@R14
	A=M
	0;JMP
	($END_GLOBAL_CALL_CODE)
	///// BEGIN_SECTION: GLOBAL RETURN CODE /////
	@$END_GLOBAL_RETURN_CODE
	0;JMP
	($GLOBAL_RETURN_CODE)
	@5
	D=A
	@LCL
	A=M-D
	D=M
	@R13
	M=D
	@SP
	AM=M-1
	D=M
	@ARG
	A=M
	M=D
	D=A
	@SP
	M=D+1
	@LCL
	D=M
	@R14
	AM=D-1
	D=M
	@THAT
	M=D
	@R14
	AM=M-1
	D=M
	@THIS
	M=D
	@R14
	AM=M-1
	D=M
	@ARG
	M=D
	@R14
	AM=M-1
	D=M
	@LCL
	M=D
	@R13
	A=M
	0;JMP
	($END_GLOBAL_RETURN_CODE)
	///// END_SECTION: GLOBAL RETURN CODE /////
	//Call Sys.init 0
	@0
	D=A
	@R13
	M=D
	@Sys.init
	D=A
	@R14
	M=D
	@Sys.init$ret0
	D=A
	//goto $GLOBAL_CALL_CODE
	@$GLOBAL_CALL_CODE
	0;JMP
	(Sys.init$ret0)
	`
	compareIgnoringCommentsAndSpaces(t, expected, sb.String())
}

func compareIgnoringCommentsAndSpaces(t *testing.T, expected, s string) {
	t.Helper()
	strArr := strings.Split(strings.TrimSpace(s), "\n")
	expectedArr := strings.Split(strings.TrimSpace(expected), "\n")
	i, j := 0, 0
	exp := ""
	actual := ""
	for {

		if i < len(strArr) {
			actual = strings.TrimSpace(strArr[i])
			if len(actual) < 1 || strings.HasPrefix(actual, "//") {
				//ignore comment and empty, try next
				i++
				continue
			}
		} else {
			actual = ""
		}
		if j < len(expectedArr) {
			exp = strings.TrimSpace(expectedArr[j])
			if len(exp) < 1 || strings.HasPrefix(exp, "//") {
				//ignore comment and empty, try next
				j++
				continue
			}
		} else {
			exp = ""
		}
		if actual != exp {
			t.Error("\"Actual: ", actual, "\" Expected: \"", exp, "\"")
		}

		if i >= len(strArr) && j >= len(expectedArr) {
			//both ended
			break
		}

		i++
		j++
	}
}
