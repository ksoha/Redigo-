package resp

import (
	"bufio"
	"fmt"
)

// WriteSimpleString writes a simple string eg. +OK\r\n
// Used for things like the reply to SET
func WriteSimpleString(w *bufio.Writer, s string) error {
	//fprintf writes formatted output to the writer(opp of sprintf)
	if _, err := fmt.Fprintf(w, "+%s\r\n", s); err != nil {
		return err
	}

	return w.Flush() //flushes the buffer to the underlying writer
}

// WriteError writes an error string eg, -ERR
// Used when command is invalid or has invalid arguments
func WriteError(w *bufio.Writer, s string) error {
	if _, err := fmt.Fprintf(w, "-%s\r\n", s); err != nil {
		return err
	}

	return w.Flush()
}

// WriteInteger writes an integer eg. :1000\r\n
// Used for things like DEL(count deleted) , INCR(new value)
func WriteInteger(w *bufio.Writer, n int) error {

	if _, err := fmt.Fprintf(w, ":%d\r\n", n); err != nil {
		return err
	}

	return w.Flush()
}

// WriteBulkString writes a bulk string eg. $3\r\nfoo\r\n
// Used for things like GET(value) , SET(value)
func WriteBulkString(w *bufio.Writer, s string) error {

	//length of the string is neccessary because the client needs to know how many bytes to read for the value
	if _, err := fmt.Fprintf(w, "$%d\r\n%s\r\n", len(s), s); err != nil {
		return err
	}

	return w.Flush()
}

// WriteNullBulkString writes a null bulk string eg. $-1\r\n
// Used for things like GET(key) when key does not exist
// is not the same as empty string which is $0\r\n\r\n (a bulk string of length 0)
func WriteNullBulkString(w *bufio.Writer) error {
	if _, err := fmt.Fprintf(w, "$-1\r\n"); err != nil {
		return err
	}

	return w.Flush()
}

//WriteArray writes an array of bulk string eg. *3\r\n$3\r\nfoo\r\n$3\r\nbar\r\n$3\r\nbaz\r\n

func WriteArray(w *bufio.Writer, items []string) error {
	if _, err := fmt.Fprintf(w, "*%d\r\n", len(items)); //write the array header with the number of items
	err != nil {
		return err
	}

	for _, item := range items {
		if _, err := fmt.Fprintf(w, "$%d\r\n%s\r\n", len(item), item); //Write each item as a bulk string
		err != nil {
			return err
		}

	}
	return w.Flush()

}
