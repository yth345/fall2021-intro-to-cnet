/*
Task:
(1) prompts the user for the input and output filenames
(2) reads from the input file one line at a time,
(3) prepends the line count to each line, and
(4) writes the line into the output file.
*/

package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func main() {
	in_fname, out_fname := "", ""
	fmt.Printf("Enter input filename: \n")
	fmt.Scanf("%s", &in_fname)
	fmt.Printf("Enter output filename: \n")
	fmt.Scanf("%s", &out_fname)

	in_f, in_err := os.Open(in_fname)
	check(in_err)
	defer in_f.Close()

	out_f, out_err := os.Create(out_fname)
	check(out_err)
	defer out_f.Close()

	scanner := bufio.NewScanner(in_f)
	writer := bufio.NewWriter(out_f)

	var aline string = ""
	var line_cnt int = 0
	for scanner.Scan() {
		line_cnt += 1
		aline = scanner.Text()
		aline = strconv.Itoa(line_cnt) + ". " + aline + "\n"
		_, _ = writer.WriteString(aline)
	}
	writer.Flush()
}
