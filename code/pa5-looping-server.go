/*
Task:
(1) listens at <your port#> until there’s an upload request
(2) reads from the socket first the file size (just the number in a single line)
(3) reads from the socket one line at a time
(4) prepend the line count to each line and store the new line into a new file: whatever.txt
(5) repeats (3) and (4) until all lines in the file is processed
(6) sends a message back that tells the client the original file and the new file size
(7) closes the connection and goes back to (1)
*/

package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func loop(ln net.Listener) {
	conn, _ := ln.Accept()
	defer conn.Close()
	fmt.Println("Welcome!")

	// specify which file to write in
	f, errf := os.Create("./whatever.txt")
	check(errf)
	fwriter := bufio.NewWriter(f)
	defer f.Close()

	reader := bufio.NewReader(conn)
	writer := bufio.NewWriter(conn)
	r_byte_cnt := 0
	w_byte_cnt := 0

	f_size_str, errr := reader.ReadString('\n')
	check(errr)
	f_size_str = strings.TrimSuffix(f_size_str, "\n")
	f_size, _ := strconv.Atoi(f_size_str)

	line_cnt := 0
	for {
		message, err := reader.ReadString('\n')
		check(err)
		line_cnt += 1
		newline := fmt.Sprintf("%d %s", line_cnt, message)
		r_byte_cnt += len(message)
		w_byte_cnt += len(newline)
		fwriter.WriteString(newline)
		fwriter.Flush()
		if r_byte_cnt >= f_size {
			break
		}
	}
	fwriter.Flush()

	// tell the client
	reply := fmt.Sprintf("%d bytes received, %d bytes generated\n", r_byte_cnt, w_byte_cnt)
	writer.WriteString(reply)
	writer.Flush()

	fmt.Printf("Upload file size: %d\n", f_size)
	fmt.Printf("Output file size: %d\n", w_byte_cnt)
}

func main() {
	fmt.Println("Launching server...")
	ln, _ := net.Listen("tcp", ":12014")
	defer ln.Close()

	for {
		loop(ln)
	}

}
