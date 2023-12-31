package parser

import (
	"fmt"
	"io"
	"strconv"

	"github.com/quickwritereader/vmTranslator/codegen"
	"github.com/quickwritereader/vmTranslator/vmlexer"
	"github.com/quickwritereader/vmTranslator/vmtoken"
)

type Parser struct {
	lexer *vmlexer.Lexer
	//currentToken vmtoken.Token
}

func NewParser(input_stream io.Reader) *Parser {
	return &Parser{
		lexer: vmlexer.New(input_stream),
	}
}

func getSegment(t vmtoken.TokenType) codegen.MemSegment {
	switch t {
	case vmtoken.ARGUMENT:
		return codegen.ARGUMENT
	case vmtoken.LOCAL:
		return codegen.LOCAL
	case vmtoken.CONSTANT:
		return codegen.CONSTANT
	case vmtoken.TEMP:
		return codegen.TEMP
	case vmtoken.THAT:
		return codegen.THAT
	case vmtoken.THIS:
		return codegen.THIS
	case vmtoken.STATIC:
		return codegen.STATIC
	case vmtoken.POINTER:
		return codegen.POINTER

	}
	return codegen.UNKNOWN
}

func returnErr(msg string, prev error) error {
	if prev != nil {
		return fmt.Errorf(msg+" (%w )", prev)
	} else {
		return fmt.Errorf(msg)
	}
}

func (p *Parser) NextCommand() (codegen.Command, error) {
	var cc codegen.Command = &codegen.CC_EMPTY{}
	var err error = nil
	tok := p.lexer.NextToken()
	switch tok.Type {
	case vmtoken.COMMENT:
		cc = &codegen.CC_COMMENT{Text: tok.Literal}
	case vmtoken.ARITHMETIC:
		cc = &codegen.CC_ARITHMETIC{Op: tok.Literal}
	case vmtoken.PUSH, vmtoken.POP:
		cpp := codegen.CC_PUSH_POP{Pop: tok.Type == vmtoken.POP}
		//find segment
		tok = p.lexer.NextToken()
		cpp.Segment = getSegment(tok.Type)
		if cpp.Segment == codegen.UNKNOWN {
			//error
			err = fmt.Errorf("segment is unkown error")
		}
		tok = p.lexer.NextToken()
		if tok.Type == vmtoken.NUMBER {
			cpp.Index, _ = strconv.Atoi(tok.Literal)
		} else {
			err = returnErr("the second argument of push/pop is not a number", err)

		}
		cc = &cpp
	case vmtoken.LABEL:
		tok = p.lexer.NextToken()
		//should be identifier
		if tok.Type == vmtoken.IDENTIFIER {
			cc = &codegen.CC_LABEL{Name: tok.Literal}
		} else {
			err = fmt.Errorf("identifier should come after %s", "label")
		}
	case vmtoken.GOTO:
		tok = p.lexer.NextToken()
		//should be identifier
		if tok.Type == vmtoken.IDENTIFIER {
			cc = &codegen.CC_GOTO{Label: tok.Literal}
		} else {
			err = fmt.Errorf("identifier should come after %s", "goto")
		}
	case vmtoken.IF_GOTO:
		tok = p.lexer.NextToken()
		//should be identifier
		if tok.Type == vmtoken.IDENTIFIER {
			cc = &codegen.CC_IF_GOTO{Label: tok.Literal}
		} else {
			err = fmt.Errorf("identifier should come after %s", "if-goto")
		}
	case vmtoken.RETURN:
		cc = &codegen.CC_RETURN{}
	case vmtoken.CALL:
		call := codegen.CC_CALL{}
		tok = p.lexer.NextToken()
		//should be identifier
		if tok.Type == vmtoken.IDENTIFIER {
			call.FunctionName = tok.Literal
		} else {
			err = fmt.Errorf("identifier should come after %s", "if-goto")
		}
		tok = p.lexer.NextToken()
		if tok.Type == vmtoken.NUMBER {
			call.ArgumentCount, _ = strconv.Atoi(tok.Literal)
		} else {
			err = returnErr("the second argument of call is not a number", err)

		}
		cc = &call
	case vmtoken.FUNCTION:
		function := codegen.CC_FUNCTION{}
		tok = p.lexer.NextToken()
		//should be identifier
		if tok.Type == vmtoken.IDENTIFIER {
			function.FunctionName = tok.Literal
		} else {
			err = fmt.Errorf("identifier should come after %s", "if-goto")
		}
		tok = p.lexer.NextToken()
		if tok.Type == vmtoken.NUMBER {
			function.LocalCount, _ = strconv.Atoi(tok.Literal)
		} else {
			err = returnErr("the second argument of function is not a number", err)
		}
		cc = &function
	}

	return cc, err
}
