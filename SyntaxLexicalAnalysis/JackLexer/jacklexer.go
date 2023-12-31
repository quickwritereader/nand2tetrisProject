package jacklexer

import (
	"fmt"
	"io"
	"strings"

	"github.com/quickwritereader/JackSyntaxAnalyser/jacktoken"
)

const (
	BUFFER_SIZE      = 4096
	MIN_VALID_BUFFER = 1
)

type Lexer struct {
	istream          io.Reader
	buffer           []byte
	r, w             int // buf read and write positions
	eof              bool
	last_line_number int
	ch               byte //current character
}

type JackTokenizerApi struct {
	lexer           *Lexer
	Current         jacktoken.Token
	LineNumber      int
	next            jacktoken.Token
	next_lineNumber int
}

func New(input_stream io.Reader) *Lexer {
	return NewBufferSize(input_stream, BUFFER_SIZE)
}

func NewBufferSize(input_stream io.Reader, buffer_size int) *Lexer {
	if buffer_size < MIN_VALID_BUFFER {
		buffer_size = MIN_VALID_BUFFER
	}
	l := &Lexer{istream: input_stream, buffer: make([]byte, buffer_size)}
	l.eof = false
	l.r = -1
	l.readChar()
	return l
}

func (l *Lexer) readBuffer() {
	buf_size := len(l.buffer)
	if !l.eof && l.r+1 >= l.w {
		//l.positionBase += l.w
		//fmt.Println("\nreadBuffer called")
		n, errx := l.istream.Read(l.buffer)
		if n > 0 {
			l.r = -1
			l.w = n
		}

		if errx != nil || n < buf_size {
			l.eof = true
		}

	}

}

func (l *Lexer) readChar() {
	l.ch = l.PeekNextChar()
	if l.r < l.w {
		l.r += 1
	}
	//remember line number
	if l.ch == '\n' {
		l.last_line_number++
	}
}

func (l *Lexer) skipWhiteSpaces() {
	for l.ch == ' ' || l.ch == '\t' || l.ch == '\r' || l.ch == '\n' {
		l.readChar()
	}
}

func (l *Lexer) readBufferConditional(predicate func(ch byte) bool) string {
	//fmt.Println("readBufferConditional:")
	var s strings.Builder
	var pos int = l.r
	if pos < 0 {
		pos = 0
		//we have to write our current character
		s.WriteByte(l.ch)
	}
	for {

		if l.r+1 >= l.w && !l.eof {
			//next read will trigger buffer reload
			//thats why we have to add buffer before it get invalidated
			if pos < l.w {
				//it means we have to write buffer
				//as next readChar will invalidate buffer
				//fmt.Println(pos, l.w, "\"", string(l.buffer[pos:l.w]), "\"", l.buffer[pos:l.w])
				s.WriteString(string(l.buffer[pos:l.w]))
			}

			pos = 0
		}

		l.readChar()
		if predicate(l.ch) {
			break
		}

	}
	if pos < l.r {
		//fmt.Println(pos, l.r, string(l.buffer[pos:l.r]))
		s.WriteString(string(l.buffer[pos:l.r]))
	}
	return s.String()
}

func (l *Lexer) readTillEOL() string {
	strx := l.readBufferConditional(
		func(ch byte) bool {
			return ch == '\n' || ch == 0
		})
	if l.PeekNextChar() == '\r' {
		//ignore
		l.readChar()
	}
	return strx
}

func (l *Lexer) readWord() string {
	return l.readBufferConditional(
		func(ch byte) bool {
			return !jacktoken.IsValidWordChar(ch)
		})
}

func (l *Lexer) NextToken() (tok jacktoken.Token, lineNumber int) {
	//var tok jacktoken.Token
	l.skipWhiteSpaces()

	if l.ch == '/' {
		if l.PeekNextChar() == '/' {
			//comment
			lineNumber = l.last_line_number
			l.readChar()
			l.readChar()
			tok.Type = jacktoken.COMMENT
			tok.Literal = l.readTillEOL()
		} else if l.PeekNextChar() == '*' {
			lineNumber = l.last_line_number
			l.readChar()
			l.readChar()
			var str strings.Builder
			tok.Type = jacktoken.MULTI_LINE_COMMENT
			for {
				strx := l.readBufferConditional(
					func(ch byte) bool {
						return ch == '*' || ch == 0
					})
				str.WriteString(strx)
				next_char := l.PeekNextChar()
				if next_char == '/' {
					l.readChar()
					tok.Literal = str.String()
					break
				} else if next_char == 0 {
					tok.Type = jacktoken.ILLEGAL
					tok.Literal = str.String()
					break
				}
			}

		} else if jacktoken.IsSymbol(l.ch) {
			lineNumber = l.last_line_number
			tok.Type = jacktoken.SYMBOL
			tok.Literal = string(l.ch)
		}
		// else {
		// 	//should not happen
		// }
	} else if l.ch == 0 {
		lineNumber = l.last_line_number
		tok.Type = jacktoken.EOF
		tok.Literal = ""
		return
	} else if jacktoken.IsSymbol(l.ch) {
		lineNumber = l.last_line_number
		tok.Type = jacktoken.SYMBOL
		tok.Literal = string(l.ch)
	} else if jacktoken.IsStringLiteralSymbol(l.ch) {
		lineNumber = l.last_line_number
		l.readChar() //ignore first
		tok.Type = jacktoken.STRING_CONSTANT
		//We allow double quote. but its not allowed in the original
		var last_char byte = 0
		tok.Literal = l.readBufferConditional(
			func(ch byte) bool {
				//allow using \" inside
				if last_char != '\\' {
					if jacktoken.IsStringLiteralSymbol(ch) {
						return true
					}
				}
				last_char = ch
				return ch == 0 || ch == '\n' || ch == '\r'
			})
		tok.Literal = strings.ReplaceAll(tok.Literal, "\\\"", "\"")
		if !jacktoken.IsStringLiteralSymbol(l.ch) {
			tok.Type = jacktoken.ILLEGAL
		}
	} else {
		lineNumber = l.last_line_number
		word := l.readWord()
		tok.Type = jacktoken.LookupType(word)
		tok.Literal = word
		return

	}
	l.readChar()
	//return tok
	return
}

func (l *Lexer) PeekNextChar() byte {
	l.readBuffer()

	if l.r+1 < l.w {
		return l.buffer[l.r+1]
	} else {
		return 0
	}
}

func (l *Lexer) PrintCharLoop() {
	for {
		fmt.Printf("'%c'", l.ch)
		l.readChar()
		if l.ch == 0 {
			break
		}
	}
}

//make it compatible with nand2tetrsi API

func NewTokenizer(input_stream io.Reader) *JackTokenizerApi {
	jAPI := JackTokenizerApi{
		lexer:           New(input_stream),
		Current:         jacktoken.Token{},
		LineNumber:      0,
		next:            jacktoken.Token{},
		next_lineNumber: 0,
	}
	jAPI.next, jAPI.next_lineNumber = jAPI.lexer.NextToken()
	return &jAPI
}

func (jAPI *JackTokenizerApi) HasMoreTokens() bool {
	return jAPI.next.Type != jacktoken.EOF
}

func (jAPI *JackTokenizerApi) Advance() {
	jAPI.Current = jAPI.next
	jAPI.LineNumber = jAPI.next_lineNumber
	jAPI.next, jAPI.next_lineNumber = jAPI.lexer.NextToken()
}

func (jAPI *JackTokenizerApi) AdvanceOnlyNext() jacktoken.Token {
	//this will be usefull for skipping comments
	tok := jAPI.next
	jAPI.next, jAPI.next_lineNumber = jAPI.lexer.NextToken()
	return tok
}

func (jAPİ *JackTokenizerApi) PeekToken() jacktoken.Token {
	token := jAPİ.next
	return token
}
