# Phase 1 — Networking & Protocols

> Networking is the substrate everything you'll build at Datadog and Apple sits on. Distributed systems are made of computers talking to each other. APIs are computers talking to each other. Microservices are computers talking to each other. If you understand the bytes flowing between them, every higher-level concept becomes obvious.

**Time:** 4–6 days.

**You'll know you're done when:** you can explain, end-to-end, what happens when you type `https://datadog.com` in your browser — including DNS, TCP, TLS, HTTP, and the chain of trust — and you've built a TLS server in Rust from a TCP socket up.

---

## Table of contents

1. [What does this even mean? — "Networking"](#what-networking-means)
2. [Module 1.1 — The OSI and TCP/IP models](#module-11--osi-and-tcpip)
3. [Module 1.2 — IP addressing](#module-12--ip-addressing)
4. [Module 1.3 — TCP vs UDP, deeply](#module-13--tcp-vs-udp)
5. [Module 1.4 — DNS, the phonebook](#module-14--dns)
6. [Module 1.5 — HTTP/1.1, HTTP/2, HTTP/3](#module-15--http)
7. [Module 1.6 — Symmetric vs asymmetric encryption](#module-16--encryption)
8. [Module 1.7 — TLS handshake, byte by byte](#module-17--tls-handshake)
9. [Module 1.8 — HTTPS = HTTP + TLS](#module-18--https)
10. [Module 1.9 — Tools: curl, dig, openssl, wireshark, nc, tcpdump](#module-19--tools)
11. [🛠️ Project: TCP → HTTP → TLS server in Rust](#project-rust-server)
12. [Exercises](#exercises)
13. [Interview question bank](#interview-questions)
14. [What you should now know](#what-you-should-now-know)

---

<a name="what-networking-means"></a>
## 🧠 What does this even mean? — "Networking"

A network is just two or more computers connected so they can send each other bytes. That's it. Everything else — Wi-Fi, the internet, your phone calling an API, Datadog ingesting metrics from millions of agents — is layered on top of this one idea.

The tricky part: there are *a lot* of layers. The cable in the wall (or the radio waves in the air) is the bottom. Your Python code making an HTTP request is near the top. Between them, there are roughly 5–7 layers, each layer pretending it's talking directly to its peer on the other machine, while actually handing data down the stack to the layer below.

Why so many layers? **Separation of concerns.** The Wi-Fi card doesn't know what HTTP is. HTTP doesn't know about Wi-Fi vs Ethernet. Each layer adds its own envelope around the data, and the receiver peels them off.

Picture it like Russian nesting dolls:
```
┌────────────────────────────────────────┐
│ Ethernet header  [ IP header  [ TCP    │
│                  [           [ header  │
│                  [           [ [HTTP   │
│                  [           [ [data]] │
│                  [           [        ]│
│                  [          ]         ]│
│                 ]                      │
└────────────────────────────────────────┘
```

When you send a request, it gets wrapped from inside out. When it arrives, it gets unwrapped from outside in.

---

<a name="module-11--osi-and-tcpip"></a>
## Module 1.1 — The OSI and TCP/IP models

### The two mental models

There are two layered models you'll hear about. Both are descriptive, not prescriptive.

#### OSI (7 layers) — academic, useful for vocabulary

```
Layer 7 - Application      HTTP, DNS, gRPC, SSH, IMAP
Layer 6 - Presentation     TLS, character encoding
Layer 5 - Session          (mostly defunct in practice)
Layer 4 - Transport        TCP, UDP
Layer 3 - Network          IP, ICMP, routing
Layer 2 - Data Link        Ethernet, Wi-Fi (MAC addresses)
Layer 1 - Physical         Cables, radio waves, fiber
```

#### TCP/IP (4 layers) — what the internet actually is

```
Application                HTTP, DNS, gRPC, TLS (TLS sits weirdly between L4 and L7)
Transport                  TCP, UDP
Internet                   IP
Link                       Ethernet, Wi-Fi
```

People casually mix the two. When someone says "L7 load balancer," they mean an application-layer LB (it can read HTTP). "L4 LB" = transport-layer (only sees TCP/IP, can't read HTTP).

> 📖 **Definition — Load balancer:** A piece of software/hardware that distributes incoming requests across multiple backend servers. L4 LBs balance based on IP/port; L7 LBs can route by URL path, hostname, headers, etc.

### Why "layered" matters in practice

You can swap any layer without changing the others. Wi-Fi replaced Ethernet. HTTP/3 replaced TCP with UDP-based QUIC. Apple swapped its underlying network stack many times. The application code (HTTP) didn't have to change because the abstraction held.

### A real example: clicking a link

```
You click "datadog.com"
   ↓
Browser (Application): "I need to send GET /. Let's open an HTTPS connection."
   ↓
TLS (Presentation): "First we negotiate encryption. I need a TCP connection."
   ↓
TCP (Transport): "Open connection to 52.x.x.x port 443. SYN, SYN-ACK, ACK."
   ↓
IP (Network): "I'll route packets from your IP to 52.x.x.x via your gateway."
   ↓
Wi-Fi/Ethernet (Link): "I'll send these packets to the router via radio/cable."
   ↓
Physical: "Modulating bits onto an electrical signal..."
```

Each arrow is the data being handed down. On the server side, it gets handed up the same stack in reverse.

---

<a name="module-12--ip-addressing"></a>
## Module 1.2 — IP addressing

### IPv4

A 32-bit number, written as 4 dotted decimals.

```
192.168.1.42
```

That's `11000000.10101000.00000001.00101010` in binary. There are 2^32 ≈ 4.3 billion possible addresses, which is why we ran out and have IPv6.

#### Special ranges

| Range | Meaning |
|---|---|
| `127.0.0.1` (`::1` in v6) | Loopback — talking to yourself |
| `10.0.0.0/8` | Private (your home/office network) |
| `172.16.0.0/12` | Private |
| `192.168.0.0/16` | Private (most home routers) |
| `169.254.0.0/16` | Link-local (auto-assigned when DHCP fails) |
| `0.0.0.0` | "All interfaces" when binding a server |

> 📖 **Definition — CIDR notation:** The `/16` after an IP block means "the first 16 bits are the network part." So `192.168.0.0/16` covers `192.168.0.0` through `192.168.255.255` — about 65k addresses.

### IPv6

128-bit. Eight groups of 4 hex digits, separated by colons. Consecutive groups of zeros can be compressed with `::`.

```
2607:f8b0:4005:080c:0000:0000:0000:200e
2607:f8b0:4005:80c::200e            (compressed)
```

### NAT (Network Address Translation)

Your router has one public IP (e.g., `73.x.x.x`). You probably have 5+ devices behind it. NAT rewrites outgoing packets so they appear to come from the router's IP, and rewrites incoming packets to deliver them to the right device.

> 🎯 **Interview note:** If asked "how does your phone reach the internet?" — answer: it gets a private IP via DHCP, the router NATs outbound traffic to its public IP, packets traverse your ISP, hit the destination, return, get NATed back to your phone.

### `localhost`, `127.0.0.1`, and `0.0.0.0` — the trio that confuses everyone

- `127.0.0.1` — a literal IP that means "this machine." Your OS short-circuits packets to it (never hits the wire).
- `localhost` — a hostname that usually resolves to `127.0.0.1` (via `/etc/hosts`).
- `0.0.0.0` — when **binding a server**, means "listen on all network interfaces." When used as a destination, it means "no particular host."

Most "why can't I connect to my server?" bugs in Docker/k8s are because the server is bound to `127.0.0.1` (only localhost) instead of `0.0.0.0` (any interface).

---

<a name="module-13--tcp-vs-udp"></a>
## Module 1.3 — TCP vs UDP, deeply

Both are transport-layer protocols. Both let processes on different machines exchange data. They differ in *guarantees*.

### TCP — reliable, ordered, connection-based

> 📖 **Definition — TCP (Transmission Control Protocol):** Gives you a *reliable, ordered byte-stream* between two endpoints. If a packet is lost, TCP retransmits. If packets arrive out of order, TCP reorders. The application sees a clean stream of bytes.

#### The TCP three-way handshake

Before any data is exchanged, the two sides synchronize sequence numbers:

```
Client                                          Server
  │                                                │
  │ ────── SYN, seq=X ──────────────────────────► │     "Want to talk?"
  │                                                │
  │ ◄───── SYN-ACK, seq=Y, ack=X+1 ────────────── │     "Sure, you said X"
  │                                                │
  │ ────── ACK, seq=X+1, ack=Y+1 ──────────────► │     "Got your Y"
  │                                                │
  │═══════════ connection established ═══════════ │
  │                                                │
  │ ────── data ──────────────────────────────► │
  │ ◄───── ACK ────────────────────────────────── │
```

> 📖 **Definition — SYN, ACK:** Flags in a TCP header. SYN = "synchronize" (start). ACK = "acknowledge" (got your packet). The first packet is SYN. The reply is SYN+ACK. The third packet is just ACK. Hence "three-way."

Once established, every packet is acknowledged. Lost packet → retransmit. Sliding windows manage flow control (don't overwhelm the receiver). Congestion control (TCP Reno, BBR, Cubic) manages how fast to send to avoid drowning the network.

#### Connection teardown — the four-way handshake

Either side can close. They send `FIN` to mean "I'm done sending."

```
Client                                          Server
  │ ────── FIN ───────────────────────────────► │
  │ ◄───── ACK ────────────────────────────────── │
  │                                                │
  │   (server may still have data to send)         │
  │                                                │
  │ ◄───── FIN ────────────────────────────────── │
  │ ────── ACK ───────────────────────────────► │
```

### UDP — fire and forget

> 📖 **Definition — UDP (User Datagram Protocol):** Send a datagram (single packet). No connection setup. No retransmit. No ordering. The receiver might get it, might not.

Why use UDP if it's worse?

1. **No setup overhead** — one packet, done.
2. **No head-of-line blocking** — TCP must deliver in order, so a lost packet stalls everything behind it. UDP doesn't care.
3. **You handle reliability yourself** — sometimes you want to (real-time games, video calls — old data is useless, drop it).

UDP is used by:
- DNS queries (small, retry at the application layer if needed)
- Video calls, gaming (low latency > reliability)
- DHCP (boot-up address assignment)
- QUIC / HTTP/3 (reinvents TCP's good parts on top of UDP, but does it better)

### Ports

> 📖 **Definition — Port:** A 16-bit number (0–65535) that identifies a specific service on a machine. An IP address gets you to the machine; a port gets you to the right program. `(IP, port)` is a *socket*.

Well-known ports:
- 22 SSH
- 53 DNS
- 80 HTTP
- 443 HTTPS
- 5432 Postgres
- 6379 Redis
- 8080 HTTP alt (common for dev servers)

Ports 0–1023 require root to bind on Unix.

---

<a name="module-14--dns"></a>
## Module 1.4 — DNS, the phonebook

You don't memorize IP addresses. DNS does.

### What happens when you resolve `datadog.com`

```
Your laptop ─► Local resolver (your router or 8.8.8.8)
                  │
                  ├─► Root nameservers (.) → "ask the .com nameservers"
                  │
                  ├─► .com nameservers → "ask datadog.com's authoritative nameserver at ns-xxx.amazonaws.com"
                  │
                  └─► datadog.com authoritative NS → "datadog.com is at 52.x.x.x"
                  
              Resolver caches result, returns to laptop.
```

### Record types

- **A** — hostname → IPv4
- **AAAA** — hostname → IPv6
- **CNAME** — alias (one name → another name)
- **MX** — mail exchange (where to send email for this domain)
- **TXT** — arbitrary text (often used for SPF, DKIM, domain verification)
- **NS** — which nameservers are authoritative for this domain
- **SOA** — start of authority
- **PTR** — reverse lookup (IP → hostname)

### TTL (Time to Live)

Every DNS record has a TTL — how long resolvers should cache it. Lower TTL = faster propagation when you change records, but more lookups. Datadog and Apple set TTLs carefully for their CDN/load balancer setups.

### Try it yourself

```bash
dig datadog.com
dig datadog.com +short
dig datadog.com MX
dig +trace datadog.com    # show the full resolution path
nslookup apple.com 8.8.8.8

# What's in your DNS cache locally?
sudo killall -INFO mDNSResponder   # Mac: dump cache to system log
```

---

<a name="module-15--http"></a>
## Module 1.5 — HTTP/1.1, HTTP/2, HTTP/3

### HTTP, the basics

> 📖 **Definition — HTTP (HyperText Transfer Protocol):** A request/response protocol for the web. The client sends a request (method, path, headers, optional body); the server replies with a response (status code, headers, optional body).

A raw HTTP/1.1 request:

```
GET /api/users/42 HTTP/1.1
Host: api.example.com
User-Agent: curl/8.0
Accept: application/json

```

(Empty line at the end signals end of headers.)

Response:

```
HTTP/1.1 200 OK
Content-Type: application/json
Content-Length: 27

{"id":42,"name":"Kvng"}
```

That's it. HTTP is a text protocol. You can speak it by hand:

```bash
nc datadog.com 80
GET / HTTP/1.1
Host: datadog.com

# (press Enter on a blank line, see the response)
```

### Methods

| Method | Idempotent? | Safe? | Has body? | Use |
|---|---|---|---|---|
| GET | yes | yes | no | read |
| HEAD | yes | yes | no | metadata only |
| POST | no | no | yes | create / arbitrary action |
| PUT | yes | no | yes | replace |
| PATCH | no (typically) | no | yes | partial update |
| DELETE | yes | no | optional | delete |
| OPTIONS | yes | yes | no | CORS preflight, capabilities |

> 📖 **Definition — Idempotent:** Same effect whether you call it once or 100 times. `DELETE /users/42` is idempotent (after the first call, the user is gone; further calls have no additional effect). `POST /users` typically creates a new user each call → not idempotent.

> 📖 **Definition — Safe:** Doesn't modify server state. GET is safe. POST/PUT/DELETE are not.

### Status codes

| Range | Meaning |
|---|---|
| 1xx | Informational (rare; e.g. 101 Switching Protocols for WebSockets) |
| 2xx | Success — 200 OK, 201 Created, 204 No Content |
| 3xx | Redirect — 301 Moved Permanently, 302 Found, 304 Not Modified |
| 4xx | Client error — 400 Bad Request, 401 Unauthorized, 403 Forbidden, 404 Not Found, 409 Conflict, 422 Unprocessable, 429 Too Many Requests |
| 5xx | Server error — 500 Internal Server Error, 502 Bad Gateway, 503 Service Unavailable, 504 Gateway Timeout |

> 🎯 **Interview note:** Know the difference between 401 (no/bad auth) and 403 (auth ok, but you can't do this). And know that 502/503/504 are infra problems, not your code's bug.

### Headers worth knowing

```
Content-Type: application/json
Content-Length: 1234
Authorization: Bearer eyJhbGc...
Accept: application/json
Cache-Control: max-age=3600, public
ETag: "abc123"
If-None-Match: "abc123"             (conditional request)
User-Agent: ...
X-Forwarded-For: 1.2.3.4            (client IP through proxies)
X-Request-ID: uuid                   (for tracing through services)
Set-Cookie: session=...
```

### HTTP/1.1

- One request at a time per TCP connection (or with keep-alive, sequential)
- Text-based
- Head-of-line blocking: if the first response is slow, others wait
- Browsers compensated by opening 6+ TCP connections per origin

### HTTP/2

Released 2015. Same semantics, different wire format.

- **Binary** instead of text (faster to parse, smaller)
- **Multiplexed** — many requests/responses interleaved on one TCP connection
- **Header compression** (HPACK)
- **Server push** (server can send things you didn't ask for; rarely used)

Still suffers from TCP head-of-line blocking — if one packet is lost, *all* multiplexed streams stall.

### HTTP/3

Released as RFC 9114 in 2022. Same semantics again, but on **QUIC** (a UDP-based protocol).

- No TCP HOL blocking — each stream is independent
- Built-in TLS 1.3 (no separate handshake)
- Faster connection setup (0-RTT for repeat connections)
- Better for mobile (connections survive IP changes)

Datadog APM and many observability platforms care a lot about HTTP/3 because mobile/IoT agents benefit hugely.

---

<a name="module-16--encryption"></a>
## Module 1.6 — Symmetric vs asymmetric encryption

This is where most people glaze over. We're going to make it concrete.

### Symmetric encryption — "shared secret"

Both sides have the **same key**. The key encrypts AND decrypts.

```
plaintext --[encrypt with key K]--> ciphertext --[decrypt with key K]--> plaintext
```

Algorithms: **AES** (industry standard), ChaCha20, DES (broken, don't use).

**Pros:**
- Fast. Hardware-accelerated on modern CPUs (AES-NI).
- Strong with a long key (256 bits is overkill).

**Cons:**
- How do you share the key in the first place? If you can't share it secretly, you can't use symmetric encryption securely.

This is the fundamental "key distribution problem." Symmetric alone isn't enough for the open internet.

### Asymmetric encryption — "public/private keys"

You have **two** keys: a public one (give to everyone) and a private one (keep secret). They're mathematically linked: what one encrypts, only the other can decrypt.

```
plaintext --[encrypt with public key]--> ciphertext --[decrypt with private key]--> plaintext
```

Or signing:
```
message --[sign with private key]--> signature --[verify with public key]--> "yes, the holder of the private key signed this"
```

Algorithms: **RSA** (older, still common), **ECDSA** / **Ed25519** (elliptic-curve, faster, smaller keys), **Diffie-Hellman** / **ECDH** (key *agreement*, not encryption).

**Pros:**
- Solves the key distribution problem.
- Enables digital signatures.

**Cons:**
- ~1000x slower than symmetric.
- Useless for encrypting large payloads.

### The genius combo

Real systems use BOTH:
1. Use **asymmetric** to agree on a shared session key safely.
2. Use that shared key with **symmetric** for the actual data (which is fast).

That is exactly how TLS works.

### Hashing — separate from encryption

> 📖 **Definition — Hash function:** Takes any input, produces a fixed-size fingerprint. One-way (you can't reverse it). Same input → same hash. Different input → different hash (with overwhelming probability).

Used for:
- Password storage (bcrypt, argon2)
- File integrity (SHA-256 of a download)
- Git commit IDs
- HMACs (hash-based message authentication)

Algorithms: **SHA-256** (use this), SHA-3, BLAKE2/3. Don't use MD5 or SHA-1 for security.

---

<a name="module-17--tls-handshake"></a>
## Module 1.7 — TLS handshake, byte by byte

The interview classic. We're going through every step.

### TLS 1.3 handshake (modern; what you'll see in 2026)

```
Client                                          Server
  │                                                │
  │ ── ClientHello ──────────────────────────────►│
  │   - TLS version: 1.3                           │
  │   - Random bytes (32 B, "client random")       │
  │   - Cipher suites the client supports          │
  │   - Supported groups (key exchange algos)      │
  │   - Key share: ephemeral public key (ECDHE)    │
  │   - SNI: "datadog.com" (which site, since      │
  │     one IP can host many)                      │
  │                                                │
  │ ◄────────── ServerHello, EncryptedExtensions, │
  │             Certificate, CertificateVerify,    │
  │             Finished ──────────────────────────│
  │   - Server's chosen cipher                     │
  │   - Server's ephemeral public key              │
  │   - Server's TLS certificate (with chain)      │
  │   - Signature proving server owns the cert     │
  │   - "Finished" message MAC over the handshake  │
  │                                                │
  │   At this point both sides:                    │
  │   - Compute the shared secret via ECDHE        │
  │   - Derive session keys via HKDF               │
  │   - Decrypt the rest                           │
  │                                                │
  │ ── Client Finished ──────────────────────────►│
  │ ── Application data (HTTP request) ──────────►│
  │ ◄───────────────────── Application data ─────│
```

That's **1 round-trip** (1-RTT). TLS 1.2 was 2-RTT. With session resumption, TLS 1.3 can do 0-RTT for repeat connections.

### What's actually happening, in plain English

1. **Client says hi** with a list of "ways I'm willing to talk securely" (cipher suites). It also throws in a fresh ephemeral public key for key exchange.

2. **Server says hi back** picking one cipher suite, sending its own ephemeral public key, and presenting its certificate. The certificate is signed by a Certificate Authority (CA) that the client's OS already trusts.

3. **Client verifies the certificate**:
   - Is it signed by a CA in my trust store? (Macs ship with ~150 trusted CAs.)
   - Is it expired?
   - Does the certificate's "Subject Alternative Names" include the hostname I asked for?
   - Is it revoked? (CRL or OCSP check, often skipped or done lazily.)

4. **Both compute the shared secret** using ECDHE (Elliptic Curve Diffie-Hellman Ephemeral). The math: even an attacker who recorded every byte can't compute the secret without one of the private keys.

5. **Both derive the same set of symmetric session keys** from that shared secret, using HKDF (HMAC-based Key Derivation Function).

6. **They start encrypting everything** with AES-GCM or ChaCha20-Poly1305 (symmetric).

### Certificate chain of trust

```
Root CA (in your OS trust store)
  └── Intermediate CA (signed by Root)
       └── Server cert for datadog.com (signed by Intermediate)
```

When the server sends its cert, it usually sends the intermediate too. The client builds the chain up to a root it trusts.

### The "ephemeral" part — Perfect Forward Secrecy

Older TLS used the server's long-term RSA key for key exchange. If the server's key was ever stolen later, an attacker who recorded old traffic could decrypt it.

ECDHE (the E is "ephemeral") generates a fresh keypair *for each connection*, used only for that handshake, then discarded. Even if the server is later compromised, past traffic stays safe. This is **Perfect Forward Secrecy (PFS)**.

### Try it yourself

```bash
# See the handshake
openssl s_client -connect datadog.com:443 -servername datadog.com

# Just the cert details
openssl s_client -connect datadog.com:443 -servername datadog.com </dev/null \
    | openssl x509 -text -noout

# See the cert chain
openssl s_client -connect datadog.com:443 -showcerts </dev/null
```

> 🎯 **Interview note:** "Walk me through what happens when I type `https://...` and press Enter" is THE classic interview question. By the end of this phase you should be able to talk for 10 minutes from this point.

---

<a name="module-18--https"></a>
## Module 1.8 — HTTPS = HTTP + TLS

That's literally it. HTTPS is HTTP, but the TCP connection is wrapped in TLS first. Everything you learned about HTTP applies; the bytes just go through an encrypted tunnel.

The full sequence for `https://datadog.com`:

```
1. DNS:   datadog.com → 52.x.x.x  (UDP to your resolver)
2. TCP:   3-way handshake to 52.x.x.x:443
3. TLS:   1-RTT handshake (cert verification, key agreement)
4. HTTP:  GET / HTTP/2 (or 1.1, or 3)  — encrypted inside TLS
5. HTTP:  200 OK, body bytes — also encrypted
6. TCP:   connection stays open for more requests (keep-alive)
7. TCP:   FIN/ACK eventually
```

When someone says "the cert is valid," they mean step 3 succeeded.

---

<a name="module-19--tools"></a>
## Module 1.9 — Tools for inspecting the network

You'll use these constantly at Datadog (an observability company).

### `curl` — make any HTTP request

```bash
# Basic GET
curl https://api.github.com/users/torvalds

# Verbose: see headers and TLS details
curl -v https://datadog.com

# See ONLY headers
curl -I https://datadog.com

# POST JSON
curl -X POST https://httpbin.org/post \
  -H "Content-Type: application/json" \
  -d '{"name":"kvng"}'

# Follow redirects
curl -L https://google.com

# Save to file
curl -o output.html https://example.com

# Custom headers, including auth
curl -H "Authorization: Bearer $TOKEN" https://api.example.com/me

# Time the request
curl -w "@-" -o /dev/null -s https://datadog.com <<'EOF'
namelookup:   %{time_namelookup}s
connect:      %{time_connect}s
appconnect:   %{time_appconnect}s   (TLS done)
starttransfer: %{time_starttransfer}s
total:        %{time_total}s
EOF
```

### `httpie` — friendlier curl

```bash
http GET https://api.github.com/users/torvalds
http POST httpbin.org/post name=kvng
```

### `dig` / `nslookup` — DNS

```bash
dig datadog.com
dig +short datadog.com
dig MX gmail.com
dig +trace datadog.com    # see every nameserver in the chain
```

### `openssl` — Swiss army knife of crypto

```bash
# Inspect a cert
openssl s_client -connect datadog.com:443 </dev/null \
    | openssl x509 -text -noout

# Check cert expiry
echo | openssl s_client -connect datadog.com:443 -servername datadog.com 2>/dev/null \
    | openssl x509 -noout -dates

# Generate an RSA keypair
openssl genrsa -out private.pem 2048
openssl rsa -in private.pem -pubout -out public.pem

# Generate an Ed25519 keypair (modern)
openssl genpkey -algorithm ed25519 -out ed25519.pem
```

### `nc` (netcat) — TCP/UDP swiss army knife

```bash
# Listen on a port
nc -l 8080

# Connect to a port and send data
echo "GET / HTTP/1.1\r\nHost: example.com\r\n\r\n" | nc example.com 80

# Test if a port is open
nc -zv datadog.com 443
```

### `tcpdump` — packet capture

```bash
# Capture HTTPS traffic on en0
sudo tcpdump -i en0 'tcp port 443' -nn -v

# Save to a file (then open in Wireshark)
sudo tcpdump -i en0 -w capture.pcap 'tcp port 443'
```

### Wireshark — visual packet inspection

GUI tool. Open Wireshark, pick your interface, hit start. Filter `tcp.port == 443`. Watch packets fly. Right-click → Follow → TCP Stream to see a full conversation reassembled. Set up a key log to actually decrypt TLS:

```bash
# Tell curl/firefox/chrome to log TLS keys
export SSLKEYLOGFILE=$HOME/sslkeys.log
# In Wireshark: Preferences → Protocols → TLS → "(Pre)-Master-Secret log filename" = $HOME/sslkeys.log
# Now Wireshark can decrypt your TLS traffic.
```

---

<a name="project-rust-server"></a>
## 🛠️ Project: TCP → HTTP → TLS server in Rust

A three-stage project where you build a server up the layers.

**See `projects/tcp-to-tls-server/` for full code.**

### Why Rust?

Rust forces you to think about bytes, lifetimes, and ownership. You can't lazily allocate everything to a string-based abstraction like Python. When you parse an HTTP request in Rust, you *feel* where the buffer ends and the next request begins — exactly the kind of intuition you need for systems work.

### What you'll build

Three binaries in one Cargo workspace:

1. **`stage1-tcp-echo`** — a raw TCP server. Accept connections, read bytes, echo them back. ~50 lines.
2. **`stage2-http`** — same TCP server, but parses HTTP/1.1 requests and returns proper HTTP responses. Routes `/`, `/healthz`, `/echo`. ~200 lines.
3. **`stage3-tls`** — same HTTP server wrapped in TLS using `rustls`. Generate a self-signed cert, watch the handshake in Wireshark. ~250 lines.

See `projects/tcp-to-tls-server/README.md` for setup.

---

<a name="exercises"></a>
## Exercises

1. **DNS detective.** Use `dig` to find: (a) the MX records for `gmail.com`; (b) the nameservers for `apple.com`; (c) all the A records returned for `cnn.com` (likely many — load balancing). Write a one-paragraph explanation of what you found.

2. **Decode an HTTP request manually.** Run `nc -l 8080` in one terminal. In another, `curl http://localhost:8080/api/users?id=42 -H "X-Test: hello" -d '{"name":"k"}'`. Read the request `nc` printed out and label each part: request line, headers, body.

3. **Watch a TLS handshake.** Open Wireshark, filter `tcp.port == 443`. In another terminal `curl https://datadog.com`. Find the ClientHello, ServerHello, Certificate. Identify the cipher suite chosen.

4. **Cert dive.** Use `openssl s_client` to grab Datadog's cert. Identify: who issued it, when it expires, what hostnames it covers (Subject Alternative Names), what algorithm signed it.

5. **Implement a TCP "knock-knock"** in Rust: a server that, when you connect, expects you to send `KNOCK\n`, then replies `WHO\n`, then you send `KVNG\n`, then it replies `HELLO\n` and closes. Use `tokio` async or std's blocking — your call.

6. **Load test your stage2 HTTP server** with `wrk` (`brew install wrk`):
   ```
   wrk -t4 -c100 -d10s http://localhost:8080/healthz
   ```
   Note req/sec. What's the bottleneck? (Single-threaded? Allocation per request? Try optimizing.)

---

<a name="interview-questions"></a>
## 🎯 Interview question bank

These show up at Apple, Datadog, and most senior-engineer-track companies. Practice saying the answers out loud.

1. **What happens when I type `https://google.com` in the browser and press Enter?**
   *(Practice the full 5-minute version.)*

2. **TCP vs UDP — pick one for: a video call, a database connection, DNS, a file download. Why?**

3. **Walk me through the TLS handshake.**

4. **Why do we use both symmetric and asymmetric crypto in TLS?**

5. **What's the difference between HTTP/1.1, HTTP/2, and HTTP/3?**

6. **Explain Perfect Forward Secrecy.**

7. **What's a CIDR block? What does `/24` mean?**

8. **Difference between 401 and 403?**

9. **Difference between a 502 and a 504?**

10. **What is SNI and why does it matter?**
    *(Hint: one IP, many sites. The cert presented depends on which hostname you asked for.)*

11. **What's a CNAME, and why can't the apex of a domain be a CNAME?**
    *(Hint: DNS rules — apex must have an SOA, CNAME forbids any other records.)*

12. **What's HSTS?**
    *(HTTP Strict Transport Security — header that tells browsers "always use HTTPS for me.")*

13. **What's a load balancer? L4 vs L7?**

14. **What's a CDN, what does it actually do?**
    *(Edge servers, geo-routing via DNS or anycast, caching static + sometimes dynamic content.)*

15. **What's CORS, why does it exist, who enforces it?**
    *(Browser-enforced. Cross-Origin Resource Sharing. Without it, a page on attacker.com could `fetch()` your bank's API using your cookies.)*

---

<a name="what-you-should-now-know"></a>
## ✅ What you should now know

- [ ] Layered network model (OSI, TCP/IP) and why layering matters
- [ ] IP addresses, CIDR, NAT, localhost trio
- [ ] TCP three-way handshake, four-way close
- [ ] TCP vs UDP trade-offs
- [ ] DNS resolution path, record types
- [ ] HTTP methods, status codes, headers
- [ ] HTTP/1.1 vs HTTP/2 vs HTTP/3 trade-offs
- [ ] Symmetric vs asymmetric encryption, why TLS uses both
- [ ] TLS 1.3 handshake step by step
- [ ] Certificate chain of trust
- [ ] Perfect Forward Secrecy
- [ ] How to use `dig`, `curl`, `openssl`, `nc`, `tcpdump`, Wireshark
- [ ] You've built a TCP server, then HTTP, then TLS in Rust

---

**Next:** [Phase 2 — OOP & Design Patterns in Go](../phase-02-oop-go/README.md)
