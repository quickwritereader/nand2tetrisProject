package vmlexer

import (
	"fmt"
	"io"
	"strings"

	"github.com/quickwritereader/vmTranslator/vmtoken"
)

const (
	BUFFER_SIZE      = 4096
	MIN_VALID_BUFFER = 1
)

// func init() {
// 	print(vmtoken.ILLEGAL + "ok\n")
// }

type Lexer struct {
	istream      io.Reader
	buffer       []byte
	r, w         int // buf read and write positions
	eof          bool
	positionBase int
	ch           byte //current character
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
		l.positionBase += l.w
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
	return l.readBufferConditional(
		func(ch byte) bool {
			return ch == '\n' || ch == '\r' || ch == 0
		})
}

func isDigit(ch byte) bool {
	return ch >= '0' && ch <= '9'
}

func isLetter(ch byte) bool {
	return (ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z')
}

func isValidChars(ch byte) bool {
	return isDigit(ch) || isLetter(ch) || ch == '_' || ch == '-' || ch == '.'
}

func (l *Lexer) readWord() string {
	return l.readBufferConditional(
		func(ch byte) bool {
			return !isValidChars(ch)
		})
}

func (l *Lexer) NextToken() vmtoken.Token {
	var tok vmtoken.Token
	l.skipWhiteSpaces()

	if l.ch == '/' && l.PeekNextChar() == '/' {
		//comment
		//if PeekNextChar causes buffer invalidation, then l.r is -1
		tok.Type = vmtoken.COMMENT
		tok.Literal = l.readTillEOL()
	} else if l.ch == 0 {
		tok.Type = vmtoken.EOF
		tok.Literal = ""
		return tok
	} else {
		word := l.readWord()
		if len(word) > 0 {
			if isNumber(word) {
				tok.Type = vmtoken.NUMBER
				tok.Literal = word
			} else {
				tok.Type = vmtoken.LookupType(word)
				tok.Literal = word
			}
		} else {
			tok.Type = vmtoken.ILLEGAL
			tok.Literal = "\"" + string(l.ch) + "\""
		}
		//in the word case, the next char was alreade read
		//so skip l.readChar advance
		return tok

	}
	//advance Position
	l.readChar()
	return tok
}

func isNumber(str string) bool {
	for _, ch := range str {
		if !isDigit(byte(ch)) {
			return false
		}
	}
	return true
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
