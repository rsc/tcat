// Copyright 2018 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Tcat is a tabular cat.
package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"unicode/utf8"
)

var rows [][]string

func usage() {
	fmt.Fprintf(os.Stderr, "usage: tcat [file...]\n")
	os.Exit(2)
}

func main() {
	log.SetPrefix("tcat: ")
	log.SetFlags(0)
	flag.Usage = usage
	flag.Parse()
	if flag.NArg() == 0 {
		read(os.Stdin)
	} else {
		for _, arg := range flag.Args() {
			f, err := os.Open(arg)
			if err != nil {
				log.Print(err)
				continue
			}
			read(f)
			f.Close()
		}
	}
	printTable(os.Stdout, rows)
}

func read(r io.Reader) {
	data, err := ioutil.ReadAll(r)
	if err != nil {
		log.Print(err)
	}
	for _, line := range strings.SplitAfter(string(data), "\n") {
		if line == "" {
			continue
		}
		row := strings.Fields(line)
		rows = append(rows, row)
	}
}

func printTable(w io.Writer, rows [][]string) {
	var max []int
	for _, row := range rows {
		for i, c := range row {
			n := utf8.RuneCountInString(c)
			if i >= len(max) {
				max = append(max, n)
			} else if max[i] < n {
				max[i] = n
			}
		}
	}

	b := bufio.NewWriter(w)
	for _, row := range rows {
		for len(row) > 0 && row[len(row)-1] == "" {
			row = row[:len(row)-1]
		}
		for i, c := range row {
			b.WriteString(c)
			if i+1 < len(row) {
				for j := utf8.RuneCountInString(c); j < max[i]+2; j++ {
					b.WriteRune(' ')
				}
			}
		}
		b.WriteRune('\n')
	}
	b.Flush()
}
