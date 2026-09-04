package command

import (
	"bufio"
	"strings"

	"github.com/ksoha/redigo/internal/resp"

	"github.com/ksoha/redigo/internal/store"
)

// Dispatch takes a parsed command (e.g. ["SET", "foo", "bar"]),
// decides what it means, runs it against the store, and writes the
// correctly-typed RESP response back to the client.

func Dispatch(args []string, s *store.Store, w *bufio.Writer) error {
	if len(args) == 0 {
		return resp.WriteError(w, "ERR empty command")
	}

	// Command names are case-insensitive in real Redis (SET, set, Set
	// all work) — so we normalize to uppercase before matching.
	cmd := strings.ToUpper(args[0])

	switch cmd {
	case "PING":
		return resp.WriteSimpleString(w, "PONG")

	case "SET":
		if len(args) != 3 {
			return resp.WriteError(w, "ERR wrong number of arguments for 'set' command")
		}
		key, value := args[1], args[2]
		s.Set(key, value)
		return resp.WriteSimpleString(w, "OK")

	case "GET":
		if len(args) != 2 {
			return resp.WriteError(w, "ERR wrong number of arguments for 'get' command")
		}
		key := args[1]
		value, ok := s.Get(key)
		if !ok {
			return resp.WriteNullBulkString(w)
		}
		return resp.WriteBulkString(w, value)

	case "DEL":
		if len(args) != 2 {
			return resp.WriteError(w, "ERR wrong number of arguments for 'del' command")
		}
		key := args[1]
		count := s.Delete(key)
		return resp.WriteInteger(w, count)

	case "EXISTS":
		if len(args) != 2 {
			return resp.WriteError(w, "ERR wrong number of arguments for 'exists' command")
		}
		key := args[1]
		if s.Exists(key) {
			return resp.WriteInteger(w, 1)
		}
		return resp.WriteInteger(w, 0)

	case "LPUSH":
		if len(args) != 3 {
			return resp.WriteError(w, "ERR wrong number of arguments for 'lpush' command")
		}
		key, value := args[1], args[2]
		count, err := s.LPush(key, value)
		if err != nil {
			return resp.WriteError(w, err.Error())
		}
		return resp.WriteInteger(w, count)

	case "RPUSH":
		if len(args) != 3 {
			return resp.WriteError(w, "ERR wrong number of arguments for 'rpush' command")
		}
		key, value := args[1], args[2]
		count, err := s.RPush(key, value)
		if err != nil {
			return resp.WriteError(w, err.Error())
		}
		return resp.WriteInteger(w, count)

	default:
		return resp.WriteError(w, "ERR unknown command '"+cmd+"'")
	}
}
