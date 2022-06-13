package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/fs"
	"io/ioutil"
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
	var err error
	switch operation {
	case "list":
		err = list(fileName, writer)
		if err != nil {
			return fmt.Errorf(err.Error())
		}
	case "add":
		item, ok := args["item"]

		if !ok || item == "" {

			return fmt.Errorf("-item flag has to be specified")
		}
		err = add(fileName, item, writer)
		if err != nil {
			return fmt.Errorf(err.Error())
		}
	case "remove":
		err = remove(args["is"], writer)
		if err != nil {
			return fmt.Errorf(err.Error())
		}
	case "findById":
		err = findById(args["id"], writer)
		if err != nil {
			return fmt.Errorf(err.Error())
		}
	default:
		return fmt.Errorf("Operation %s not allowed!", operation)
	}

	_, errf := os.Open(fileName)
	if errf != nil {
		return fmt.Errorf("!")
	}

	return nil
}

func findById(s string, writer io.Writer) error {
	panic("unimplemented")
}

func remove(s string, writer io.Writer) error {
	panic("unimplemented")
}

func add(fileName string, strItem string, writer io.Writer) error {

	b := []byte(strItem)
	var item Item
	err := json.Unmarshal(b, &item)
	if err != nil {
		return err
	}

	var filePermission fs.FileMode = 0644
	file, err := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE, filePermission)

	if err != nil {
		return err
	}

	bytes, err := ioutil.ReadAll(file)
	if err != nil {
		return err
	}

	var items []Item
	err = json.Unmarshal(bytes, &items)
	if err != nil {
		return err
	}
	for _, i := range items {

		if fmt.Sprint(item.Id) == i.Id {

			writer.Write([]byte("Item with id 1 already exists"))
			file.Close()
			return nil
		}
	}

	items = append(items, item)
	bt, err := json.Marshal(items)
	if err != nil {
		return err
	}
	writer.Write(bt)

	file.Write([]byte(bt))
	file.Close()

	return nil
}

func list(fileName string, writer io.Writer) error {
	var filePermission fs.FileMode = 0644

	file, err := os.OpenFile(fileName, os.O_RDONLY, filePermission)

	if err != nil {
		return err
	}

	bytes, err := ioutil.ReadAll(file)

	if err != nil {
		return err
	}
	err = file.Close()
	if err != nil {
		return err
	}
	writer.Write(bytes)

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
