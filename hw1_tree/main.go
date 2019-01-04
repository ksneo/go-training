package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
)

const treeNode = "├"
const treeNodeLast = "└"
const treeHor = "───"
const treeVert = "│"

func filter(vs []os.FileInfo, printFiles bool) []os.FileInfo {
	vsf := make([]os.FileInfo, 0)
	for _, v := range vs {
		if !v.IsDir() && !printFiles {
			continue
		}
		vsf = append(vsf, v)
	}
	return vsf
}

func drawFilePath(file string, level int, isLast bool, drawPathString string) string {
	treeNodeString := treeNode
	if isLast {
		treeNodeString = treeNodeLast
	}
	return drawPathString + treeNodeString + treeHor + file
}

func walkTree(out io.Writer, path string, printFiles bool, level int, drawPathString string) error {
	files, err := ioutil.ReadDir(path)
	files = filter(files, printFiles)
	sort.Slice(files, func(i, j int) bool { return files[i].Name() < files[j].Name() })
	curPathString := ""
	newPathString := ""
	if err != nil {
		return err
	}
	elemCount := len(files)
	for idx, file := range files {
		if file.IsDir() {
			if idx+1 == elemCount {
				curPathString = treeNodeLast
			} else {
				curPathString = treeNode
			}
			curPathString = curPathString + treeHor
			fmt.Fprintln(out, drawPathString+curPathString+file.Name()+"	")
			if idx+1 == elemCount {
				newPathString = drawPathString + " 	"
			} else {
				newPathString = drawPathString + treeVert + "	"
			}
			err = walkTree(out, filepath.Join(path, file.Name()), printFiles, level+1, newPathString)
			if err != nil {
				return err
			}
		} else if printFiles {
			fmt.Fprintln(out, drawFilePath(file.Name(), level, idx+1 == elemCount, drawPathString))
		}
	}
	return nil
}

func dirTree(out io.Writer, path string, printFiles bool) error {
	drawPathString := ""
	level := 0
	walkTree(out, path, printFiles, level, drawPathString)
	return nil
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
