package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
  "errors"
)

// Load a csv file.
func open_csv_file(filePath string) *os.File {
	f, err := os.Open(filePath)
	if err != nil {
		fmt.Println(filePath, "not found.")
		os.Exit(-1)
	}
  return f
}

func read_file(v *csvView, f *os.File) error {
	r := csv.NewReader(f)


  //peak first line
  header, err := r.Read()
  if err == io.EOF {
    return errors.New("Empty file")
  }

	if mv.has_header {
		mv.header = header
	}

	fields_n := len(header)
	widths := make([]int, fields_n)
	width_ratios := make([]float64, fields_n)
  max_widths := make([]int, fields_n)

  //if we don't have a header unget the first line
  if !mv.has_header {
    f.Seek(0,0)
    r = csv.NewReader(f)
  }

	cnt := 0
	var records [][]string
rloop:
	for ; ; cnt++ {
		rec, err := r.Read()
		if err == io.EOF {
      if cnt == 0 {
        return errors.New("No records found")
      }
			break rloop
		}

		for i := 0; i < fields_n; i++ {
			lr := len(rec[i]) + 1 //??? \0
			widths[i] += lr
      if(lr >= max_widths[i]){
        max_widths[i] = lr
      }
		}
		records = append(records, rec)
	}

	//ratios
	for i := 0; i < fields_n; i++ {
		widths[i] /= cnt
	}
	tot := 0
	for i := 0; i < fields_n; i++ {
		tot += widths[i]
	}

	for i := 0; i < fields_n; i++ {
		width_ratios[i] = float64(widths[i]) / float64(tot)
	}

	//
	v.fields_n = fields_n
	v.records_len = cnt
	v.records = records
	v.width_ratios = width_ratios
  v.max_widths = max_widths

  return nil
}
