// Async TCP port scanner. Demonstrates tokio + Semaphore for backpressure.
//
// Run: cargo run -- --host scanme.nmap.org --start 1 --end 1024 --conc 200

use clap::Parser;
use std::sync::Arc;
use std::time::Duration;
use tokio::net::TcpStream;
use tokio::sync::Semaphore;
use tokio::time::timeout;

#[derive(Parser, Debug)]
struct Args {
    #[arg(long)]
    host: String,
    #[arg(long, default_value_t = 1)]
    start: u16,
    #[arg(long, default_value_t = 1024)]
    end: u16,
    #[arg(long, default_value_t = 200)]
    conc: usize,
    #[arg(long, default_value_t = 1000)]
    timeout_ms: u64,
}

#[tokio::main]
async fn main() {
    let args = Args::parse();

    // Semaphore caps how many simultaneous connect attempts can be outstanding.
    // Without this, scanning 60K ports = 60K open sockets = OS limits hit hard.
    let sem = Arc::new(Semaphore::new(args.conc));

    let host = Arc::new(args.host);
    let to = Duration::from_millis(args.timeout_ms);
    let mut handles = Vec::new();

    for port in args.start..=args.end {
        let permit = sem.clone().acquire_owned().await.unwrap();
        let host = host.clone();
        let h = tokio::spawn(async move {
            let _permit = permit; // held until task completes
            let addr = format!("{}:{}", host, port);
            match timeout(to, TcpStream::connect(&addr)).await {
                Ok(Ok(_)) => Some(port),     // open
                _ => None,                    // closed or timeout
            }
        });
        handles.push(h);
    }

    let mut open = Vec::new();
    for h in handles {
        if let Ok(Some(p)) = h.await {
            open.push(p);
            println!("OPEN  {}", p);
        }
    }
    println!("\n{} open of {} scanned", open.len(),
        args.end as i32 - args.start as i32 + 1);
}
