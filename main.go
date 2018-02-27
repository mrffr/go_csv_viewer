package main

import (
  "os"
  "fmt"
)

var fields_n int;
var records [][]string;

func main(){
  filePath := ""
  if len(os.Args) > 1 {
    filePath = os.Args[1]
  } else {
    //TODO print usage
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
