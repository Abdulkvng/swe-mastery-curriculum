// stage2-http/src/main.rs
//
// An HTTP/1.1 server, hand-rolled. No frameworks.
// We parse the request line, headers, and body ourselves.
//
// HTTP/1.1 wire format (the EXACT bytes on the wire):
//
//   GET /healthz HTTP/1.1\r\n
//   Host: localhost:7878\r\n
//   User-Agent: curl/8.0\r\n
//   Accept: */*\r\n
//   \r\n
//   <body bytes if any>
//
// The empty line (`\r\n\r\n`) separates headers from body.
// The Content-Length header tells us how many body bytes to read.

use std::collections::HashMap;
use std::io::{BufRead, BufReader, Read, Write};
use std::net::{TcpListener, TcpStream};
use std::thread;

fn main() -> std::io::Result<()> {
    let addr = "0.0.0.0:7878";
    let listener = TcpListener::bind(addr)?;
    println!("[stage2] listening on {addr}");

    for stream in listener.incoming() {
        match stream {
            Ok(stream) => {
                thread::spawn(move || {
                    if let Err(e) = handle_client(stream) {
                        eprintln!("[stage2] client error: {e}");
                    }
                });
            }
            Err(e) => eprintln!("[stage2] accept error: {e}"),
        }
    }

    Ok(())
}

// A parsed HTTP request.
struct Request {
    method: String,
    path: String,
    #[allow(dead_code)] // we don't act on version, but parse it
    version: String,
    headers: HashMap<String, String>,
    body: Vec<u8>,
}

// A response we'll serialize to bytes.
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

    // Write the response to a stream in HTTP/1.1 wire format.
    fn write_to<W: Write>(&self, w: &mut W) -> std::io::Result<()> {
        // Status line
        write!(w, "HTTP/1.1 {} {}\r\n", self.status, self.status_text)?;
        // Headers
        for (k, v) in &self.headers {
            write!(w, "{k}: {v}\r\n")?;
        }
        // Blank line separating headers from body
        write!(w, "\r\n")?;
        // Body
        w.write_all(&self.body)?;
        Ok(())
    }
}

fn handle_client(stream: TcpStream) -> std::io::Result<()> {
    let peer = stream.peer_addr()?;
    println!("[stage2] connected: {peer}");

    // Wrap in BufReader for line-based reading.
    // BufReader buffers reads internally, so read_line doesn't make a syscall per byte.
    let mut reader = BufReader::new(stream.try_clone()?);
    let mut writer = stream;

    let req = match parse_request(&mut reader)? {
        Some(r) => r,
        None => {
            // Client closed before sending anything. Just bail.
            return Ok(());
        }
    };

    println!("[stage2] {} {} {}", req.method, req.path, req.version);

    let resp = route(&req);
    resp.write_to(&mut writer)?;
    writer.flush()?;
    // We send Connection: close, so we just close. No keep-alive logic.
    Ok(())
}

// Parse an HTTP/1.1 request from a BufReader.
// Returns None if the stream is empty (client disconnected).
fn parse_request(reader: &mut BufReader<TcpStream>) -> std::io::Result<Option<Request>> {
    // === Request line ===
    // e.g. "GET /healthz HTTP/1.1\r\n"
    let mut line = String::new();
    let n = reader.read_line(&mut line)?;
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

    // === Headers ===
    // Read lines until we hit an empty line (\r\n alone).
    let mut headers: HashMap<String, String> = HashMap::new();
    loop {
        let mut header_line = String::new();
        let n = reader.read_line(&mut header_line)?;
        if n == 0 {
            break; // EOF mid-headers; treat as end.
        }
        let trimmed = header_line.trim_end_matches("\r\n").trim_end_matches('\n');
        if trimmed.is_empty() {
            break; // end of headers
        }
        if let Some((k, v)) = trimmed.split_once(':') {
            // Header names are case-insensitive per RFC. Lowercase for lookup.
            headers.insert(k.trim().to_ascii_lowercase(), v.trim().to_string());
        }
    }

    // === Body ===
    // If Content-Length is set, read exactly that many bytes.
    // (We're ignoring chunked encoding for brevity. Real servers handle it.)
    let mut body = Vec::new();
    if let Some(cl) = headers.get("content-length") {
        if let Ok(len) = cl.parse::<usize>() {
            body.resize(len, 0);
            reader.read_exact(&mut body)?;
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

// Tiny router: dispatch on (method, path).
fn route(req: &Request) -> Response {
    match (req.method.as_str(), req.path.as_str()) {
        ("GET", "/") => Response::ok("Hello from stage2-http! Try /healthz, /echo\n"),
        ("GET", "/healthz") => Response::ok("ok\n"),
        ("POST", "/echo") => {
            // Echo the body back.
            Response::ok(req.body.clone())
        }
        ("GET", "/headers") => {
            // Show the request headers, useful for debugging.
            let mut s = String::new();
            for (k, v) in &req.headers {
                s.push_str(&format!("{k}: {v}\n"));
            }
            Response::ok(s)
        }
        _ => Response::not_found(),
    }
}
