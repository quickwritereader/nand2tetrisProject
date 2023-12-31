package main

import (
	"flag"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"

	jacklexer "github.com/quickwritereader/JackSyntaxAnalyser/JackLexer"
	jackhelper "github.com/quickwritereader/JackSyntaxAnalyser/jackHelper"
	"github.com/quickwritereader/JackSyntaxAnalyser/jacktoken"
	"github.com/quickwritereader/JackSyntaxAnalyser/parser"
)

func tokenizeToXmlFile(path_list []string, commentFlag bool) {

	if len(path_list) < 1 {
		return
	}
	for _, open_file := range path_list {
		if len(open_file) > 0 {
			file, err := os.Open(open_file)
			if err != nil {
				fmt.Fprintln(os.Stderr, "Error:", err)
				return
			}
			defer file.Close()
			output_path := strings.TrimSuffix(open_file, filepath.Ext(open_file))
			writeOutput, err := os.OpenFile(output_path+".TOKEN.xml", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, os.FileMode(uint32(644)))
			if err != nil {
				fmt.Fprintln(os.Stderr, "Error:", err)
				return
			}
			defer writeOutput.Close()
			fname := path.Base(open_file)
			writeOutput.WriteString("<tokens>\n")
			tokenAPI := jacklexer.NewTokenizer(file)
			for tokenAPI.HasMoreTokens() {

				tokenAPI.Advance()

				tok, ln := tokenAPI.Current, tokenAPI.LineNumber
				if tok.Type == jacktoken.ILLEGAL {
					fmt.Fprintf(os.Stderr, "%s:%d Error \n", fname, ln)
					break
				}
				if tok.Type == jacktoken.COMMENT || tok.Type == jacktoken.MULTI_LINE_COMMENT {
					if commentFlag {
						jackhelper.WriteTokenAsXmlComment(writeOutput, tok)
					}
					continue
				}
				jackhelper.WriteTokenAsXmlNode(writeOutput, tok)

			}

			writeOutput.WriteString("</tokens>")
		}
	}

}

func parseToXmlFile(path_list []string, commentFlag bool) {

	if len(path_list) < 1 {
		return
	}
	for _, open_file := range path_list {
		if len(open_file) > 0 {
			file, err := os.Open(open_file)
			if err != nil {
				fmt.Fprintln(os.Stderr, "Error:", err)
				return
			}
			defer file.Close()
			output_path := strings.TrimSuffix(open_file, filepath.Ext(open_file))
			writeOutput, err := os.OpenFile(output_path+".PARSE.xml", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, os.FileMode(uint32(644)))
			if err != nil {
				fmt.Fprintln(os.Stderr, "Error:", err)
				return
			}
			defer writeOutput.Close()

			p := parser.NewRecursiveDescentParser(file, writeOutput)

			p.OutputComments = commentFlag
			p.CompileAsXml()

		}
	}

}

func main() {
	tokenizeFlag := flag.Bool("lexical", false, "tokenize and output xml")
	parseFlag := flag.Bool("syntactic", false, "parse and output xml")
	commentFlag := flag.Bool("syntactic", false, "output comments in xml")
	flag.Parse()

	path_list := jackhelper.GetInputFiles(flag.Arg(0), ".jack")

	if *tokenizeFlag {
		tokenizeToXmlFile(path_list, *commentFlag)
	}

	if *parseFlag {
		parseToXmlFile(path_list, *commentFlag)
	}

}
