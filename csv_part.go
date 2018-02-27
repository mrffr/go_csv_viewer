package main

import (
  "fmt"
  "os"
  "bufio"
  "encoding/csv"
)

func read_file(filePath string) [][]string{
  // Load a csv file.
  f, err := os.Open(filePath)
  if err != nil {
    fmt.Println(filePath, "not found.")
    os.Exit(-1)
  }

  r := csv.NewReader(bufio.NewReader(f))

  //TODO read indivdual to calc field width for layout
  dat, err := r.ReadAll()
  if err != nil { panic(err) }

  f.Close()
  return dat
}
