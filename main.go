package main

import (
	"flag"
	"fmt"
	"os"
)

//Ideas
//sorting -- checking types
//highlight full line across columns possible???
//don't like column sizing idea. not sure fixable

type csvView struct {
	fields_n     int
	records_len  int
	has_header   bool
	header       []string
	records      [][]string //actual data
	max_widths   []int      //actual widths
	width_ratios []float64  //testing column sizing calculations
}

var mv csvView

func main() {
	filePath := ""

	var has_header = flag.Bool("h", false, "the file has a header")
	flag.Parse()

	if len(os.Args) > 1 {
		//TODO better way to do this may lie in flag
		filePath = os.Args[len(os.Args)-1] //assume file comes last
	} else {
		fmt.Println("Usage:", os.Args[0], "FLAGS <csv file>")
		flag.PrintDefaults()
		os.Exit(-1)
	}

	mv.has_header = *has_header
  f := open_csv_file(filePath)
  err := read_file(&mv, f)
  if err != nil {
    fmt.Println("Error", err)
    os.Exit(-1)
  }
  f.Close()

	run_ui()
}
