package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
)

var (
	protoPath = flag.String("proto-path",
		"/usr/share/xcb", "path to directory of X protocol XML files")
	packageName = flag.String("package", "xkb", "name of package to generate")
	gofmt       = flag.Bool("gofmt", true,
		"When disabled, gofmt will not be run before outputting Go code")
)

func usage() {
	basename := os.Args[0]
	if lastSlash := strings.LastIndex(basename, "/"); lastSlash > -1 {
		basename = basename[lastSlash+1:]
	}
	log.Printf("Usage: %s [flags] xml-file", basename)
	flag.PrintDefaults()
	os.Exit(1)
}

func init() {
	log.SetFlags(0)
}

func main() {
	flag.Usage = usage
	flag.Parse()

	// Read the single XML file into []byte
	xmlBytes, err := os.ReadFile(flag.Arg(0))
	if err != nil {
		log.Fatal(err)
	}

	// Initialize the buffer, parse it, and filter it through gofmt.
	c := newContext()
	c.Morph(xmlBytes)

	outFile, err := os.Create(fmt.Sprintf("%s/%s.go", *packageName, *packageName))
	if err != nil {
		panic(err)
	}
	outFile.Truncate(0)
	defer outFile.Close()
	out := bufio.NewWriter(outFile)

	if !*gofmt {
		c.out.WriteTo(out)
	} else {
		cmdGofmt := exec.Command("gofmt")
		cmdGofmt.Stdin = c.out
		cmdGofmt.Stdout = out
		cmdGofmt.Stderr = os.Stderr
		err = cmdGofmt.Run()
		if err != nil {
			log.Fatal(err)
		}
	}

	out.Flush()
}
