package jacklexer

import (
	"strings"
	"testing"

	jackhelper "github.com/quickwritereader/JackSyntaxAnalyser/jackHelper"
	"github.com/quickwritereader/JackSyntaxAnalyser/jacktoken"
)

func TestNextTokenAddition(t *testing.T) {
	// We allow double quote unlike the original
	var test_code = `
	main()
	"s \" s"
	(){}
	/* multi
	comment
	test
	*/
	/* multi
	illegal

	`

	buffer_size := 4096
	myReader1 := strings.NewReader(test_code)
	lexer1 := NewBufferSize(myReader1, buffer_size)

	expectedTokList := [...]jacktoken.Token{
		{Type: "identifier", Literal: "main"},
		{Type: "symbol", Literal: "("},
		{Type: "symbol", Literal: ")"},
		{Type: "stringConstant", Literal: "s \" s"},
		{Type: "symbol", Literal: "("},
		{Type: "symbol", Literal: ")"},
		{Type: "symbol", Literal: "{"},
		{Type: "symbol", Literal: "}"},
		{Type: "MULTI_LINE_COMMENT", Literal: " multi\n\tcomment\n\ttest\n\t"},
		{Type: "Illegal", Literal: " multi\n\tillegal\n\n\t"},
	}

	for _, expTok := range expectedTokList {
		tok, _ := lexer1.NextToken()
		if tok != expTok {
			t.Errorf("Actual: %#v Expected %#v", tok, expTok)
		}

		if tok.Type == jacktoken.EOF {
			break
		}

	}
	tok, _ := lexer1.NextToken()
	if tok.Type != jacktoken.EOF {
		t.Errorf("Last token should be EOF, %#v", tok)
	}

}

func TestTokenizerApi(t *testing.T) {
	var str = `if (x < 153)
	{let city="Paris";}`

	tokenApi := NewTokenizer(strings.NewReader(str))

	writer := strings.Builder{}
	for tokenApi.HasMoreTokens() {
		tokenApi.Advance()
		jackhelper.WriteTokenAsXmlNode(&writer, tokenApi.Current)

	}
	expected := `<keyword> if </keyword>
				<symbol> ( </symbol>
				<identifier> x </identifier>
				<symbol> &lt; </symbol>
				<integerConstant> 153 </integerConstant>
				<symbol> ) </symbol>
				<symbol> { </symbol>
				<keyword> let </keyword>
				<identifier> city </identifier>
				<symbol> = </symbol>
				<stringConstant> Paris </stringConstant>
				<symbol> ; </symbol>
				<symbol> } </symbol>`

	isOk, msg := jackhelper.CompareStringLines(writer.String(), expected)

	if !isOk {
		t.Error(msg)
	}

}

func TestTokenizerApi2(t *testing.T) {
	var str = `
	// (derived from projects/09/Square/Main.jack, with testing additions)

	/** Initializes a new Square Dance game and starts running it. */
	class Main {
		static boolean test;    // Added for testing -- there is no static keyword
								// in the Square files.
		function void main() {
		  var SquareGame game;
		  let game = SquareGame.new();
		  do game.run();
		  do game.dispose();
		  return;
		}
	
		function void more() {  // Added to test Jack syntax that is not used in
			var int i, j;       // the Square files.
			var String s;
			var Array a;
			if (false) {
				let s = "string constant";
				let s = null;
				let a[1] = a[2];
			}
			else {              // There is no else keyword in the Square files.
				let i = i * (-j);
				let j = j / (-2);   // note: unary negate constant 2
				let i = i | j;
			}
			return;
		}
	}	
	
	`

	tokenApi := NewTokenizer(strings.NewReader(str))

	writer := strings.Builder{}
	jackhelper.WriteStringAsBeginNode(&writer, "tokens")
	for tokenApi.HasMoreTokens() {
		tokenApi.Advance()
		if tokenApi.Current.Type != jacktoken.COMMENT && tokenApi.Current.Type != jacktoken.MULTI_LINE_COMMENT {
			jackhelper.WriteTokenAsXmlNode(&writer, tokenApi.Current)
		}
	}
	jackhelper.WriteStringAsEndNode(&writer, "tokens")
	expected := `
			<tokens>
			<keyword> class </keyword>
			<identifier> Main </identifier>
			<symbol> { </symbol>
			<keyword> static </keyword>
			<keyword> boolean </keyword>
			<identifier> test </identifier>
			<symbol> ; </symbol>
			<keyword> function </keyword>
			<keyword> void </keyword>
			<identifier> main </identifier>
			<symbol> ( </symbol>
			<symbol> ) </symbol>
			<symbol> { </symbol>
			<keyword> var </keyword>
			<identifier> SquareGame </identifier>
			<identifier> game </identifier>
			<symbol> ; </symbol>
			<keyword> let </keyword>
			<identifier> game </identifier>
			<symbol> = </symbol>
			<identifier> SquareGame </identifier>
			<symbol> . </symbol>
			<identifier> new </identifier>
			<symbol> ( </symbol>
			<symbol> ) </symbol>
			<symbol> ; </symbol>
			<keyword> do </keyword>
			<identifier> game </identifier>
			<symbol> . </symbol>
			<identifier> run </identifier>
			<symbol> ( </symbol>
			<symbol> ) </symbol>
			<symbol> ; </symbol>
			<keyword> do </keyword>
			<identifier> game </identifier>
			<symbol> . </symbol>
			<identifier> dispose </identifier>
			<symbol> ( </symbol>
			<symbol> ) </symbol>
			<symbol> ; </symbol>
			<keyword> return </keyword>
			<symbol> ; </symbol>
			<symbol> } </symbol>
			<keyword> function </keyword>
			<keyword> void </keyword>
			<identifier> more </identifier>
			<symbol> ( </symbol>
			<symbol> ) </symbol>
			<symbol> { </symbol>
			<keyword> var </keyword>
			<keyword> int </keyword>
			<identifier> i </identifier>
			<symbol> , </symbol>
			<identifier> j </identifier>
			<symbol> ; </symbol>
			<keyword> var </keyword>
			<identifier> String </identifier>
			<identifier> s </identifier>
			<symbol> ; </symbol>
			<keyword> var </keyword>
			<identifier> Array </identifier>
			<identifier> a </identifier>
			<symbol> ; </symbol>
			<keyword> if </keyword>
			<symbol> ( </symbol>
			<keyword> false </keyword>
			<symbol> ) </symbol>
			<symbol> { </symbol>
			<keyword> let </keyword>
			<identifier> s </identifier>
			<symbol> = </symbol>
			<stringConstant> string constant </stringConstant>
			<symbol> ; </symbol>
			<keyword> let </keyword>
			<identifier> s </identifier>
			<symbol> = </symbol>
			<keyword> null </keyword>
			<symbol> ; </symbol>
			<keyword> let </keyword>
			<identifier> a </identifier>
			<symbol> [ </symbol>
			<integerConstant> 1 </integerConstant>
			<symbol> ] </symbol>
			<symbol> = </symbol>
			<identifier> a </identifier>
			<symbol> [ </symbol>
			<integerConstant> 2 </integerConstant>
			<symbol> ] </symbol>
			<symbol> ; </symbol>
			<symbol> } </symbol>
			<keyword> else </keyword>
			<symbol> { </symbol>
			<keyword> let </keyword>
			<identifier> i </identifier>
			<symbol> = </symbol>
			<identifier> i </identifier>
			<symbol> * </symbol>
			<symbol> ( </symbol>
			<symbol> - </symbol>
			<identifier> j </identifier>
			<symbol> ) </symbol>
			<symbol> ; </symbol>
			<keyword> let </keyword>
			<identifier> j </identifier>
			<symbol> = </symbol>
			<identifier> j </identifier>
			<symbol> / </symbol>
			<symbol> ( </symbol>
			<symbol> - </symbol>
			<integerConstant> 2 </integerConstant>
			<symbol> ) </symbol>
			<symbol> ; </symbol>
			<keyword> let </keyword>
			<identifier> i </identifier>
			<symbol> = </symbol>
			<identifier> i </identifier>
			<symbol> | </symbol>
			<identifier> j </identifier>
			<symbol> ; </symbol>
			<symbol> } </symbol>
			<keyword> return </keyword>
			<symbol> ; </symbol>
			<symbol> } </symbol>
			<symbol> } </symbol>
			</tokens>
	`

	isOk, msg := jackhelper.CompareStringLines(writer.String(), expected)

	if !isOk {
		t.Error(msg)
	}

}
