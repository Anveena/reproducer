## OpenSSL QUIC Double‑Free Reproducer (Issue #28501)

This repository provides a minimal QUIC server (Go, quic-go) and a multi-threaded QUIC client (C, OpenSSL 3.5.x) to reproduce a suspected double-free in OpenSSL QUIC error handling under concurrent connections.

Reference: https://github.com/openssl/openssl/issues/28501

---

## 1) Server

Run:

```bash
cd openSSL/service
GO111MODULE=on go run main.go
```

Behavior (generation barrier):

- The server listens on UDP 127.0.0.1:45678 (ALPN: `anAlpnForTest`, TLS 1.3 only).
- Each client sends two little-endian uint32 values: `generation` and `total`.
- The server groups clients by `generation`. Once the number of clients in a generation reaches `total`, it cleanly closes all those connections/streams in that generation.
- No interactive input is required.

Relevant files: `service/main.go`, `service/quic/main_loop.go`, `service/quic/manager.go`, `service/quic/tls.go`.

---

## 2) Client

Build (using the provided OpenSSL install via CMake config files):

```bash
cd openSSL/client
cmake -S . -B build
cmake --build build -j
```

Optional: build OpenSSL locally for static linking and sanitizers (defaults to 3.5.3). Adjust flags as needed inside the script:

```bash
cd openSSL/client/scripts
bash build_OpenSSL.sh
```

Run:

```bash
cd openSSL/client/build
./client
```

Behavior:

- Spawns multiple threads (default `THREAD_COUNT = 10`).
- All threads use the same `generation` and set `total = THREAD_COUNT`.
- Each thread performs a QUIC handshake to 127.0.0.1:45678, sends the 8-byte header, reads a small response, and exits.
- Expected outcome: all threads finish and exit cleanly without crashes.

Adjust concurrency by editing `THREAD_COUNT` in `client/src/main.c`.

---

## 3) Reproduction Procedure

1. Start the server (see above).
2. Build and run the client.
3. The server will wait until all clients of the same generation have connected, then cleanly close those connections.
4. Expected: clean termination on both sides. If the bug is present, AddressSanitizer may report a double-free in OpenSSL error handling during concurrent QUIC operations.

Example (abridged) ASan report observed with OpenSSL 3.5.x:

```text
ERROR: AddressSanitizer: attempting double-free ...
#0 free
#1 ERR_pop_to_mark
#2 ... demux_recv ...
... another thread ...
#1 ERR_clear_error
```

---

## 4) Environment

- OS: Linux x86_64
- Compiler: GCC 13 (client)
- Server deps: Go (per `service/go.mod`), `github.com/quic-go/quic-go@v0.54.0`
- Client deps: OpenSSL 3.5.x (script defaults to 3.5.3), pthread
- Network: localhost (127.0.0.1), UDP
- CMake adds `-fsanitize=address,undefined` by default for the client

---

## 5) Notes

- Client follows the non-blocking QUIC demo pattern (`SSL_net_{read,write}_desired`, `SSL_get_event_timeout`, `select`).
- SNI and verification hostname are set; certificate verification is disabled for simplicity (`SSL_VERIFY_NONE`).
- ALPN must match the server: `anAlpnForTest`.
- The server closes all connections in a generation only after all have connected, so the expected behavior is a clean, orderly shutdown with no errors. Any double-free indicates an issue in error-state handling under concurrent QUIC operations.
