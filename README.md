# Trabalho 

Alunos: 
 
 - Vinicius Clemente de Sousa Reis, vinicius.clemente@ufu.br
 
 - Natan Rodovalho, natan.rodovalho@ufu.br
 
 - Alexandre Pereira Marcos, alexandrepm2810@ufu.br
 

Para decifrar os arquivos, criamos um dicionário de acordo com os seguintes 
assuntos: star wars, simpsons e twin peaks. Para criarmos o dicionário
foi feita uma pesquisa sobre cidades, personagens e objetos sobre tais assuntos.

As palavras encontradas na pesquisa foram armazenadas no arquivo **base-dictionary.txt**

Para gerar um dicionário mais abrangente foi criado um código em python:

```python
import itertools


def main():
    base_dictionary = open("base-dictionary.txt", "r")
    words = base_dictionary.read().split("\n")

    dictionary = set({})

    modifiers = generate_modifiers()

    for word in words:
        for i in itertools.product([False, True], repeat=len(modifiers)):
            new_word = word
            for j in range(len(modifiers)):
                if i[j]:
                    new_word = modifiers[j](new_word)

            dictionary.add(new_word + "\n")

    final_dictionary = open("dictionary.txt", "w")
    for w in dictionary:
        if w != "":
            final_dictionary.write(w)

    final_dictionary.close()
    base_dictionary.close()


def generate_modifiers():
    upper = (lambda w: w.upper())
    lower = (lambda w: w.lower())
    capitalize = (lambda w: w.capitalize())
    replace_dash = (lambda w: w.replace("-", " "))
    replace_spaces = (lambda w: w.replace(" ", ""))
    replace_a = (lambda w: w.replace("a", "@").replace("A", "@"))
    replace_e = (lambda w: w.replace("e", "3").replace("E", "3"))
    replace_i = (lambda w: w.replace("i", "!").replace("I", "!"))
    to_pascal_case = (lambda word: ''.join(w for w in word.title() if not w.isspace()))
    modifiers = [upper, lower, capitalize, replace_dash, replace_spaces, replace_a, replace_e, replace_i, to_pascal_case]
    return modifiers


if __name__ == "__main__":
    main()

```

nesse código criamos os modificadores: 
 
 - upper: transforma a palavra em maiúscula
 
 - lower: transforma a palavra em minúscula
 
 - capitalize: transforma a primeira letra da palavra em maiúscula
 
 - replace_dash: remove o carácter '-'
 
 - replace_spaces: remove os espaços
 
 - replace_a: troca 'a' por '@' 
 
 - replace_e: troca 'e' por '3'
 
 - replace_i: troca 'i' por '!'
 
Após criar os modificadores, aplicamos esses para cada palavra do arquivo base-dictionary.txt 
para fazer a combinação dos modificadores, utilizamos a função nativa do python `itertools.product`
algumas combinações vão gerar o mesmo resultado e para não ter um dicionário 
com palavras repetidas salvamos o resultado em um set.

Depois de salvar todas as combinações dos modificadores no set, transformamos o set
no arquivo dictionary.txt

Após gerar o arquivo dictionary.txt com o comando `python generate_dictionary.py`
Utilizamos agora o código escrito em go:

```go
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

```

Nesse código:

 1. lemos as `words` do arquivo dictionary.txt
 
 2. listamos os `files` da pasta arquivos
 
 3. Criamos um goroutine para cada file.
 
 4. Cada goroutine tem o objetivo de tenta decodificar o seu arquivo utilizando as `words`
 

Utilizamos goroutine para rodar as decodificações em paralelo aumentando assim
a performance do nosso código.

Para tentar decodificar o arquivo utilizamos a função `func decode(file string, pass string) ([]byte, error)`
essa função chama o `openssl` com os argumentos `enc -d -aes-256-cbc -pbkdf2 -salt -in $file -pass pass:$pass`
caso o `openssl` retorne algum resultado verrificamos se existe algum 
caracter invalido, caso exista desconsideramos esse resultado.
Caso o decode não retorne erro salvamos no arquivo result.txt a msg:
`File:%s -> Pass:%s Data: %s\n`


Para rodar o arquivo em go é só executar `go run main.go`

Após 4h e 30min podemos verrificar que o arquivo de result.txt apresentou:

```
File:file0.enc -> Pass:teste Data: teste

File:file15.enc -> Pass:dagobah Data: teste

File:file5.enc -> Pass:b@rt Data: teste

File:file7.enc -> Pass:davidlynch Data: teste

File:file8.enc -> Pass:TwinPeaks Data: teste

File:file10.enc -> Pass:burns Data: teste

File:file27.enc -> Pass:Springfield Data: teste

File:file12.enc -> Pass:DaleCooper Data: teste

File:file29.enc -> Pass:falcon Data: teste

File:file2.enc -> Pass:laurapalmer Data: teste

```

Como devemos escolher apenas 5 palavras para obter os pontos de exclusividades
escolhemos:

```
File:file5.enc -> Pass:b@rt Data: teste
File:file29.enc -> Pass:falcon Data: teste
File:file12.enc -> Pass:DaleCooper Data: teste
File:file2.enc -> Pass:laurapalmer Data: teste
File:file15.enc -> Pass:dagobah Data: teste
```

o código completo pode ser encontrado em 
`https://github.com/viniciusCSreis/ufu-openssl-dicionary`