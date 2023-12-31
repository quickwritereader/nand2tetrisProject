package codegen

import (
	"fmt"
	"io"
)

const (
	STATIC   = "STATIC"
	LOCAL    = "LOCAL"
	CONSTANT = "CONSTANT"
	TEMP     = "TEMP"
	THIS     = "THIS"
	THAT     = "THAT"
	ARGUMENT = "ARGUMENT"
	UNKNOWN  = "UNKNOWN"
	POINTER  = "POINTER"
)

type MemSegment string

type CodeGen struct {
	wrt                io.StringWriter
	name               string
	unique_index       int
	function_ret_index map[string]int

	global_function_call_code_generated   bool
	global_function_return_code_generated bool
}

type Command interface {
	EmitAsm(cg *CodeGen)
}

func NewCodeGen() CodeGen {
	return NewCodeGenFile(nil, "")
}

func NewCodeGenFile(writer io.StringWriter, name string) CodeGen {
	return CodeGen{writer, name, 0, make(map[string]int), false, false}
}

func (cg *CodeGen) SetWriter(writer io.StringWriter) {
	cg.wrt = writer
}

func (cg *CodeGen) SetName(name string) {
	cg.name = name
}

func (cg *CodeGen) Process(cmd Command) {
	cmd.EmitAsm(cg)
}

func (cg *CodeGen) emit(asm string) {
	cg.wrt.WriteString(asm)
}

func (cg *CodeGen) PopHeadToD() {
	//*SP--; D=*SP
	//  A=SP-1 ;*SP=*SP-1;  D=*SP;
	cg.wrt.WriteString("@SP\nAM=M-1\nD=M\n")
}

func (cg *CodeGen) PushD() {
	//*SP=*addr; SP++
	cg.wrt.WriteString("@SP\nA=M\nM=D\n@SP\nM=M+1\n")
}

func (cg *CodeGen) WriteInit() {
	//SP=256
	cg.emit("/////////////////////INITIALIZATION////////////////////\n")
	cg.emit("@256\nD=A\n@SP\nM=D\n")
	cg.GenerateGlobalCallCode()
	cg.GenerateGlobalReturnCode()
	cg.emit("//Call Sys.init 0\n")
	call := &CC_CALL{FunctionName: "Sys.init", ArgumentCount: 0}
	call.EmitAsm(cg)
	cg.emit("/////////////////////////////////////////////////////\n")
	cg.emit("/////////////////////////////////////////////////////\n")
}

func (cg *CodeGen) GenerateGlobalReturnCode() {
	//LETS GENERATE GLOBAL RETURN FUNCTIONALITY
	if cg.global_function_return_code_generated {
		return
	}
	cg.global_function_return_code_generated = true
	//lets code skip, jump over it
	cg.emit("///// BEGIN_SECTION: GLOBAL RETURN CODE /////\n")
	cg.emit("@$END_GLOBAL_RETURN_CODE\n")
	cg.emit("0;JMP\n")
	//our GLOBAL_RETURN_CODE
	cg.emit("($GLOBAL_RETURN_CODE)\n")
	//   ---                <-   ARG
	//   ---
	//   0   RETURN
	//	 1	Saved LCL
	//	 2	Saved ARG
	//	 3	Saved THIS
	//	 4	Saved THAT
	//	 end_frame                <- LOCAL LCL
	//as we see we can get others through LCL content
	// R13 = *(LCL-5)  //R13 will contain return address
	cg.emit("@5\nD=A\n")
	cg.emit("@LCL\nA=M-D\nD=M\n@R13\nM=D\n")
	//*ARG = pop() // we should store result into ARG pointed stack
	cg.PopHeadToD()
	cg.emit("@ARG\nA=M\nM=D\n")
	//store ARG value into D as well
	cg.emit("D=A\n")
	// SP = ARG+1 // new SP should point after ARG
	cg.emit("@SP\nM=D+1\n")

	//let's restore THAT = *(LCL-1) ; THIS=*(LCL-2)
	// lets store LCL into temp
	//tmp = LCL ; so that we can use it as --tmp way, and also restore LCL itself later
	//cg.emit("@LCL\nD=M\n@R14\nM=D\n")

	for index, restore_reg := range [...]string{"THAT", "THIS", "ARG", "LCL"} {

		if index == 0 {
			//first time we can get tmp from LCL
			//tmp=LCL-1
			cg.emit("@LCL\nD=M\n")
			cg.emit("@R14\nAM=D-1\n")
		} else {
			//tmp=tmp-1
			cg.emit("@R14\nAM=M-1\n")
		}

		//restore reg; reg=*tmp;
		//D=*tmp; reg=D
		cg.emit("D=M\n@")
		cg.emit(restore_reg)
		cg.emit("\nM=D\n")
	}
	//jmp to return address.
	//return address was in R13
	cg.emit("@R13\nA=M\n0;JMP\n")

	//mark end of it for skip purpose
	cg.emit("($END_GLOBAL_RETURN_CODE)\n")
	cg.emit("///// END_SECTION: GLOBAL RETURN CODE /////\n")

}

func (cg *CodeGen) GenerateGlobalCallCode() {
	//LETS GENERATE GLOBAL CALL FUNCTIONALITY
	if cg.global_function_call_code_generated {
		return
	}
	cg.global_function_call_code_generated = true
	//D should be return address
	//R13 should be arg size
	//R14 should be called function address

	//lets code skip, jump over it
	cg.emit("///// BEGIN_SECTION: GLOBAL CALL CODE /////\n")
	cg.emit("@$END_GLOBAL_CALL_CODE\n")
	cg.emit("0;JMP\n")
	//our GLOBAL_CALL_CODE
	cg.emit("($GLOBAL_CALL_CODE)\n")
	//   ---                <-   ARG
	//   ---
	//   0   RETURN
	//	 1	Saved LCL
	//	 2	Saved ARG
	//	 3	Saved THIS
	//	 4	Saved THAT
	//	 end_frame                <- LOCAL LCL

	//push return Address (it should be in D)
	cg.PushD()
	//push LCL,ARG,THIS, THAT
	for _, reg := range [...]string{"LCL", "ARG", "THIS", "THAT"} {
		cg.emit("@")
		cg.emit(reg)
		cg.emit("\nD=M\n")
		cg.PushD()
	}
	//reposition ARG
	//ARG = SP-5-nArgs(R13)
	cg.emit("@SP\nD=M\n@R13\nD=D-M\n@5\nD=D-A\n")
	cg.emit("@ARG\nM=D\n")
	//reposition LCL
	//LCL=SP
	cg.emit("@SP\nD=M\n@LCL\nM=D\n")
	//goto Function Address (we stored it in R14)
	cg.emit("@R14\nA=M\n0;JMP\n")
	//mark end of it for skip purpose
	cg.emit("($END_GLOBAL_CALL_CODE)\n")
	cg.emit("///// END_SECTION: GLOBAL CALL CODE /////\n")
}

type CC_EMPTY struct{}

type CC_COMMENT struct {
	Text string
}

type CC_PUSH_POP struct {
	Pop     bool
	Segment MemSegment
	Index   int
}

type CC_LABEL struct {
	Name string
}

type CC_GOTO struct {
	Label string
}

type CC_IF_GOTO struct {
	Label string
}

type CC_ARITHMETIC struct {
	Op string
}

type CC_FUNCTION struct {
	FunctionName string
	LocalCount   int
}

type CC_CALL struct {
	FunctionName  string
	ArgumentCount int
}

type CC_RETURN struct {
}

func (cc *CC_EMPTY) EmitAsm(cg *CodeGen) {

}

func (cc *CC_COMMENT) EmitAsm(cg *CodeGen) {
	cg.emit(cc.Text)
	cg.emit("\n")
}

func writeFuncPrefix(cg *CodeGen) {
	cg.emit(cg.name)
	cg.emit("_")
}

func (cc *CC_LABEL) EmitAsm(cg *CodeGen) {
	cg.emit("(")
	writeFuncPrefix(cg)
	cg.emit(cc.Name)
	cg.emit(")\n")
}

func (cc *CC_GOTO) EmitAsm(cg *CodeGen) {
	cg.emit("@")
	writeFuncPrefix(cg)
	cg.emit(cc.Label)
	cg.emit("\n0;JMP\n")
}

func (cc *CC_IF_GOTO) EmitAsm(cg *CodeGen) {
	cg.PopHeadToD()
	cg.emit("@")
	writeFuncPrefix(cg)
	cg.emit(cc.Label)
	cg.emit("\nD;JNE\n")
}

func (cc *CC_FUNCTION) EmitAsm(cg *CodeGen) {
	//for each function we have to get Label

	cg.emit(fmt.Sprintf("(%s)\n", cc.FunctionName))
	//we have to push  local args as zeroed
	//so if it was unintialized we could just subtract SP
	//so lets execute simple loop for zeroing
	//Generate code if only if LocalCount > 0
	if cc.LocalCount > 0 {
		cg.emit(fmt.Sprintf("@%d\nD=A\n", cc.LocalCount))
		//lets zero contents of SP points *SP++=0
		cg.emit(fmt.Sprintf("(%s$LOOP_PUSH_ZEROED_LOCALS)\n", cc.FunctionName))
		cg.emit("D=D-1\n@SP\nAM=M+1\nA=A-1\nM=0\n@")
		cg.emit(cc.FunctionName)
		cg.emit("$LOOP_PUSH_ZEROED_LOCALS\nD;JGT\n")
	}

}

func (cc *CC_RETURN) EmitAsm(cg *CodeGen) {
	//FRAME SHOULD BE RESET
	//as this procedure is almost the same
	//we will generate one global and just jump to it
	if !cg.global_function_return_code_generated {
		cg.GenerateGlobalReturnCode()
	}
	//lets emit jmp into $END_GLOBAL_RETURN_CODE
	cg.emit("@$GLOBAL_RETURN_CODE\n0;JMP\n")

}

func (cc *CC_CALL) EmitAsm(cg *CodeGen) {
	if !cg.global_function_call_code_generated {
		cg.GenerateGlobalCallCode()
	}
	cg.emit("//call Function Gen:\n")
	//R13 should be arg size
	cg.emit(fmt.Sprintf("@%d\n", cc.ArgumentCount))
	cg.emit("D=A\n@R13\nM=D\n")
	//R14 should be called function address
	cg.emit("@" + cc.FunctionName)
	cg.emit("\nD=A\n@R14\nM=D\n")
	//D should be return address
	//generate return LABEl and put it in D
	ret_index := cg.function_ret_index[cc.FunctionName]
	return_label_name := fmt.Sprintf("%s$ret%d", cc.FunctionName, ret_index)
	cg.function_ret_index[cc.FunctionName] = ret_index + 1

	cg.emit("@")
	cg.emit(return_label_name)
	cg.emit("\nD=A\n")
	//jmp to $GLOBAL_CALL_CODE
	cg.emit("//goto $GLOBAL_CALL_CODE\n")
	cg.emit("@$GLOBAL_CALL_CODE\n0;JMP\n")

	cg.emit("(")
	cg.emit(return_label_name)
	cg.emit(")\n")

}

func (cc *CC_PUSH_POP) EmitAsm(cg *CodeGen) {
	if cc.Pop {
		cc.emitPop(cg)
	} else {
		cc.emitPush(cg)
	}
}

func getSeg(t MemSegment) string {
	if t == ARGUMENT {
		return "ARG"
	} else if t == LOCAL {
		return "LCL"
	} else if t == THIS {
		return "THIS"
	} else if t == THAT {
		return "THAT"
	} else {
		//error
		return ""
	}
}

func (cc *CC_PUSH_POP) emitPush(cg *CodeGen) {

	switch cc.Segment {
	case CONSTANT:
		cg.emit(fmt.Sprintf("@%v\nD=A\n", cc.Index))
	case TEMP:
		cg.emit(fmt.Sprintf("@%v\nD=M\n", cc.Index+5))
	case POINTER:
		th := "THAT"
		if cc.Index == 0 {
			th = "THIS"
		}
		cg.emit(fmt.Sprintf("@%v\nD=M\n", th))
	case STATIC:
		cg.emit(fmt.Sprintf("@%v.%v\nD=M\n", cg.name, cc.Index))
	case ARGUMENT, LOCAL, THIS, THAT:
		//addr = [[SegmentPointer] + index ]
		seg := getSeg(cc.Segment)
		cg.emit(fmt.Sprintf("@%v\nD=A\n@%v\nA=D+M\nD=M\n", cc.Index, seg))

	}
	cg.PushD()
}

func (cc *CC_PUSH_POP) emitPop(cg *CodeGen) {

	//calculate segment address into temp R13
	if cc.Segment == ARGUMENT || cc.Segment == LOCAL || cc.Segment == THIS || cc.Segment == THAT {
		seg := getSeg(cc.Segment)
		//D contains addr= ([[SegmentPointer] + index])
		cg.emit(fmt.Sprintf("@%v\nD=A\n@%v\nD=D+M\n", cc.Index, seg))
		//store Address into [R13] ,
		cg.emit("@R13\nM=D\n")
	}

	//*SP--; D=* // result: A=SP-1 ;*SP=*SP-1;  D=*SP;
	cg.PopHeadToD()

	//get addr to be stored
	switch cc.Segment {
	case TEMP:
		cg.emit(fmt.Sprintf("@%v\n", cc.Index+5))
	case POINTER:
		th := "THAT"
		if cc.Index == 0 {
			th = "THIS"
		}
		cg.emit(fmt.Sprintf("@%v\n", th))
	case STATIC:
		cg.emit(fmt.Sprintf("@%v.%v\n", cg.name, cc.Index))
	case ARGUMENT, LOCAL, THIS, THAT:
		//restore address from [R13] into A
		cg.emit("@R13\nA=M\n")

	}
	//*addr=*SP -> *addr = D
	cg.emit("M=D\n")

}

func (cc *CC_ARITHMETIC) EmitAsm(cg *CodeGen) {
	//pop, pop , operate, push

	if cc.Op != "neg" && cc.Op != "not" {
		//result: A=SP-1 ;*SP=*SP-1;  D=*SP;
		cg.PopHeadToD()
		//pop next. our A is already SP=SP-1
		//so we just need to subtract 1 to point to the second
		cg.emit("A=A-1\n")
		//*SP will be just M
		if cc.Op == "lt" || cc.Op == "gt" || cc.Op == "eq" {
			//store M =0 as failure , but later change
			cg.emit("D=M-D\nM=0\n")
			//lets examine D if its not greater or lt, or not eq
			//then we will push SP -1
			//if
			label := fmt.Sprintf("$LABEL%d_%s", cg.unique_index, cc.Op)
			cg.emit("@")
			cg.emit(label)
			cg.unique_index += 1 //for uniq
			jmp_condition := "JNE"
			if cc.Op == "gt" {
				jmp_condition = "JLE"
			} else if cc.Op == "lt" {
				jmp_condition = "JGE"
			}
			cg.emit("\nD;")
			cg.emit(jmp_condition)
			//lets, push -1 in that case
			//our SP needs to be reload with -1,
			cg.emit("\n@SP\nA=M-1\nM=-1\n")
			cg.emit("(")
			cg.emit(label)
			cg.emit(")\n")

		} else {
			switch cc.Op {
			case "add":
				cg.emit("D=D+M\n")
			case "sub":
				cg.emit("D=M-D\n")
			case "and":
				cg.emit("D=D&M\n")
			case "or":
				cg.emit("D=D|M\n")
			}
			//our A is SP pointing to the second argument
			//so we can just store the result there
			// and we do not need to increase SP, as it was changed to the first arg
			// after push it is the same
			cg.emit("M=D\n")
		}
	} else {
		//decrease SP in A, but do not pop
		cg.emit("@SP\nA=M-1\n")
		//now M points to [SP-1]
		if cc.Op == "neg" {
			cg.emit("D=-M\n")
		} else {
			cg.emit("D=!M\n")
		}
		cg.emit("M=D\n")
		//No need to increase SP, as we did not decrease its value in memory

	}

}
