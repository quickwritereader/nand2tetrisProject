package jacktoken

type TokenType string

type Token struct {
	Type    TokenType
	Literal string
}

const (
	ILLEGAL            = "Illegal"
	COMMENT            = "COMMENT"
	MULTI_LINE_COMMENT = "MULTI_LINE_COMMENT"
	EOF                = "EOF"
	IDENTIFIER         = "identifier"
	KEYWORD            = "keyword"
	SYMBOL             = "symbol"
	INTEGER_CONSTANT   = "integerConstant"
	STRING_CONSTANT    = "stringConstant"
)

var valid_symbols = []byte("{}()[],.;=+-/*><&|~")

var token_map = map[string]TokenType{
	"class":       KEYWORD,
	"constructor": KEYWORD,
	"function":    KEYWORD,
	"method":      KEYWORD,
	"field":       KEYWORD,
	"static":      KEYWORD,
	"var":         KEYWORD,
	"int":         KEYWORD,
	"char":        KEYWORD,
	"boolean":     KEYWORD,
	"void":        KEYWORD,
	"true":        KEYWORD,
	"false":       KEYWORD,
	"null":        KEYWORD,
	"this":        KEYWORD,
	"let":         KEYWORD,
	"do":          KEYWORD,
	"if":          KEYWORD,
	"else":        KEYWORD,
	"while":       KEYWORD,
	"return":      KEYWORD,
}

func isDigit(ch byte) bool {
	return ch >= '0' && ch <= '9'
}

func isLetter(ch byte) bool {
	return (ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z')
}

func isValidIdentifier(str string) bool {
	if len(str) < 1 {
		return false
	}
	if isLetter(str[0]) || str[0] == '_' {
		for _, ch := range str {
			if !isDigit(byte(ch)) && !isLetter(byte(ch)) && ch != '_' {
				return false
			}
		}
		return true
	}

	return false
}

func isNumber(str string) bool {
	for _, ch := range str {
		if !isDigit(byte(ch)) {
			return false
		}
	}
	return true
}

func IsStringLiteralSymbol(ch byte) bool {
	return ch == '"'
}

func IsSymbol(ch byte) bool {
	for _, x := range valid_symbols {
		if ch == x {
			return true
		}
	}
	return false
}

func IsValidWordChar(ch byte) bool {
	if isDigit(ch) || isLetter(ch) || ch == '_' || ch == '"' {
		return true
	}
	return false
}

func LookupType(identifier string) TokenType {
	if isNumber(identifier) {
		return INTEGER_CONSTANT
	} else if tok, ok := token_map[identifier]; ok {
		return tok
	} else if isValidIdentifier(identifier) {
		return IDENTIFIER
	}

	return ILLEGAL
}
