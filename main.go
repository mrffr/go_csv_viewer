package main

import (
  "os"
  "fmt"
  "flag"
)

//TODO
//sorting -- checking types
//highlight full line across columns

type csvView struct {
  fields_n int
  has_header bool
  records [][]string //actual data
  width_ratios []float64 //testing column sizing calculations
}

var mv csvView

func main(){
  filePath := ""

  var has_header = flag.Bool("h", false, "the file has a header")
  flag.Parse()

  if len(os.Args) > 1 {
    filePath = os.Args[len(os.Args)-1] //assume file comes last
  } else {
    fmt.Println("Usage:",os.Args[0],"FLAGS <csv file>")
    flag.PrintDefaults()
    os.Exit(-1)
  }

  mv.has_header = *has_header
  read_file(&mv, filePath)


  /*
  if mv.has_header {
    for i:=0;i<mv.fields_n;i++ { fmt.Println(mv.width_ratios[i])}
  }
  */

  run_ui()
}
