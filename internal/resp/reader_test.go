package resp

import (
	"bufio"
	"strings"
	"testing"
)

func TestReadCommandSET(t *testing.T) {
	input := "*3\r\n$3\r\nSET\r\n$3\r\nfoo\r\n$3\r\nbar\r\n"

	reader := bufio.NewReader(strings.NewReader(input))

	got, err := ReadCommand(reader)
	if err != nil {
		t.Fatalf("ReadCommand returned error: %v", err)
	}

	expected := []string{"SET", "foo", "bar"}

	if len(got) != len(expected) {
		t.Fatalf("expected %d arguments, got %d", len(expected), len(got))
	}

	for i := range expected {
		if got[i] != expected[i] {
			t.Errorf("expected %q at index %d, got %q", expected[i], i, got[i])
		}
	}
}

func TestReadCommandGET(t *testing.T) {
	input := "*2\r\n$3\r\nGET\r\n$3\r\nfoo\r\n"

	reader := bufio.NewReader(strings.NewReader(input))

	got, err := ReadCommand(reader)
	if err != nil {
		t.Fatalf("ReadCommand returned error: %v", err)
	}

	expected := []string{"GET", "foo"}

	if len(got) != len(expected) {
		t.Fatalf("expected %d arguments, got %d", len(expected), len(got))
	}

	for i := range expected {
		if got[i] != expected[i] {
			t.Errorf("expected %q at index %d, got %q", expected[i], i, got[i])
		}
	}
}

func TestReadCommandInvalidArray(t *testing.T) {
	input := "+OK\r\n"

	reader := bufio.NewReader(strings.NewReader(input))

	_, err := ReadCommand(reader)

	if err == nil {
		t.Error("expected error for invalid RESP array")
	}
}

func TestReadCommandInvalidBulkString(t *testing.T) {
	input := "*1\r\n+OK\r\n"

	reader := bufio.NewReader(strings.NewReader(input))

	_, err := ReadCommand(reader)

	if err == nil {
		t.Error("expected error for invalid bulk string")
	}
}

func TestReadCommandEmptyArray(t *testing.T) {
	input := "*0\r\n"

	reader := bufio.NewReader(strings.NewReader(input))

	got, err := ReadCommand(reader)
	if err != nil {
		t.Fatalf("ReadCommand returned error: %v", err)
	}

	if len(got) != 0 {
		t.Errorf("expected empty command, got %v", got)
	}
}
