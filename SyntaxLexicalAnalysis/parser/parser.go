package parser

import (
	"fmt"
	"io"
	"os"

	jacklexer "github.com/quickwritereader/JackSyntaxAnalyser/JackLexer"
	jackhelper "github.com/quickwritereader/JackSyntaxAnalyser/jackHelper"
	"github.com/quickwritereader/JackSyntaxAnalyser/jacktoken"
)

type RecursiveDescentParser struct {
	tokenizer      *jacklexer.JackTokenizerApi
	writer         io.StringWriter
	errorHappend   bool
	OutputComments bool
}

var keywordConstant = []string{"true", "false", "null", "this"}
var operations = []string{"+", "-", "*", "/", "&", "|", "<", ">", "="}

func NewRecursiveDescentParser(input_stream io.Reader, output io.StringWriter) *RecursiveDescentParser {
	if output == nil {
		//dummy Writer
		output = &jackhelper.DummyWriter{}
	}
	return &RecursiveDescentParser{
		tokenizer:      jacklexer.NewTokenizer(input_stream),
		writer:         output,
		errorHappend:   false,
		OutputComments: true,
	}
}
func (p *RecursiveDescentParser) outputErrorLine() {
	if p.errorHappend {
		return
	}
	fmt.Fprintf(os.Stderr, "Error in line %d , last token %#v\n", p.tokenizer.LineNumber, p.tokenizer.Current)
	p.errorHappend = true
}

func (p *RecursiveDescentParser) beginNode(str string) {
	jackhelper.WriteStringAsBeginNode(p.writer, str)
}

func (p *RecursiveDescentParser) endNode(str string) {
	jackhelper.WriteStringAsEndNode(p.writer, str)
}

func (p *RecursiveDescentParser) wholeNode(t jacktoken.Token) {
	jackhelper.WriteTokenAsXmlNode(p.writer, t)
}

func (p *RecursiveDescentParser) readToken() jacktoken.Token {
	tokenType := p.tokenizer.Current.Type
	for tokenType == jacktoken.COMMENT || tokenType == jacktoken.MULTI_LINE_COMMENT {
		if p.OutputComments {
			jackhelper.WriteTokenAsXmlComment(p.writer, p.tokenizer.Current)
		}
		if p.tokenizer.HasMoreTokens() {
			//advance
			p.tokenizer.Advance()
			tokenType = p.tokenizer.Current.Type
		}
	}
	return p.tokenizer.Current
}

func (p *RecursiveDescentParser) peekToken() jacktoken.Token {
	//before peeking we should skip if they are comments
	for p.tokenizer.HasMoreTokens() {

		t := p.tokenizer.PeekToken().Type

		if t == jacktoken.COMMENT || t == jacktoken.MULTI_LINE_COMMENT {
			//skip next
			tok := p.tokenizer.AdvanceOnlyNext()
			//lets print what we skipped
			if p.OutputComments {
				jackhelper.WriteTokenAsXmlComment(p.writer, tok)
			}

		} else {
			break
		}

	}
	return p.tokenizer.PeekToken()
}

func (p *RecursiveDescentParser) outputCurrentAndAdvance() {
	p.wholeNode(p.tokenizer.Current)
	p.tokenizer.Advance()
}

func (p *RecursiveDescentParser) checkKeywordList(literals []string) bool {
	t := p.readToken()
	if t.Type == jacktoken.KEYWORD {
		if jackhelper.LocalIndex(literals, t.Literal) >= 0 {
			return true
		}
	}
	return false
}

func (p *RecursiveDescentParser) checkSymbolList(literals []string) bool {
	t := p.readToken()
	if t.Type == jacktoken.SYMBOL {
		if jackhelper.LocalIndex(literals, t.Literal) >= 0 {
			return true
		}
	}
	return false
}

func (p *RecursiveDescentParser) checkKeyword(keyword string) bool {
	t := p.readToken()
	return t.Type == jacktoken.KEYWORD && t.Literal == keyword
}

func (p *RecursiveDescentParser) checkSymbol(keyword string) bool {
	t := p.readToken()
	return t.Type == jacktoken.SYMBOL && t.Literal == keyword
}

func (p *RecursiveDescentParser) checkKeywordAndAdvance(literal string) {
	if p.checkKeyword(literal) {
		p.outputCurrentAndAdvance()
	} else {
		p.outputErrorLine()
	}
}

func (p *RecursiveDescentParser) checkSymbolAndAdvance(literal string) {
	if p.checkSymbol(literal) {
		p.outputCurrentAndAdvance()
	} else {
		p.outputErrorLine()
	}
}

func (p *RecursiveDescentParser) checkIdentifierAndAdvance() {
	if p.readToken().Type == jacktoken.IDENTIFIER {
		p.outputCurrentAndAdvance()
	} else {
		p.outputErrorLine()
	}
}

func (p *RecursiveDescentParser) checkType0(checkVoid bool) bool {
	//type: 'int'|'char'|'boolean' | className
	t := p.readToken()
	if t.Type == jacktoken.IDENTIFIER {
		return true
	} else if t.Type == jacktoken.KEYWORD && (t.Literal == "int" || t.Literal == "char" || t.Literal == "boolean") {
		return true
	} else if checkVoid && t.Literal == "void" {
		return true
	}

	return false

}

func (p *RecursiveDescentParser) checkType() bool {
	return p.checkType0(false)

}

func (p *RecursiveDescentParser) checkTypeAll() bool {
	return p.checkType0(true)
}

func (p *RecursiveDescentParser) checkTypeAndAdvance() {
	if p.checkType() {
		p.outputCurrentAndAdvance()
	} else {
		p.outputErrorLine()
	}
}

func (p *RecursiveDescentParser) checkTypeAllAndAdvance() {
	if p.checkTypeAll() {
		p.outputCurrentAndAdvance()
	} else {
		p.outputErrorLine()
	}
}

func (p *RecursiveDescentParser) CompileAsXml() {
	if p.tokenizer.HasMoreTokens() {
		p.tokenizer.Advance()
	}
	if p.checkKeyword("class") {
		p.beginNode("class")
		p.compileClass()
		p.endNode("class")
	} else {
		p.outputErrorLine()
	}
	//check if there are inputs left
	//as we have one class per File , we have to output error
	if !p.errorHappend && p.readToken().Type != jacktoken.EOF {
		p.outputErrorLine()
	}

}

func (p *RecursiveDescentParser) compileClass() {
	//'class' className '{' classVarDec* subroutineDec* '}'
	p.checkKeywordAndAdvance("class")
	p.checkIdentifierAndAdvance()
	p.checkSymbolAndAdvance("{")

	// we could use map to functions as well, like in non recursive table
	for !p.errorHappend && p.checkKeywordList([]string{"static", "field"}) {
		p.compileClassVarDec()
	}
	for !p.errorHappend && p.checkKeywordList([]string{"constructor", "method", "function"}) {
		p.compileSubroutine()
	}
	//check }
	p.checkSymbolAndAdvance("}")

}

func (p *RecursiveDescentParser) checkVarDecList() {
	//check type
	//process the first default
	p.checkTypeAndAdvance()
	p.checkIdentifierAndAdvance()

	t := p.readToken()
	//, varName
	for !p.errorHappend && t.Type == jacktoken.SYMBOL && t.Literal == "," {
		p.outputCurrentAndAdvance()

		p.checkIdentifierAndAdvance()

		t = p.readToken()
	}
	p.checkSymbolAndAdvance(";")
}

func (p *RecursiveDescentParser) compileGeneralVarDec(nodeLabel string) {
	//('static'|'field' | 'var' ) type varName (',' varName)* ';'
	p.beginNode(nodeLabel)
	//jus output and advance, we already checked it
	p.outputCurrentAndAdvance()
	p.checkVarDecList()
	p.endNode(nodeLabel)

}

func (p *RecursiveDescentParser) compileClassVarDec() {
	//('static'|'field') type varName (',' varName)* ';'
	p.compileGeneralVarDec("classVarDec")

}

func (p *RecursiveDescentParser) compileVarDec() {
	//varDec: 'var' type varName (',' varName)* ';'
	p.compileGeneralVarDec("varDec")
}

func (p *RecursiveDescentParser) checkParameterList() {
	//((type varName) (',' type varName)*)?
	p.beginNode("parameterList")
	if !p.checkType() || p.peekToken().Type != jacktoken.IDENTIFIER {
		//empty, just return
		p.endNode("parameterList")
		return
	}

	p.checkTypeAndAdvance()
	p.checkIdentifierAndAdvance()
	for p.checkSymbol(",") {
		//optional list
		p.outputCurrentAndAdvance()
		p.checkTypeAndAdvance()
		p.checkIdentifierAndAdvance()

	}
	p.endNode("parameterList")
}

func (p *RecursiveDescentParser) compileSubroutine() {
	/*
		('constructor'|'function'|'method')
		('void' | type) subroutineName '(' parameterList ')'
		subroutineBody
	*/
	p.beginNode("subroutineDec")
	p.outputCurrentAndAdvance()

	//'void' | type
	p.checkTypeAllAndAdvance()
	p.checkIdentifierAndAdvance()
	p.checkSymbolAndAdvance("(")
	p.checkParameterList()
	p.checkSymbolAndAdvance(")")
	p.compileSubroutineBody()
	p.endNode("subroutineDec")

}

func (p *RecursiveDescentParser) compileSubroutineBody() {
	// '{' varDec* statements '}'
	p.beginNode("subroutineBody")
	p.checkSymbolAndAdvance("{")

	//varDec*

	for !p.errorHappend && p.checkKeyword("var") {
		p.compileVarDec()
	}
	p.compileStatements()
	p.checkSymbolAndAdvance("}")
	p.endNode("subroutineBody")
}

func (p *RecursiveDescentParser) compileStatements() {
	// statements: statement*
	// 	statement: letStatement | ifStatement | whileStatement |
	// doStatement | returnStatement
	p.beginNode("statements")
	tok := p.readToken()

	for tok.Type == jacktoken.KEYWORD {
		if tok.Literal == "if" {
			p.compileIf()
		} else if tok.Literal == "do" {
			p.compileDo()
		} else if tok.Literal == "let" {
			p.compileLet()
		} else if tok.Literal == "return" {
			p.compileReturn()
		} else if tok.Literal == "while" {
			p.compileWhile()
		} else {
			break
		}

		tok = p.readToken()

	}

	p.endNode("statements")
}

func (p *RecursiveDescentParser) compileDo() {
	// doStatement: 'do' subroutineCall ';'
	p.beginNode("doStatement")
	p.outputCurrentAndAdvance()
	p.compileSubroutineCall()
	p.checkSymbolAndAdvance(";")
	p.endNode("doStatement")
}

func (p *RecursiveDescentParser) compileLet() {
	// letStatement: 'let' varName ('[' expression ']')? '=' expression ';'
	p.beginNode("letStatement")
	p.outputCurrentAndAdvance()
	p.checkIdentifierAndAdvance()
	if p.checkSymbol("[") {
		//just ouput as its checked
		p.outputCurrentAndAdvance()
		p.compileExpression()
		p.checkSymbolAndAdvance("]")
	}
	p.checkSymbolAndAdvance("=")
	p.compileExpression()
	p.checkSymbolAndAdvance(";")
	p.endNode("letStatement")
}

func (p *RecursiveDescentParser) compileWhile() {
	// whileStatement: 'while' '(' expression ')' '{' statements '}'
	p.beginNode("whileStatement")
	//while is already checked, just output and advance
	p.outputCurrentAndAdvance()
	p.checkSymbolAndAdvance("(")
	p.compileExpression()
	p.checkSymbolAndAdvance(")")
	p.checkSymbolAndAdvance("{")
	p.compileStatements()
	p.checkSymbolAndAdvance("}")
	p.endNode("whileStatement")
}

func (p *RecursiveDescentParser) compileReturn() {

	// ReturnStatement 'return' expression? ';'
	p.beginNode("returnStatement")
	p.outputCurrentAndAdvance()
	if p.checkSymbol(";") {
		p.outputCurrentAndAdvance()
	} else {
		p.compileExpression()
		p.checkSymbolAndAdvance(";")

	}
	p.endNode("returnStatement")

}

func (p *RecursiveDescentParser) compileIf() {
	// ifStatement: 'if' '(' expression ')' '{' statements '}'
	// ('else' '{' statements '}')?
	p.beginNode("ifStatement")
	p.outputCurrentAndAdvance()
	p.checkSymbolAndAdvance("(")
	p.compileExpression()
	p.checkSymbolAndAdvance(")")
	p.checkSymbolAndAdvance("{")
	p.compileStatements()
	p.checkSymbolAndAdvance("}")

	if !p.checkKeyword("else") {
		//if no else , just return
		p.endNode("ifStatement")
		return
	}
	//its else
	p.outputCurrentAndAdvance()
	//
	p.checkSymbolAndAdvance("{")
	p.compileStatements()
	p.checkSymbolAndAdvance("}")
	p.endNode("ifStatement")
}

func (p *RecursiveDescentParser) compileExpression() {
	// term (op term)*
	if p.errorHappend {
		return
	}
	p.beginNode("expression")
	p.compileTerm()
	for !p.errorHappend && p.checkSymbolList(operations) {
		p.outputCurrentAndAdvance()
		p.compileTerm()

	}
	p.endNode("expression")
}

func (p *RecursiveDescentParser) compileTerm() {
	// integerConstant | stringConstant | keywordConstant |
	// varName | varName '[' expression ']' | subroutineCall |
	// '(' expression ')' | unaryOp term
	p.beginNode("term")
	tok := p.readToken()
	peek := p.peekToken()
	if tok.Type == jacktoken.INTEGER_CONSTANT || tok.Type == jacktoken.STRING_CONSTANT {
		p.outputCurrentAndAdvance()
	} else if p.checkKeywordList(keywordConstant) {
		p.outputCurrentAndAdvance()
	} else if tok.Type == jacktoken.IDENTIFIER && peek.Type == jacktoken.SYMBOL && peek.Literal == "[" {
		//varName '[' expression ']'
		p.outputCurrentAndAdvance()
		p.outputCurrentAndAdvance()
		p.compileExpression()
		p.checkSymbolAndAdvance("]")

	} else if tok.Type == jacktoken.SYMBOL && (tok.Literal == "-" || tok.Literal == "~") && peek.Type != jacktoken.STRING_CONSTANT {
		//unlike in the original we disallow unaryOp for  string
		p.outputCurrentAndAdvance()
		p.compileTerm()
	} else if p.checkSymbol("(") {
		p.outputCurrentAndAdvance()
		p.compileExpression()
		p.checkSymbolAndAdvance(")")
	} else if tok.Type == jacktoken.IDENTIFIER && peek.Type == jacktoken.SYMBOL && (peek.Literal == "(" || peek.Literal == ".") {
		//try subroutine Call
		p.compileSubroutineCall()
	} else if tok.Type == jacktoken.IDENTIFIER {
		//varName
		p.outputCurrentAndAdvance()
	} else {
		p.outputErrorLine()
	}
	p.endNode("term")
}

func (p *RecursiveDescentParser) compileSubroutineCall() {
	// subroutineName '(' expressionList ')' |
	// (className | varName) '.' subroutineName '(' expressionList ')'

	p.checkIdentifierAndAdvance()

	if p.checkSymbol("(") {
		p.outputCurrentAndAdvance()
	} else if p.checkSymbol(".") {
		p.outputCurrentAndAdvance()
		//subroutineName '('
		p.checkIdentifierAndAdvance()
		p.checkSymbolAndAdvance("(")
	}

	p.compileExpressionList()
	p.checkSymbolAndAdvance(")")

}

func (p *RecursiveDescentParser) compileExpressionList() {
	//(expression (',' expression)* )?
	//? means optional
	p.beginNode("expressionList")
	if !p.checkSymbol(")") {
		p.compileExpression()
		for !p.errorHappend && p.checkSymbol(",") {
			p.outputCurrentAndAdvance()
			p.compileExpression()
		}

	}
	p.endNode("expressionList")

}
