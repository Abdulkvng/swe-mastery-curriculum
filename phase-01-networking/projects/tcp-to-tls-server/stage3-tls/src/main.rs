// stage3-tls/src/main.rs
//
// Same HTTP server as stage 2, but every connection is wrapped in TLS.
//
// What changes from stage 2:
//   - Before reading/writing HTTP, we wrap the TcpStream in a rustls::StreamOwned.
//   - rustls handles the TLS handshake automatically on first read/write.
//   - We need a certificate + private key. See gen-cert.sh.
//
// What stays the same:
//   - HTTP parsing logic is byte-for-byte identical.
//   - The application doesn't know or care that it's encrypted. That's the
//     beauty of layering.
//
// To run:
//   ./gen-cert.sh
//   cargo run
//
// To test:
//   curl -k https://localhost:7878/         # -k = ignore cert (it's self-signed)
//   openssl s_client -connect localhost:7878 </dev/null

use std::collections::HashMap;
use std::fs::File;
use std::io::{BufRead, BufReader, Read, Write};
use std::net::{TcpListener, TcpStream};
use std::sync::Arc;
use std::thread;

use rustls::pki_types::{CertificateDer, PrivateKeyDer};
use rustls::{ServerConfig, ServerConnection, StreamOwned};

fn main() -> std::io::Result<()> {
    // === Load TLS config ===
    // ServerConfig is shared across all connections, so we wrap in Arc (atomic ref-counted).
    let config = Arc::new(load_tls_config()?);

    let addr = "0.0.0.0:7878";
    let listener = TcpListener::bind(addr)?;
    println!("[stage3] listening on {addr} (TLS)");
    println!("[stage3] try: curl -k https://localhost:7878/");

    for stream in listener.incoming() {
        match stream {
            Ok(tcp_stream) => {
                let config = config.clone();
                thread::spawn(move || {
                    if let Err(e) = handle_client(tcp_stream, config) {
                        eprintln!("[stage3] client error: {e}");
                    }
                });
            }
            Err(e) => eprintln!("[stage3] accept error: {e}"),
        }
    }

    Ok(())
}

// Build a rustls ServerConfig from cert + key files on disk.
fn load_tls_config() -> std::io::Result<ServerConfig> {
    // Read cert.pem (the PEM-encoded certificate chain).
    let cert_file = File::open("cert.pem").map_err(|e| {
        std::io::Error::new(
            e.kind(),
            "cert.pem not found. Run ./gen-cert.sh first.",
        )
    })?;
    let mut cert_reader = BufReader::new(cert_file);
    let certs: Vec<CertificateDer<'static>> = rustls_pemfile::certs(&mut cert_reader)
        .collect::<Result<Vec<_>, _>>()?;

    // Read key.pem (the PEM-encoded private key).
    let key_file = File::open("key.pem").map_err(|e| {
        std::io::Error::new(e.kind(), "key.pem not found. Run ./gen-cert.sh first.")
    })?;
    let mut key_reader = BufReader::new(key_file);
    // private_key returns the first key found (PKCS#8, RSA, or EC).
    let key: PrivateKeyDer<'static> = rustls_pemfile::private_key(&mut key_reader)?
        .ok_or_else(|| std::io::Error::new(std::io::ErrorKind::InvalidData, "no key in key.pem"))?;

    // Build the config. with_no_client_auth = we don't require client certs (mTLS).
    let config = ServerConfig::builder()
        .with_no_client_auth()
        .with_single_cert(certs, key)
        .map_err(|e| std::io::Error::new(std::io::ErrorKind::InvalidData, e.to_string()))?;

    Ok(config)
}

fn handle_client(tcp_stream: TcpStream, config: Arc<ServerConfig>) -> std::io::Result<()> {
    let peer = tcp_stream.peer_addr()?;
    println!("[stage3] tcp connected: {peer}");

    // Create a TLS connection bound to this server config.
    let tls_conn = ServerConnection::new(config)
        .map_err(|e| std::io::Error::new(std::io::ErrorKind::Other, e.to_string()))?;

    // StreamOwned ties the TLS state and TCP stream together.
    // It implements Read + Write, so the rest of the code is identical to stage 2.
    let mut tls_stream = StreamOwned::new(tls_conn, tcp_stream);

    // The handshake happens lazily on the first read or write.
    // Once it completes, all I/O is encrypted transparently.

    let req = match parse_request(&mut tls_stream)? {
        Some(r) => r,
        None => return Ok(()),
    };

    println!("[stage3] {} {} {}", req.method, req.path, req.version);

    let resp = route(&req);
    resp.write_to(&mut tls_stream)?;
    tls_stream.flush()?;
    Ok(())
}

// === Below is identical to stage 2 — proves layering works. ===

struct Request {
    method: String,
    path: String,
    version: String,
    headers: HashMap<String, String>,
    body: Vec<u8>,
}

struct Response {
    status: u16,
    status_text: &'static str,
    headers: Vec<(String, String)>,
    body: Vec<u8>,
}

impl Response {
    fn ok(body: impl Into<Vec<u8>>) -> Self {
        let body = body.into();
        Response {
            status: 200,
            status_text: "OK",
            headers: vec![
                ("Content-Type".into(), "text/plain; charset=utf-8".into()),
                ("Content-Length".into(), body.len().to_string()),
                ("Connection".into(), "close".into()),
            ],
            body,
        }
    }

    fn not_found() -> Self {
        let body = b"404 not found\n".to_vec();
        Response {
            status: 404,
            status_text: "Not Found",
            headers: vec![
                ("Content-Type".into(), "text/plain".into()),
                ("Content-Length".into(), body.len().to_string()),
                ("Connection".into(), "close".into()),
            ],
            body,
        }
    }

    fn write_to<W: Write>(&self, w: &mut W) -> std::io::Result<()> {
        write!(w, "HTTP/1.1 {} {}\r\n", self.status, self.status_text)?;
        for (k, v) in &self.headers {
            write!(w, "{k}: {v}\r\n")?;
        }
        write!(w, "\r\n")?;
        w.write_all(&self.body)?;
        Ok(())
    }
}

fn parse_request<R: Read>(reader: &mut R) -> std::io::Result<Option<Request>> {
    let mut buf = BufReader::new(reader);

    let mut line = String::new();
    let n = buf.read_line(&mut line)?;
    if n == 0 {
        return Ok(None);
    }

    let line = line.trim_end_matches("\r\n").trim_end_matches('\n');
    let parts: Vec<&str> = line.splitn(3, ' ').collect();
    if parts.len() != 3 {
        return Err(std::io::Error::new(
            std::io::ErrorKind::InvalidData,
            format!("bad request line: {line:?}"),
        ));
    }
    let method = parts[0].to_string();
    let path = parts[1].to_string();
    let version = parts[2].to_string();

    let mut headers: HashMap<String, String> = HashMap::new();
    loop {
        let mut header_line = String::new();
        let n = buf.read_line(&mut header_line)?;
        if n == 0 {
            break;
        }
        let trimmed = header_line.trim_end_matches("\r\n").trim_end_matches('\n');
        if trimmed.is_empty() {
            break;
        }
        if let Some((k, v)) = trimmed.split_once(':') {
            headers.insert(k.trim().to_ascii_lowercase(), v.trim().to_string());
        }
    }

    let mut body = Vec::new();
    if let Some(cl) = headers.get("content-length") {
        if let Ok(len) = cl.parse::<usize>() {
            body.resize(len, 0);
            buf.read_exact(&mut body)?;
        }
    }

    Ok(Some(Request {
        method,
        path,
        version,
        headers,
        body,
    }))
}

fn route(req: &Request) -> Response {
    match (req.method.as_str(), req.path.as_str()) {
        ("GET", "/") => Response::ok("Hello from stage3-tls (HTTPS)!\n"),
        ("GET", "/healthz") => Response::ok("ok\n"),
        ("POST", "/echo") => Response::ok(req.body.clone()),
        _ => Response::not_found(),
    }
}
