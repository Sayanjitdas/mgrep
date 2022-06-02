package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

type ParseInfo struct {
	searchString string
	pathOfFile   string
}

var SearchList []ParseInfo

func (prse ParseInfo) ParseFile(wt *sync.WaitGroup) {

	f, err := os.Open(prse.pathOfFile)
	checkErr(err)
	defer f.Close()
	defer wt.Done()
	lineNo := 1
	scnner := bufio.NewScanner(f)

	for scnner.Scan() {
		if strings.Contains(scnner.Text(), prse.searchString) {
			_, fileName := filepath.Split(prse.pathOfFile)
			fmt.Printf(">> file: %20s | lineNo: %5d | matchedline: %s\n", fileName, lineNo, scnner.Text())
		}
		lineNo += 1
	}
}

func checkErr(err error) {
	if err != nil {
		if strings.Contains(err.Error(), "open dir: no such file or directory") {
			fmt.Println(err.Error())
			os.Exit(1)
		} else {
			panic(err.Error())
		}
	}
}

func ParseFileRunner() {
	var wt sync.WaitGroup
	for _, fileInfo := range SearchList {
		wt.Add(1)
		go fileInfo.ParseFile(&wt)
	}
	wt.Wait()
}

func WalkDirectoryAndFiles(searchString, directoryName string) {

	f, err := os.Open(directoryName)
	checkErr(err)
	defer f.Close()

	files, err := f.ReadDir(0)
	checkErr(err)

	for _, v := range files {

		if v.IsDir() {
			WalkDirectoryAndFiles(searchString, filepath.Join(directoryName, v.Name()))
		} else {
			p := ParseInfo{
				searchString: searchString,
				pathOfFile:   filepath.Join(directoryName, v.Name()),
			}

			SearchList = append(SearchList, p)
		}
	}
}

func main() {

	//receiving command line arguments
	cmdArgs := os.Args
	var searchString, directoryName string
	if len(cmdArgs) == 3 {
		searchString = cmdArgs[1]
		directoryName = cmdArgs[2]
	} else {
		fmt.Printf(`ERROR: required two args given %d  
USAGE: mgrep <"search_string"> <search_dir(absolute-path)>`, len(cmdArgs)-1)
		fmt.Println()
		os.Exit(1)
	}

	WalkDirectoryAndFiles(searchString, directoryName)
	ParseFileRunner()

}
