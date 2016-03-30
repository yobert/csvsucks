package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strings"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Printf("\nusage:\n\t%s (file names)\n\n", os.Args[0])
		os.Exit(1)
	}

	for i, fpath := range os.Args {
		if i == 0 {
			continue
		}

		err := do(fpath)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}
}

func do(fpath string) error {
	fh, err := os.Open(fpath)
	if err != nil {
		return err
	}
	defer fh.Close()

	out_fpath := fpath + ".tsv"
	out_fh, err := os.Create(out_fpath)
	if err != nil {
		return err
	}
	defer out_fh.Close()

	csvr := csv.NewReader(fh)
	csvr.FieldsPerRecord = -1
	linecount := 0
	for {
		row, err := csvr.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		for i, v := range row {
			if strings.IndexByte(v, '\t') != -1 {
				return fmt.Errorf("file %#v line %d has a tab in cell %d, so we can't encode that as TSV")
			}
			if i != 0 {
				fmt.Fprint(out_fh, "\t")
			}
			fmt.Fprint(out_fh, v)
		}
		fmt.Fprintln(out_fh)
		linecount++
	}

	fmt.Printf("Converted %#v to %#v (%d lines)\n", fpath, out_fpath, linecount)
	return nil
}
