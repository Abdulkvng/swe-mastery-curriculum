// stage1-tcp-echo/src/main.rs
//
// A raw TCP echo server. The simplest possible thing.
// Whatever bytes a client sends, we send them back.
//
// Concepts you should internalize:
//   - TcpListener "binds" to (IP, port). 0.0.0.0 = all interfaces.
//   - listener.accept() blocks until a client connects, then returns a TcpStream.
//   - A TcpStream implements Read + Write — you treat it like a file.
//   - We spawn a thread per connection so multiple clients can be served.
//     (At Datadog scale you'd use async / tokio instead — coming in stage 3.)

use std::io::{Read, Write};
use std::net::{TcpListener, TcpStream};
use std::thread;

fn main() -> std::io::Result<()> {
    // Bind to 0.0.0.0:7878. 7878 = "rust" on a phone keypad. Cute.
    let addr = "0.0.0.0:7878";
    let listener = TcpListener::bind(addr)?;
    println!("[stage1] listening on {addr}");

    // Accept loop. Every iteration accepts one connection.
    // listener.incoming() yields Result<TcpStream> values.
    for stream in listener.incoming() {
        match stream {
            Ok(stream) => {
                // Spawn a thread to handle this connection so we can keep
                // accepting more in parallel.
                thread::spawn(move || {
                    if let Err(e) = handle_client(stream) {
                        eprintln!("[stage1] client error: {e}");
                    }
                });
            }
            Err(e) => {
                eprintln!("[stage1] accept error: {e}");
            }
        }
    }

    Ok(())
}

// Handle one connection: read bytes, echo them back, close when client closes.
fn handle_client(mut stream: TcpStream) -> std::io::Result<()> {
    let peer = stream.peer_addr()?;
    println!("[stage1] connected: {peer}");

    // 4 KiB buffer. In production you'd be more careful (BufReader, max sizes, etc.)
    let mut buf = [0u8; 4096];

    loop {
        // Read up to buf.len() bytes. Returns number of bytes read (0 = EOF).
        let n = stream.read(&mut buf)?;
        if n == 0 {
            // Client closed the connection.
            println!("[stage1] disconnected: {peer}");
            return Ok(());
        }

        // Echo back exactly the bytes we got.
        // write_all loops on partial writes until the whole slice is sent.
        stream.write_all(&buf[..n])?;
    }
}
