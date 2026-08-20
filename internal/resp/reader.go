package resp

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
)

//ReadCommand reads one full RESP command from the connection
//and returns a slice of string ["set" , "foo", "bar"] for example

//every command a client sends is always a RESP array of bulk strings.
//eg :
//	*3\r\n
//	$3\r\nSET\r\n
//	$3\r\nfoo\r\n
//	$3\r\nbar\r\n

func ReadCommand(r *bufio.Reader) ([]string, error) {
	//step 1: read the first line , it should start with *
	//and tell how many elements are in the array
	line, err := readLine(r)
	if err != nil {
		return nil, err
	}

	if len(line) == 0 || line[0] != '*' {
		return nil, fmt.Errorf("expected array, got: %q", line)
	}

	//stripe the * and parse the number after it
	numElements, err := strconv.Atoi(line[1:]) //atoi converts string to int
	if err != nil {
		return nil, fmt.Errorf("invalid array length: %q", line)
	}

	//step 2: read the exact bulk string
	args := make([]string, 0, numElements) //crearting an empty string slice
	for i := 0; i < numElements; i++ {
		arg, err := readBulkString(r)
		if err != nil {
			return nil, err
		}
		args = append(args, arg)
	}
	return args, nil
}

// readBulkString reads a single RESP Bulk String, e.g.:
//
//	$3\r\n
//	SET\r\n

func readBulkString(r *bufio.Reader) (string, error) {
	line, err := readLine(r)
	if err != nil {
		return "", err
	}

	if len(line) == 0 || line[0] != '$' {
		return "", fmt.Errorf("expected bulk string, got: %q", line)
	}

	// line looks like "$3" — the length of the string that follows.
	length, err := strconv.Atoi(line[1:])
	if err != nil {
		return "", fmt.Errorf("invalid bulk string length: %q", line)
	}

	//read exactly the length bytes thats the actual string content
	buf := make([]byte, length)
	if _, err := io.ReadFull(r, buf); err != nil {
		return "", err
	}

	//read the trailing \r\n after the bulk string
	if _, err := readLine(r); err != nil {
		return "", err
	}

	return string(buf), nil
}

// readLine reads bytes up until \r\n, and returns the line WITHOUT
// the trailing \r\n. bufio.Reader.ReadString does the heavy lifting;
// we just trim off the 2 control characters at the end.

func readLine(r *bufio.Reader) (string, error) {
	line, err := r.ReadString('\n')
	if err != nil {
		return "", err
	}
	//line ends with \r\n, so we trim the last 2 bytes
	return line[:len(line)-2], nil
}
