package main

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/quickwritereader/vmTranslator/codegen"
	"github.com/quickwritereader/vmTranslator/parser"
)

func processFile(cg *codegen.CodeGen, open_file string) {
	if len(open_file) > 0 {
		file, err := os.Open(open_file)
		if err != nil {
			fmt.Println("Error:", err)
			return
		}
		defer file.Close()
		cg.SetName(filepath.Base(strings.TrimSuffix(open_file, filepath.Ext(open_file))))
		fmt.Println("Process file:", open_file)
		pp := parser.NewParser(file)
	exit:
		for {
			cmd, err := pp.NextCommand()
			if err != nil {
				fmt.Printf("Err %s", err.Error())
				break exit
			}
			switch cmd.(type) {
			case *codegen.CC_EMPTY:
				break exit

			}
			cg.Process(cmd)
		}

	}
}

func getInputFiles(path string, ext string) []string {
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

func main() {
	var path string
	if len(os.Args) > 1 {
		path = os.Args[1]
	}
	cg := codegen.NewCodeGen()

	fileInfo, err := os.Stat(path)

	if err == nil {
		output_path := path
		if !fileInfo.IsDir() {
			output_path = strings.TrimSuffix(path, filepath.Ext(path))
		} else {
			output_path = output_path + string(os.PathSeparator) + filepath.Base(path)
		}
		fmt.Print("Output: ")
		fmt.Println(output_path + ".asm")
		writeOutput, err := os.OpenFile(output_path+".asm", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, os.FileMode(uint32(644)))
		if err != nil {
			fmt.Println("Error:", err)
			return
		}
		defer writeOutput.Close()

		cg.SetWriter(writeOutput)

		path_list := getInputFiles(path, ".vm")

		if len(path_list) > 1 {
			//multiple vm files
			cg.WriteInit()
		}
		for _, file_path := range path_list {
			processFile(&cg, file_path)
		}
	}

}
