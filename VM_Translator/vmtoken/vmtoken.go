package vmtoken

type TokenType string

type Token struct {
	Type    TokenType
	Literal string
}

const (
	ILLEGAL    = "Illegal"
	EOF        = "EOF"
	PUSH       = "PUSH"
	POP        = "POP"
	ARITHMETIC = "ARITHMETIC"
	STATIC     = "STATIC"
	LOCAL      = "LOCAL"
	CONSTANT   = "CONSTANT"
	THIS       = "THIS"
	THAT       = "THAT"
	TEMP       = "TEMP"
	ARGUMENT   = "ARGUMENT"
	COMMENT    = "COMMENT"
	NUMBER     = "NUMBER"
	POINTER    = "POINTER"
	GOTO       = "GOTO"
	LABEL      = "LABEL"
	IF_GOTO    = "IF-GOTO"
	RETURN     = "RETURN"
	CALL       = "CALL"
	FUNCTION   = "FUNCTION"
	IDENTIFIER = "IDENTIFIER"
)

var keywords = map[string]TokenType{
	"push":     PUSH,
	"pop":      POP,
	"static":   STATIC,
	"local":    LOCAL,
	"this":     THIS,
	"that":     THAT,
	"argument": ARGUMENT,
	"pointer":  POINTER,
	"temp":     TEMP,
	"constant": CONSTANT,
	"add":      ARITHMETIC,
	"sub":      ARITHMETIC,
	"or":       ARITHMETIC,
	"and":      ARITHMETIC,
	"neg":      ARITHMETIC,
	"not":      ARITHMETIC,
	"lt":       ARITHMETIC,
	"gt":       ARITHMETIC,
	"eq":       ARITHMETIC,
	"goto":     GOTO,
	"if-goto":  IF_GOTO,
	"label":    LABEL,
	"function": FUNCTION,
	"return":   RETURN,
	"call":     CALL,
}

func LookupType(identifier string) TokenType {
	if tok, ok := keywords[identifier]; ok {
		return tok
	}
	return IDENTIFIER
}
