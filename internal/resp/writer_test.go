package resp

import (
	"bufio"
	"bytes"
	"testing"
)

func TestWriteSimpleString(t *testing.T) {

	var buf bytes.Buffer
	writer := bufio.NewWriter(&buf)

	err := WriteSimpleString(writer, "OK")
	if err != nil {
		t.Fatalf("WriteSimpleString failed: %v", err)
	}

	expected := "+OK\r\n"
	if buf.String() != expected {
		t.Errorf("WriteSimpleString output = %q; want %q", buf.String(), expected)
	}
}

func TestWriteError(t *testing.T) {
	var buf bytes.Buffer

	writer := bufio.NewWriter(&buf)

	err := WriteError(writer, "ERR unknown command")
	if err != nil {
		t.Fatalf("WriteError failed: %v", err)
	}

	expected := "-ERR unknown command\r\n"
	if buf.String() != expected {
		t.Errorf("WriteError output = %q; want %q", buf.String(), expected)
	}
}

func TestWriteInteger(t *testing.T) {
	var buf bytes.Buffer

	writer := bufio.NewWriter(&buf)

	err := WriteInteger(writer, 100)
	if err != nil {
		t.Fatalf("WriteInteger failed: %v", err)
	}

	expected := ":100\r\n"
	if buf.String() != expected {
		t.Errorf("WriteInteger output = %q; want %q", buf.String(), expected)
	}
}

func TestWriteBulkString(t *testing.T) {
	var buf bytes.Buffer

	writer := bufio.NewWriter(&buf)

	err := WriteBulkString(writer, "Hello")
	if err != nil {
		t.Fatalf("WriteBulkString failed: %v", err)
	}

	expected := "$5\r\nHello\r\n"
	if buf.String() != expected {
		t.Errorf("WriteBulkString output = %q; want %q", buf.String(), expected)
	}
}

func TestWriteNullBulkString(t *testing.T) {
	var buf bytes.Buffer

	writer := bufio.NewWriter(&buf)

	err := WriteNullBulkString(writer)
	if err != nil {
		t.Fatalf("WriteNullBulkString failed: %v", err)
	}

	expected := "$-1\r\n"
	if buf.String() != expected {
		t.Errorf("WriteNullBulkString output = %q; want %q", buf.String(), expected)
	}
}

func TestWriteArray(t *testing.T) {
	var buf bytes.Buffer

	writer := bufio.NewWriter(&buf)

	items := []string{"Hello", "World"}
	err := WriteArray(writer, items)
	if err != nil {
		t.Fatalf("WriteArray failed: %v", err)
	}

	expected := "*2\r\n$5\r\nHello\r\n$5\r\nWorld\r\n"
	if buf.String() != expected {
		t.Errorf("WriteArray output = %q; want %q", buf.String(), expected)
	}
}
