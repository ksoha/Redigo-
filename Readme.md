# Redigo — Development Roadmap

Redigo is a Redis-compatible in-memory key-value store built from scratch in Go.

The project is divided into five development phases. Each phase introduces one major systems concept and ends with a clear implementation milestone.

---

## Phase 1 — RESP Protocol, TCP Server & Core Commands

**Timeline:** Days 1–3
**Status:** Completed

### Objective

Build the basic Redis-compatible server infrastructure:

* Implement the RESP protocol.
* Create a TCP server.
* Support multiple concurrent client connections.
* Parse incoming commands.
* Encode outgoing responses.
* Implement basic string commands.

### Features Implemented

* RESP request parsing.
* RESP response encoding.
* TCP listener on port `6380`.
* One goroutine per client connection.
* Command dispatcher.
* Thread-safe in-memory store using `sync.RWMutex`.

### Supported Commands

| Command         | Description                        |
| --------------- | ---------------------------------- |
| `PING`          | Checks whether the server is alive |
| `SET key value` | Stores a string value              |
| `GET key`       | Retrieves a string value           |
| `DEL key`       | Deletes a key                      |
| `EXISTS key`    | Checks whether a key exists        |

### Example

```bash
redis-cli -p 6380 PING
```

```text
PONG
```

```bash
redis-cli -p 6380 SET name Redigo
```

```text
OK
```

```bash
redis-cli -p 6380 GET name
```

```text
"Redigo"
```

### Technical Concepts

* TCP networking in Go.
* `net.Listener` and `net.Conn`.
* RESP protocol design.
* `bufio.Reader` and `bufio.Writer`.
* Goroutines.
* Mutexes and concurrent map access.
* Command routing.

### Completion Milestone

The following command works successfully:

```bash
redis-cli -p 6380 SET foo bar
```

The server can accept multiple clients and correctly process core string commands.

---

## Phase 2 — Richer Data Types

**Timeline:** Days 4–6
**Status:** In Progress

### Objective

Extend the in-memory store beyond simple strings by implementing Redis-style collections.

This phase focuses on:

* Internal data representation.
* Type checking.
* Collection operations.
* Correct error handling.
* Testing commands against `redis-cli`.

---

### 2.1 Lists

Lists are internally represented using Go slices.

#### Supported Commands

| Command                 | Description                             |
| ----------------------- | --------------------------------------- |
| `LPUSH key value`       | Adds a value to the beginning of a list |
| `RPUSH key value`       | Adds a value to the end of a list       |
| `LRANGE key start stop` | Returns a range of list elements        |
| `LPOP key`              | Removes and returns the first element   |
| `RPOP key`              | Removes and returns the last element    |

#### Example

```bash
redis-cli -p 6380 LPUSH numbers 10
redis-cli -p 6380 RPUSH numbers 20
redis-cli -p 6380 LRANGE numbers 0 -1
```

Expected result:

```text
10
20
```

---

### 2.2 Hashes

Hashes are internally represented using:

```go
map[string]string
```

Each Redis key stores a collection of field-value pairs.

#### Supported Commands

| Command                | Description                         |
| ---------------------- | ----------------------------------- |
| `HSET key field value` | Creates or updates a hash field     |
| `HGET key field`       | Retrieves the value of a hash field |
| `HGETALL key`          | Returns all fields and values       |

#### Example

```bash
redis-cli -p 6380 HSET user:1 name Soha
redis-cli -p 6380 HSET user:1 role BackendEngineer
redis-cli -p 6380 HGET user:1 name
redis-cli -p 6380 HGETALL user:1
```

---

### 2.3 Sets

Sets are internally represented using a Go map where the values act as unique members.

A possible representation is:

```go
map[string]struct{}
```

#### Supported Commands

| Command                | Description                    |
| ---------------------- | ------------------------------ |
| `SADD key member`      | Adds a member to a set         |
| `SMEMBERS key`         | Returns all set members        |
| `SISMEMBER key member` | Checks whether a member exists |

#### Example

```bash
redis-cli -p 6380 SADD skills Go
redis-cli -p 6380 SADD skills Redis
redis-cli -p 6380 SISMEMBER skills Go
redis-cli -p 6380 SMEMBERS skills
```

---

### Type Safety

Redigo should reject invalid operations on incompatible data types.

For example:

```bash
SET user:1 Soha
HGET user:1 name
```

The second command should return a WRONGTYPE-style error because `user:1` contains a string instead of a hash.

### Testing Requirements

Each data type should include:

* Unit tests for store operations.
* Dispatcher tests.
* RESP response tests.
* Missing-key tests.
* Overwrite/update tests.
* Wrong-type tests.
* Concurrent access tests.
* Race detector validation.

Run:

```bash
go test ./...
```

Run with the race detector:

```bash
go test -race ./...
```

### Completion Milestone

Phase 2 is complete when:

* Lists work correctly.
* Hashes work correctly.
* Sets work correctly.
* Wrong-type operations return errors.
* Unit and integration tests pass.
* Commands work through `redis-cli`.

---

## Phase 3 — Key Expiry & TTL

**Timeline:** Days 7–8
**Status:** Planned

### Objective

Implement automatic key expiration.

Keys should be able to expire after a specified duration and should no longer be accessible after expiration.

### Planned Commands

| Command                    | Description                                |
| -------------------------- | ------------------------------------------ |
| `EXPIRE key seconds`       | Sets expiry in seconds                     |
| `PEXPIRE key milliseconds` | Sets expiry in milliseconds                |
| `TTL key`                  | Returns remaining lifetime in seconds      |
| `PTTL key`                 | Returns remaining lifetime in milliseconds |
| `PERSIST key`              | Removes expiry from a key                  |

### Expiry Representation

Each key may store expiration metadata alongside its value.

Conceptually:

```go
type Entry struct {
    Value     Value
    ExpiresAt time.Time
}
```

Keys without expiry will have no expiration timestamp.

### Passive Expiry

Passive expiry happens when a key is accessed.

Example:

```text
GET key
    ↓
Check whether key is expired
    ↓
If expired, delete key
    ↓
Return nil / key-not-found response
```

This prevents expired keys from being returned to clients.

### Active Expiry

Active expiry uses a background goroutine that periodically checks expired keys.

The first implementation may use a simple periodic scan. If time permits, it can be improved using sampling-based expiration to reduce unnecessary work.

### Important Design Considerations

* Expired keys must not be returned.
* `TTL` should return the remaining lifetime.
* Updating a key should correctly handle its expiry.
* `PERSIST` should remove the expiration.
* Expiry operations must be thread-safe.
* The background goroutine should shut down cleanly.

### Testing Requirements

Test:

* Key expires after the expected duration.
* Expired keys cannot be retrieved.
* `TTL` returns a decreasing value.
* `PEXPIRE` works with millisecond precision.
* `PERSIST` removes expiry.
* Updating a key resets or preserves expiry according to the chosen design.
* Background expiry removes expired keys.
* Concurrent reads and expiry operations are safe.

### Completion Milestone

This command should work:

```bash
redis-cli -p 6380 SET temporary value
redis-cli -p 6380 EXPIRE temporary 5
```

After approximately five seconds:

```bash
redis-cli -p 6380 GET temporary
```

The key should no longer exist.

---

## Phase 4 — AOF Persistence

**Timeline:** Days 9–11
**Status:** Planned

### Objective

Add durability using an Append-Only File.

Every mutating command will be written to a log file. When Redigo restarts, the log will be replayed to rebuild the in-memory state.

### Why AOF?

AOF is selected instead of snapshot-based persistence because:

* It is simpler to implement for this project.
* It records operations in sequence.
* It connects with write-ahead logging concepts.
* It makes the recovery process easy to understand.
* It provides a strong systems-design discussion topic.

### Write Flow

```text
Client sends command
        ↓
RESP parser
        ↓
Command dispatcher
        ↓
Update in-memory store
        ↓
Append command to AOF
        ↓
Return response to client
```

### Recovery Flow

```text
Server starts
        ↓
Open AOF file
        ↓
Read commands sequentially
        ↓
Replay commands
        ↓
Rebuild in-memory store
        ↓
Start accepting clients
```

### Example

```bash
redis-cli -p 6380 SET name Soha
redis-cli -p 6380 HSET user:1 role backend
redis-cli -p 6380 LPUSH tasks coding
```

After restarting the server, these values should still be available.

### Planned AOF Features

* Append mutating commands.
* Create the AOF file if it does not exist.
* Replay commands during startup.
* Handle malformed or incomplete final log entries.
* Add configurable AOF file path.
* Flush writes according to the selected durability policy.
* Avoid logging read-only commands such as `GET` and `PING`.

### Fsync Policies

The implementation should document the selected write policy.

Possible policies:

| Policy       | Description                        |
| ------------ | ---------------------------------- |
| Always       | Flush every write immediately      |
| Every second | Flush periodically                 |
| OS-managed   | Rely on operating-system buffering |

The initial version may use a simple policy, while the README documents the durability and performance tradeoff.

### Important Design Considerations

* AOF replay should not append replayed commands again.
* Startup should handle an empty AOF file.
* Partial final commands should not crash the server.
* Mutating commands must be logged consistently.
* Expiry behavior during replay must be defined.
* AOF writes should be synchronized with store updates.
* Persistence errors should be handled clearly.

### Completion Milestone

The following workflow should work:

```bash
docker compose up --build
```

```bash
redis-cli -p 6380 SET name Redigo
redis-cli -p 6380 HSET user:1 role backend
```

Stop and restart the server:

```bash
docker compose down
docker compose up --build
```

The previously stored data should be restored from the AOF file.

---

## Phase 5 — Benchmarking, Documentation & Final Polish

**Timeline:** Days 12–13
**Status:** Planned

### Objective

Measure Redigo's performance, document the architecture, and prepare the project for GitHub and interviews.

### Benchmarking

Use `redis-benchmark` to measure:

* Requests per second.
* Average latency.
* Throughput under concurrent clients.
* Performance of read commands.
* Performance of write commands.
* Performance with persistence enabled.
* Performance with persistence disabled, if supported.

### Example

```bash
redis-benchmark -p 6380 -t get,set -n 10000
```

The actual benchmark results should be recorded honestly after testing.

### Suggested Benchmark Scenarios

#### Basic SET/GET

Measure the performance of:

```text
SET
GET
```

#### Concurrent Clients

Run benchmarks with different client counts to observe how the server behaves under concurrency.

#### Data Type Commands

Benchmark:

```text
LPUSH
LRANGE
HSET
HGET
SADD
SMEMBERS
```

#### Persistence Impact

Compare performance:

```text
AOF disabled
AOF enabled
```

This will help demonstrate the cost of durability.

### Documentation

The final README should include:

* Project overview.
* Motivation.
* Features.
* Architecture diagram.
* Project structure.
* Supported commands.
* RESP protocol explanation.
* Concurrency model.
* Store design.
* Expiry implementation.
* AOF design.
* Recovery process.
* Benchmark results.
* Testing instructions.
* Docker instructions.
* Design tradeoffs.
* Limitations.
* Future improvements.

### Final Testing Checklist

* [ ] `go test ./...`
* [ ] `go test -race ./...`
* [ ] Manual `redis-cli` testing
* [ ] Concurrent client testing
* [ ] Expiry testing
* [ ] AOF recovery testing
* [ ] Docker build testing
* [ ] Docker restart testing
* [ ] Benchmark execution
* [ ] README completed
* [ ] GitHub repository cleaned up
* [ ] Project structure documented
* [ ] Known limitations documented

### Completion Milestone

Redigo should have:

* A working Redis-compatible command interface.
* Multiple supported data types.
* Key expiration.
* AOF persistence.
* Automated tests.
* Race detector validation.
* Benchmark numbers.
* Docker support.
* Complete GitHub documentation.

---

## Day 14 — Buffer & Final Review

The final day is reserved for:

* Bug fixes.
* Docker improvements.
* README corrections.
* Code cleanup.
* Better error messages.
* Additional tests.
* Benchmark reruns.
* Final GitHub push.

---

## Future Work

The following features are intentionally outside the current project scope:

* Replication.
* Leader-follower architecture.
* Clustering.
* Sharding.
* Pub/Sub.
* Transactions such as `MULTI` and `EXEC`.
* LRU eviction.
* Pluggable storage backends.
* Advanced Redis data types.
* Distributed consensus.

These features may be explored in a future version after the single-node implementation is stable and well-tested.

---

## Final Project Goal

The goal of Redigo is not to replace Redis.

The goal is to demonstrate a practical understanding of:

* Network protocols.
* TCP servers.
* Concurrent systems.
* In-memory data structures.
* Command dispatching.
* Expiration strategies.
* Persistence and recovery.
* Performance measurement.
* Engineering tradeoffs.

A small, correct, tested, benchmarked, and documented systems project is more valuable than a large project with incomplete or unreliable features.
