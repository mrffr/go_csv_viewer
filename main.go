package main

import (
  "os"
  "fmt"
)

//TODO
//basic determine if there is a header -- checking types
//sorting -- checking types

var fields_n int;
var records [][]string;

func main(){
  filePath := ""
  if len(os.Args) > 1 {
    filePath = os.Args[1]
  } else {
    fmt.Println("Usage:",os.Args[0],"<csv file>")
    os.Exit(-1)
  }

  records = read_file(filePath)

  if len(records) == 0 {
    fmt.Println("No records found in", filePath);
    os.Exit(-1)
  }

  fields_n = len(records[0]) //TODO check if this is wrong csv format

  run_ui()
}
