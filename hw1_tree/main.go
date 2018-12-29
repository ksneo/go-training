package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

const treeNode = "├"
const treeNodeLast = "└"
const treeHor = "───"
const treeVert = "│"

func drawFilePath(file string, level int, isLast bool) string {
	treeNodeString := treeNode
	if isLast {
		treeNodeString = treeNodeLast
	}
	return strings.Repeat(treeVert+"\t", level) + treeNodeString + treeHor + file
}

func walkTree(out io.Writer, path string, printFiles bool, level int, drawPathString string) error {
	files, err := ioutil.ReadDir(path)
	if err != nil {
		return err
	}
	elemCount := len(files)
	for idx, file := range files {
		if file.IsDir() {
			fmt.Fprintln(out, drawPathString+file.Name())
			err = walkTree(out, filepath.Join(path, file.Name()), printFiles, level+1, drawPathString)
			if err != nil {
				return err
			}
		} else if printFiles {
			fmt.Fprintln(out, drawFilePath(file.Name(), level, idx+1 == elemCount))
		}
	}

	return nil
}

func dirTree(out io.Writer, path string, printFiles bool) error {
	drawPathString := treeNode + treeHor
	level := 0
	return walkTree(out, path, printFiles, level, drawPathString)
}

func main() {
	out := os.Stdout
	if !(len(os.Args) == 2 || len(os.Args) == 3) {
		panic("usage go run main.go . [-f]")
	}
	path := os.Args[1]
	printFiles := len(os.Args) == 3 && os.Args[2] == "-f"
	err := dirTree(out, path, printFiles)
	if err != nil {
		panic(err.Error())
	}
}
