package main

import (
	"fmt"
	"os"
	"bufio"
	"strings"
)

var (
	hadError = false
)


func main(){
	args := os.Args
	if len(args) > 2 {
		fmt.Println("Usage: glox [script]")
	}else if len(args) == 2{
		runFile(args[1])
	}else{
		runPrompt()
	}
}

func runFile(fileName string){
	file, err := os.ReadFile(fileName)
	if err!=nil{
		panic(err)
	}
	source := string(file)
	run(source)
	
}

func runPrompt(){
	for{
		fmt.Println("> ")
		inputReader := bufio.NewReader(os.Stdin)
		input, _ := inputReader.ReadString('\n')
		if input == ""{
			fmt.Println("Exiting")
			break
		}
		run(input)
	}
}

func run(source string){
	scanner := bufio.NewScanner(strings.NewReader(source))
	scanner.Split(bufio.ScanWords)

	for scanner.Scan(){
		scan := NewScanner(scanner.Text())
		fmt.Println(scan)
	}
}

func report(line int, where string, message string){
	fmt.Println("[line %d] Error %s : %s",line,where,message)
}

func error(line int, message string){
	report(line, "", message)
}

