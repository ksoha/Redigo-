package command

import (
	"bufio"
	"bytes"
	"testing"

	"github.com/ksoha/redigo/internal/store"
)

// helper function:
// creates a Store and an in-memory writer for each test.
func setup() (*store.Store, *bufio.Writer, *bytes.Buffer) {
	s := store.New()

	var buf bytes.Buffer
	writer := bufio.NewWriter(&buf)

	return s, writer, &buf
}

// PING -> PONG
func TestDispatchPING(t *testing.T) {
	s, writer, buf := setup()

	err := Dispatch([]string{"PING"}, s, writer)
	if err != nil {
		t.Fatalf("Dispatch returned error: %v", err)
	}

	expected := "+PONG\r\n"

	if buf.String() != expected {
		t.Errorf("expected %q, got %q", expected, buf.String())
	}
}

// SET foo bar -> OK
func TestDispatchSET(t *testing.T) {
	s, writer, buf := setup()

	err := Dispatch([]string{"SET", "foo", "bar"}, s, writer)
	if err != nil {
		t.Fatalf("Dispatch returned error: %v", err)
	}

	// Check the response.
	expected := "+OK\r\n"

	if buf.String() != expected {
		t.Errorf("expected %q, got %q", expected, buf.String())
	}

	// Also check that SET actually stored the value.
	value, ok := s.Get("foo")

	if !ok {
		t.Fatal("expected key 'foo' to exist")
	}

	if value != "bar" {
		t.Errorf("expected value 'bar', got %q", value)
	}
}

// GET foo -> bar
func TestDispatchGET(t *testing.T) {
	s, writer, buf := setup()

	// Put something in the store first.
	s.Set("foo", "bar")

	err := Dispatch([]string{"GET", "foo"}, s, writer)
	if err != nil {
		t.Fatalf("Dispatch returned error: %v", err)
	}

	expected := "$3\r\nbar\r\n"

	if buf.String() != expected {
		t.Errorf("expected %q, got %q", expected, buf.String())
	}
}

// GET missing key -> Null Bulk String
func TestDispatchGETMissingKey(t *testing.T) {
	s, writer, buf := setup()

	err := Dispatch([]string{"GET", "missing"}, s, writer)
	if err != nil {
		t.Fatalf("Dispatch returned error: %v", err)
	}

	expected := "$-1\r\n"

	if buf.String() != expected {
		t.Errorf("expected %q, got %q", expected, buf.String())
	}
}

// DEL foo -> 1
func TestDispatchDEL(t *testing.T) {
	s, writer, buf := setup()

	s.Set("foo", "bar")

	err := Dispatch([]string{"DEL", "foo"}, s, writer)
	if err != nil {
		t.Fatalf("Dispatch returned error: %v", err)
	}

	expected := ":1\r\n"

	if buf.String() != expected {
		t.Errorf("expected %q, got %q", expected, buf.String())
	}

	// Make sure the key was actually deleted.
	_, ok := s.Get("foo")

	if ok {
		t.Error("expected key 'foo' to be deleted")
	}
}

// DEL missing -> 0
func TestDispatchDELMissingKey(t *testing.T) {
	s, writer, buf := setup()

	err := Dispatch([]string{"DEL", "missing"}, s, writer)
	if err != nil {
		t.Fatalf("Dispatch returned error: %v", err)
	}

	expected := ":0\r\n"

	if buf.String() != expected {
		t.Errorf("expected %q, got %q", expected, buf.String())
	}
}

// EXISTS foo -> 1
func TestDispatchEXISTS(t *testing.T) {
	s, writer, buf := setup()

	s.Set("foo", "bar")

	err := Dispatch([]string{"EXISTS", "foo"}, s, writer)
	if err != nil {
		t.Fatalf("Dispatch returned error: %v", err)
	}

	expected := ":1\r\n"

	if buf.String() != expected {
		t.Errorf("expected %q, got %q", expected, buf.String())
	}
}

// EXISTS missing -> 0
func TestDispatchEXISTSMissingKey(t *testing.T) {
	s, writer, buf := setup()

	err := Dispatch([]string{"EXISTS", "missing"}, s, writer)
	if err != nil {
		t.Fatalf("Dispatch returned error: %v", err)
	}

	expected := ":0\r\n"

	if buf.String() != expected {
		t.Errorf("expected %q, got %q", expected, buf.String())
	}
}

// Unknown command -> error
func TestDispatchUnknownCommand(t *testing.T) {
	s, writer, buf := setup()

	err := Dispatch([]string{"UNKNOWN"}, s, writer)
	if err != nil {
		t.Fatalf("Dispatch returned error: %v", err)
	}

	expected := "-ERR unknown command 'UNKNOWN'\r\n"

	if buf.String() != expected {
		t.Errorf("expected %q, got %q", expected, buf.String())
	}
}

// SET with wrong number of arguments -> error
func TestDispatchSETWrongArguments(t *testing.T) {
	s, writer, buf := setup()

	err := Dispatch([]string{"SET", "foo"}, s, writer)
	if err != nil {
		t.Fatalf("Dispatch returned error: %v", err)
	}

	expected := "-ERR wrong number of arguments for 'set' command\r\n"

	if buf.String() != expected {
		t.Errorf("expected %q, got %q", expected, buf.String())
	}
}

// GET with wrong number of arguments -> error
func TestDispatchGETWrongArguments(t *testing.T) {
	s, writer, buf := setup()

	err := Dispatch([]string{"GET"}, s, writer)
	if err != nil {
		t.Fatalf("Dispatch returned error: %v", err)
	}

	expected := "-ERR wrong number of arguments for 'get' command\r\n"

	if buf.String() != expected {
		t.Errorf("expected %q, got %q", expected, buf.String())
	}
}

// Commands should be case-insensitive.
func TestDispatchCaseInsensitive(t *testing.T) {
	s, writer, buf := setup()

	err := Dispatch([]string{"set", "foo", "bar"}, s, writer)
	if err != nil {
		t.Fatalf("Dispatch returned error: %v", err)
	}

	expected := "+OK\r\n"

	if buf.String() != expected {
		t.Errorf("expected %q, got %q", expected, buf.String())
	}
}

// Empty command -> error
func TestDispatchEmptyCommand(t *testing.T) {
	s, writer, buf := setup()

	err := Dispatch([]string{}, s, writer)
	if err != nil {
		t.Fatalf("Dispatch returned error: %v", err)
	}

	expected := "-ERR empty command\r\n"

	if buf.String() != expected {
		t.Errorf("expected %q, got %q", expected, buf.String())
	}
}
