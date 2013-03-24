# Tempered

Read docs [at godocs.org](http://godoc.org/github.com/janne/tempered).

## Simple example w/o error handling

    package main

    import "fmt"
    import "github.com/janne/tempered"

    func main() {
      t, _ := tempered.New()
      defer t.Close()
      sensing, _ := t.Devices[0].Sense()
      fmt.Printf("%.2f°C %.1f%%RH\n", sensing.TempC, sensing.RelHum)
    }

