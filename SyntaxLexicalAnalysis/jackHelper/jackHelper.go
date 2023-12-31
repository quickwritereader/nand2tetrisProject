package jackhelper

import (
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/quickwritereader/JackSyntaxAnalyser/jacktoken"
)

type DummyWriter struct {
}

func (d *DummyWriter) WriteString(s string) (n int, err error) {
	return len(s), nil
}

func GetInputFiles(path string, ext string) []string {
	fileInfo, err := os.Stat(path)
	path_list := []string{}
	if err == nil {
		if fileInfo.IsDir() {

			filepath.WalkDir(path, func(inner_path string, d fs.DirEntry, err error) error {
				if !d.IsDir() {
					if strings.HasSuffix(inner_path, ext) {
						path_list = append(path_list, inner_path)
					}
				}
				return nil
			})

		} else {
			path_list = append(path_list, path)
		}

	}
	return path_list
}

func encodeXml(in string) string {
	in = strings.ReplaceAll(in, "&", "&amp;")
	in = strings.ReplaceAll(in, ">", "&gt;")
	in = strings.ReplaceAll(in, "<", "&lt;")
	in = strings.ReplaceAll(in, "\"", "&quote;")
	return in
}

func WriteTokenAsXmlNode(writeOutput io.StringWriter, tok jacktoken.Token) {
	tok_type := string(tok.Type)
	writeOutput.WriteString("<")
	writeOutput.WriteString(tok_type)
	writeOutput.WriteString("> ")
	writeOutput.WriteString(encodeXml(tok.Literal))
	writeOutput.WriteString(" </")
	writeOutput.WriteString(tok_type)
	writeOutput.WriteString(">\n")
}

func WriteTokenAsXmlComment(writeOutput io.StringWriter, tok jacktoken.Token) {
	writeOutput.WriteString("<!--")
	tok_type := string(tok.Type)
	writeOutput.WriteString("<")
	writeOutput.WriteString(tok_type)
	writeOutput.WriteString("> ")
	writeOutput.WriteString(encodeXml(tok.Literal))
	writeOutput.WriteString(" </")
	writeOutput.WriteString(tok_type)
	writeOutput.WriteString("> -->\n")

}

func WriteStringAsBeginNode(writeOutput io.StringWriter, str string) {
	writeOutput.WriteString("<")
	writeOutput.WriteString(str)
	writeOutput.WriteString(">\n")
}
func WriteStringAsEndNode(writeOutput io.StringWriter, str string) {
	writeOutput.WriteString("</")
	writeOutput.WriteString(str)
	writeOutput.WriteString(">\n")
}

func returnErr(msg string, prev error) error {
	if prev != nil {
		return fmt.Errorf(msg+" (%w )", prev)
	} else {
		return fmt.Errorf(msg)
	}
}

func LocalIndex[S ~[]E, E comparable](s S, v E) int {
	for i := range s {
		if v == s[i] {
			return i
		}
	}
	return -1
}

func CompareStringLines(actual string, expected string) (bool, string) {
	strArr := strings.Split(strings.TrimSpace(actual), "\n")
	expectedArr := strings.Split(strings.TrimSpace(expected), "\n")
	i, j := 0, 0
	exp := ""
	act := ""
	for {

		if i < len(strArr) {
			act = strings.TrimSpace(strArr[i])
			if len(act) < 1 {
				//empty case, increment
				i++
				continue
			}
		} else {
			act = ""
		}
		if j < len(expectedArr) {
			exp = strings.TrimSpace(expectedArr[j])
			if len(exp) < 1 {
				//empty case, increment
				j++
				continue
			}
		} else {
			exp = ""
		}
		if act != exp {
			return false, fmt.Sprintf("Actual: \"%s\" Expected: \"%s\"", act, exp)
		}

		if i >= len(strArr) && j >= len(expectedArr) {
			//both ended
			break
		}

		i++
		j++
	}
	return true, ""
}
