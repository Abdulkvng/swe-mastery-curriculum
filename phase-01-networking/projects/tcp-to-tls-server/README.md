# tcp-to-tls-server

> Build a server up the network stack: raw TCP → HTTP/1.1 → HTTPS (with TLS).

## Why three stages?

Each stage adds exactly one new concept on top of the last. By stage 3 you'll have a real working HTTPS server you wrote yourself, and you'll deeply understand each layer.

## Project layout

```
tcp-to-tls-server/
├── Cargo.toml          (workspace root)
├── stage1-tcp-echo/    Raw TCP echo server
├── stage2-http/        HTTP/1.1 server (parses requests, routes paths)
└── stage3-tls/         Same HTTP server, wrapped in TLS via rustls
```

## Prereqs

```bash
brew install rust
rustc --version  # should be 1.75+
```

## Build everything

```bash
cd tcp-to-tls-server
cargo build --release
```

## Run each stage

```bash
# Stage 1: in one terminal
cargo run -p stage1-tcp-echo

# In another:
nc localhost 7878
hello
hello                    # echoed back
^D                       # close

# Stage 2:
cargo run -p stage2-http
curl http://localhost:7878/
curl http://localhost:7878/healthz
curl -X POST http://localhost:7878/echo -d "hi from kvng"

# Stage 3: first generate a self-signed cert
cd stage3-tls
./gen-cert.sh
cargo run

# In another terminal (-k = ignore cert validity since it's self-signed):
curl -k https://localhost:7878/
curl -k https://localhost:7878/healthz

# Watch the TLS handshake:
openssl s_client -connect localhost:7878 -servername localhost </dev/null
```

## What you should learn

### Stage 1
- Rust's `std::net::TcpListener` and `TcpStream`
- Reading bytes from a socket into a buffer
- The accept loop pattern (`loop { listener.accept() }`)
- Spawning a thread per connection
- What `nc` actually does

### Stage 2
- Parsing a text protocol (HTTP/1.1) by hand
- The request line, headers, body separation by `\r\n\r\n`
- Routing based on path
- Writing a Content-Length-aware response
- Why HTTP/1.1 keep-alive matters (we don't implement it — note the cost)
- Buffer ownership: who owns the bytes, when do they get freed

### Stage 3
- `rustls` — the modern Rust TLS library
- Loading a cert + private key
- Wrapping a `TcpStream` in a `rustls::ServerConnection`
- The handshake happens automatically on first read/write
- Watching the handshake in Wireshark / `openssl s_client`
- Why self-signed certs trigger warnings (no CA in trust chain)
