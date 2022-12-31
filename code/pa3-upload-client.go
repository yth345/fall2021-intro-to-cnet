/*
Task:
(1) connects to the server that runs on the workstation at port 12014
(2) prompts the user for the upload filename
(3) sends first the file size (just the number in a single line)
(4) sends next the file content (the entire file)
(5) receives a message back from the server
(6) prints what the server says
(7) closes the connection and terminates the program
*/

package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"os"
	"strconv"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func main() {
	conn, errc := net.Dial("tcp", "127.0.0.1:12014") // net.Dial("connection type", "IP address:port")
	check(errc)
	defer conn.Close()

	f_name := ""
	fmt.Printf("Enter the upload filename:")
	fmt.Scanf("%s", &f_name)

	// get uplaod file size without opening the file
	f_stat, errs := os.Stat(f_name)
	check(errs)
	f_size := int(f_stat.Size()) // file.Size() returns int64, change type to int

	// send file size and file content to the server
	// NOTE: remember to add \n for strings sending to the server
	writer := bufio.NewWriter(conn)

	_, errw := writer.WriteString(strconv.Itoa(f_size) + "\n") // or use WriterString(fmt.Sprintf("%d\n", f_size))
	check(errw)
	fmt.Printf("Send file size: %d\n", f_size)
	writer.Flush()

	f, errf := os.Open(f_name)
	check(errf)
	defer f.Close()

	f_reader := bufio.NewReader(f)
	for {
		text, errr := f_reader.ReadString('\n')
		if errr == io.EOF {
			break
		}
		_, errw := writer.WriteString(text)
		check(errw)
		writer.Flush()
	}
	writer.Flush()

	f_scanner := bufio.NewScanner(f)
	for f_scanner.Scan() {
		text := fmt.Sprintf("%s\r\n", f_scanner.Text())
		_, err := writer.WriteString(text)
		check(err)
		writer.Flush()
	}
	writer.Flush()

	// catch what the server says and print it out
	scanner := bufio.NewScanner(conn)
	if scanner.Scan() {
		fmt.Printf("Server replies: %s\n", scanner.Text())
	}
}
