package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/fs"
	"io/ioutil"
	"os"
	"strings"
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
		err = remove(fileName,args["id"], writer)
		if err != nil {
			return fmt.Errorf(err.Error())
		}
	case "findById":
		err = findById(fileName,args["id"], writer)
		if err != nil {
			return fmt.Errorf(err.Error())
		}
	default:
		return fmt.Errorf("Operation %s not allowed!", operation)
	}	

	return nil
}

func findById(fileName string,id string, writer io.Writer) error {
	if id == "" {
		return fmt.Errorf("-id flag has to be specified")
	}

    var filePermission fs.FileMode = 0644

	file, err := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE, filePermission)
	if err != nil {
		return err
	}
    items,err:=getAllItems(file)

	if err != nil {
		return err
	}	

	for _, item := range items {
		if id == item.Id {

			bt, err := json.Marshal(item)
			if err != nil {
				
				return err
			}	
			val:=strings.ToLower(string(bt))
			writer.Write([]byte(val))
			file.Close()
			
			return nil
		}
	}
	
	file.Close()	

	return nil
}

func remove(fileName string,id string, writer io.Writer) error {

	if id == "" {
			return fmt.Errorf("-id flag has to be specified")
	}
	var filePermission fs.FileMode = 0644

	file, err := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE, filePermission)
	if err != nil {
		return err
	}
    items,err:=getAllItems(file)

	if err != nil {
		return err
	}	
	file.Close()

	for i := 0; i < (len(items) - 1); {
		if id == items[i].Id {
            items=removeItemByIndex(items, i)
			bt, err := json.Marshal(items)
			if err != nil {				
				return err
		}

		val:=strings.ToLower(string(bt))
		writer.Write([]byte(val))
		file, err := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE|os.O_TRUNC, filePermission)
		if err != nil {
			return err
		}
		file.Write([]byte(val))
		file.Close()			
		return nil
		} 
	}
	
	file.Close()
	writer.Write([]byte(fmt.Sprintf("Item with id %s not found",id)))

	return nil
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
	items,err:=getAllItems(file)
	if err != nil {
		return err
	}	

	for _, i := range items {
		if fmt.Sprint(item.Id) == i.Id {

			writer.Write([]byte(fmt.Sprintf("Item with id %s already exists",i.Id)))
			file.Close()
			return nil
		}
	}	

	items = append(items, item)

	bt, err := json.Marshal(items)
	if err != nil {		
		return err
	}

	val:=strings.ToLower(string(bt))
	file.Write([]byte(val))
	file.Close()

	writer.Write([]byte(val))

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

func getAllItems(file *os.File) ([]Item,error){
	bytes, err := ioutil.ReadAll(file)
	if err != nil {
		return nil,err
	}
	var items []Item
	if len(bytes)!=0 {
		err = json.Unmarshal(bytes, &items)
	if err != nil {
		return nil,err
	}
		return items,nil
    }
   return nil,nil
}

func removeItemByIndex(s []Item, index int) []Item {
	return append(s[:index], s[index+1:]...)
}

func main() {
	err := Perform(parseArgs(), os.Stdout)
	if err != nil {
		panic(err)
	}
}
