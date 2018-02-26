package main

import (
  "fmt"
  "os"
  "bufio"
  "encoding/csv"
)

func print_record(record []string){
  for value := range record {
    fmt.Printf("  %v\n", record[value])
  }
}

// We failed opening file so try to create the file instead
func fail_create(filePath string) *os.File{
  fmt.Println(filePath, "was not found creating new file", filePath)
  f, err := os.Create(filePath)
  if err != nil { panic(err) }
  return f
}

func read_file(filePath string) [][]string{
  // Load a csv file.
  f, e := os.Open(filePath)
  if e != nil {
    f = fail_create(filePath)
  }

  // Create a new reader.
  r := csv.NewReader(bufio.NewReader(f))
  r.Comma = ';'

  dat, err := r.ReadAll()
  if err != nil { panic(err) }

  f.Close()
  return dat
}

func write_file(filePath string, dat [][]string){
  f, e := os.OpenFile(filePath, os.O_RDWR, 0644) //TODO look at flags
  if e != nil {
    f = fail_create(filePath)
  }

  w := csv.NewWriter(bufio.NewWriter(f))
  w.Comma = ';'

  err := w.WriteAll(dat)
  if err != nil { panic(err) }

  if err := w.Error(); err != nil { panic(err) }
  f.Close()
}
