# Glossary

Every term defined across the curriculum, in alphabetical order. When you see a `📖 Definition` callout in any phase, the term ends up here.

---

### API (Application Programming Interface)
A contract between two pieces of software. One side exposes functions/endpoints; the other side calls them. "API" usually means a network API (REST, gRPC) but it also means a library API (the methods you call on a class).

### ACID
Four properties a database transaction should have: **Atomicity** (all or nothing), **Consistency** (the database stays valid), **Isolation** (concurrent transactions don't trample each other), **Durability** (once committed, it survives crashes).

### Asymmetric encryption
Encryption that uses two different keys: a **public key** anyone can have, and a **private key** only you have. What one encrypts, only the other can decrypt. Slow but solves "how do we agree on a secret over an insecure channel?" Used in TLS handshake.

### Bash
The shell (command-line program) most Linux/Mac users interact with. You type commands, it runs them. Also a scripting language.

### CI/CD (Continuous Integration / Continuous Delivery)
**CI**: Every time you push code, automated tests run. **CD**: If tests pass, the code automatically deploys (or is automatically *ready* to deploy). Goal: catch bugs fast, ship often.

### Concurrency
Multiple things *making progress* at the same time. Doesn't necessarily mean simultaneously — could be fast-switching on one CPU. (Contrast with parallelism.)

### Connection pool
A reusable set of pre-opened database connections. Opening a new DB connection is expensive (network round-trips, auth). A pool keeps N of them alive and hands them out on demand.

### Container
A lightweight, isolated environment for running an application. Think "tiny VM that shares the host's kernel." Docker is the most famous container runtime.

### Distributed system
A set of computers that coordinate over a network to look (to a user) like one system. Examples: Google Search, Datadog itself, your bank's website.

### DNS (Domain Name System)
The phonebook of the internet. Translates human-readable names (`google.com`) to IP addresses (`142.250.190.46`).

### Embedded system
A computer that lives inside something that isn't a "computer" — your car's ECU, a microwave, a smart watch. Usually has tight memory/CPU/power constraints.

### Git
A version control system. Tracks changes to files over time, lets multiple people collaborate without overwriting each other.

### GitHub / GitLab
Websites that host Git repositories and add collaboration features (pull requests, code review, CI/CD pipelines). Datadog uses **GitLab** internally.

### gRPC
A way for services to call each other over a network, using **Protocol Buffers** (protobuf) for the message format and HTTP/2 for transport. Faster and more typed than REST. Datadog uses it heavily for internal services.

### HTTP (HyperText Transfer Protocol)
The protocol the web speaks. Request/response over TCP. Versions: HTTP/1.1 (text-based), HTTP/2 (binary, multiplexed), HTTP/3 (over UDP via QUIC).

### HTTPS
HTTP wrapped in TLS encryption. The `S` is "Secure."

### IP address
A number that identifies a computer on a network. IPv4: `142.250.190.46`. IPv6: `2607:f8b0:4005:80c::200e`.

### Kafka
A distributed message queue / event streaming platform. Producers write to "topics," consumers read from them. Built for huge throughput. Datadog uses Kafka extensively for ingesting metrics.

### Kubernetes (k8s)
An orchestrator for containers. You tell it "I want 3 copies of this app running, healthy, behind a load balancer," and it makes that happen.

### Linux
An open-source operating system kernel. Most servers run Linux. macOS isn't Linux but is Unix-like, so most commands work the same.

### Parallelism
Multiple things happening *literally simultaneously*, on multiple CPU cores. A subset of concurrency.

### Process
A running program, with its own memory space. Two processes can't accidentally read each other's memory.

### Protocol
A set of rules for how two computers should talk. HTTP is a protocol. TCP is a protocol. Even shaking hands is a "protocol."

### Repository (repo)
A folder tracked by Git. Has a `.git` subdirectory containing the history.

### Shell
A program that takes typed commands and runs them. Bash, zsh, fish — all shells. Different syntaxes, similar idea.

### Symmetric encryption
Encryption where the same key encrypts and decrypts. Fast. Problem: how do you share the key without someone snooping? (Answer: use asymmetric to share it. See **TLS handshake**.)

### TCP (Transmission Control Protocol)
A protocol that gives you a reliable, ordered byte-stream between two computers. If a packet is lost, TCP retransmits. HTTP runs on TCP.

### TLS (Transport Layer Security)
The protocol that makes HTTPS secure. Encrypts data in transit, verifies server identity via certificates. Successor to SSL.

### Thread
A path of execution within a process. Multiple threads in one process **share memory** — which is powerful but also the source of every concurrency bug ever.

### UDP (User Datagram Protocol)
Like TCP but unreliable and unordered — you fire packets and hope they arrive. Faster, used for video calls, DNS, gaming.

### YAML
A human-readable config file format. Indentation matters. Used by Kubernetes, GitLab CI, GitHub Actions, Docker Compose, etc.

---

*This glossary grows with every phase. When you forget a term, come here.*
