package main

import (
  "fmt"
  "os"
  "io"
  "encoding/csv"
)

func read_file(v *csvView, filePath string){
  // Load a csv file.
  f, err := os.Open(filePath)
  if err != nil {
    fmt.Println(filePath, "not found.")
    os.Exit(-1)
  }

  r := csv.NewReader(f)

  header, err := r.Read()
  if err == io.EOF {
    fmt.Println("No records found in file")
    os.Exit(-1)
  }

  var records [][]string
  if mv.has_header {
    mv.header = header
  }else{
    records = append(records, header)
  }

  fields_n := len(header)
  widths := make([]int, fields_n)
  width_ratios := make([]float64, fields_n)


  cnt := 0
  rloop:
  for ;;cnt++{
    rec, err := r.Read()
    if err == io.EOF {
      cnt -= 1
      break rloop
    }

    for i:= 0; i<fields_n;i++ {
      lr := len(rec[i])
      widths[i] += lr
    }
    records = append(records, rec)
  }

  //ratios
  for i :=0;i<fields_n;i++ { widths[i] /= cnt }
  tot := 0
  for i :=0;i<fields_n;i++ { tot += widths[i] }

  for i :=0;i<fields_n;i++ { width_ratios[i] = float64(widths[i]) / float64(tot) }

  //
  v.fields_n = fields_n
  v.records_len = cnt
  v.records = records
  v.width_ratios = width_ratios

  f.Close()
}

