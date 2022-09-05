package main

import (
  "bufio"
  "flag"
  "fmt"
  "github.com/tarm/serial"
  "log"
  "os"
)

// Command Line Args
// EX: ./serial -com=COM3 -baud=115200
var (
  com  string
  baud int;
)
//var  s *serial.Port

// Set Defaults
func init() {
  flag.StringVar(&com, "com", "COM3", "The COM port to read data from")
  flag.IntVar(&baud, "baud", 115200, "The Baud Rate")
}

func openCom() (*serial.Port) {

  flag.Parse()

  // Program Info
  log.Printf("Starting with port %s at baud rate %d", com, baud)

  // Setup Serial Port
  c := &serial.Config{Name: com, Baud: baud}
  s, err := serial.OpenPort(c)
  if err != nil {
    log.Fatal(err)
  }

  log.Printf("Opened %s", com)
  return s
}


func readcom(s *serial.Port) {

  log.Print("Reading Bytes in...")

  buf := make([]byte, 2048)

  r, err := s.Read(buf)
  if err != nil {
    log.Fatal(err)
  }

  // Print Bytes Read
  log.Printf("% #x ", buf[:r])
  file, err := os.Create("test.txt")
  if err != nil {
    log.Fatal(err)
  }

  defer file.Close()

  w := bufio.NewWriter(file)
  for _, abyte := range buf {
    fmt.Fprintf(w, "% #x ", abyte)
  }
  write_err := w.Flush()

  if write_err != nil {
    log.Fatal(write_err)
  }
}
