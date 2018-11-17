package main

import (
  "bytes"
  "bufio"
  "github.com/gotk3/gotk3/gtk"
  "github.com/gotk3/gotk3/glib"
  "log"
  "net"
)

func ScanNulls(data []byte, atEOF bool) (advance int, token []byte, err error) {
  if atEOF && len(data) == 0 {
    return 0, nil, nil
  }
  if i := bytes.IndexByte(data, '\000'); i >= 0 {
    return i + 1, data[0:i], nil
  }
  if atEOF {
    return len(data), data, nil
  }
  return 0, nil, nil
}

func HandleConnection(output chan<- string, conn net.Conn) {
  defer conn.Close()
  scanner := bufio.NewScanner(conn)
  scanner.Split(ScanNulls)
  for scanner.Scan() {
    output <- scanner.Text()
  }
}

func Serve(output chan<- string) {
  ln, err := net.Listen("tcp", ":5986")
  if err != nil {
    panic(err)
  }
  for {
    conn, err := ln.Accept()
    if err != nil {
      panic(err)
    }
    go HandleConnection(output, conn)
  }
}

func HandleInput(input <-chan string) {
  for {
    newText, more := <-input
    if !more {
      return
    }
    log.Print(newText)
    _, err := glib.IdleAdd(ChangeText, newText)
    if err != nil {
      panic(err)
    }
  }
}

func ChangeText(s string) {
  text.SetText(s)
}

var win *gtk.Window
var text *gtk.Label

func main() {
  var err error
  gtk.Init(nil)
  win, err = gtk.WindowNew(gtk.WINDOW_TOPLEVEL)
  if err != nil {
    panic(err)
  }
  win.SetTitle("Subaru")
  win.Connect("destroy", func() {
    gtk.MainQuit()
  })
  text, err = gtk.LabelNew("subaru subaru subaru kyun")
  if err != nil {
    panic(err)
  }
  text.SetFont("Lato 42")
  win.Add(text)
  prescott := make(chan string, 5)
  go Serve(prescott)
  go HandleInput(prescott)
  win.ShowAll()
  gtk.Main()
}
