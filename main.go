package main

import (
	"flag"
	"fmt"
	"io"
	"os"
)

type Arguments map[string]string

type Item struct {
	Id    string
	Email string
	Age   int
}

func Perform(args Arguments, writer io.Writer) error {

	operation, ok := args["operation"]

	if !ok || operation == "" {

		return fmt.Errorf("-operation flag has to be specified")
	}

	fileName, ok := args["fileName"]

	if !ok || fileName == "" {
		return fmt.Errorf("-fileName flag has to be specified")
	}
	

	return nil
}



func parseArgs() Arguments {
	id := flag.Int("id", 0, "help")
	fileName := flag.String("fileName", "", "help")
	operation := flag.String("operation", "", "help")
	item := flag.String("item", "", "help")
	flag.Parse()
	args := Arguments{
		"id":        fmt.Sprint(*id),
		"operation": *operation,
		"item":      *item,
		"fileName":  *fileName,
	}
	return args
}

func main() {
	err := Perform(parseArgs(), os.Stdout)
	if err != nil {
		panic(err)
	}
}
