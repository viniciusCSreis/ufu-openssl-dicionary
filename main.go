package main

import (
	"bufio"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

const (
	filesDir       = "arquivos"
	dictionaryFile = "dictionary.txt"
	resultFilePath = "result.txt"
)

var resultFile *os.File

func main() {

	createResultFile()
	defer resultFile.Close()

	files := getFilesToDecode()
	words := getWords()

	//Create a wait group to wait goroutines
	wgFile := sync.WaitGroup{}
	wgFile.Add(len(files))

	for _, f := range files {
		fProxy := f
		go func() {
			for n, word := range words {

				if n%100 == 0 {
					fmt.Printf("file [%d]: %s wordLine %d\n", time.Now().Unix(), fProxy.Name(), n)
				}

				data, err := decode(filepath.Join(filesDir, fProxy.Name()), word)
				if err == nil {
					data := fmt.Sprintf("File:%s -> Pass:%s Data: %s\n", fProxy.Name(), word, string(data))
					WriteToFile(data)
				}

			}
			//Finish goroutine
			wgFile.Done()
		}()
	}

	//wait goroutines
	wgFile.Wait()

}

func getWords() []string {
	wordsBytes, err := ioutil.ReadFile(dictionaryFile)
	if err != nil {
		panic(err.Error())
	}
	words := strings.Split(string(wordsBytes), "\n")
	return words
}

func getFilesToDecode() []os.FileInfo {
	files, err := ioutil.ReadDir(filesDir)
	if err != nil {
		panic(err.Error())
	}
	return files
}

func createResultFile() {
	var err error
	resultFile, err = os.OpenFile(resultFilePath, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
	if err != nil {
		resultFile, err = os.Create(resultFilePath)
		if err != nil {
			panic(err)
		}
	}
}

var mu sync.Mutex

func WriteToFile(data string) {
	mu.Lock()
	defer mu.Unlock()
	w := bufio.NewWriter(resultFile)

	if _, err := w.Write([]byte(data)); err != nil {
		panic(err)
	}

	if err := w.Flush(); err != nil {
		panic(err)
	}
}

func decode(file string, pass string) ([]byte, error) {
	args := []string{
		"enc",
		"-d",
		"-aes-256-cbc",
		"-pbkdf2",
		"-salt",
		"-in",
		file,
		"-pass",
		fmt.Sprintf("pass:%s", pass),
	}
	cmd := exec.Command("openssl", args...)

	data, err := cmd.Output()

	if err != nil {
		return nil, err
	}

	for _, d := range data {
		if (d < ' ' || d > '~') && d != '\n' {
			return nil, errors.New("decode to a not valid data")
		}
	}
	return data, nil

}
