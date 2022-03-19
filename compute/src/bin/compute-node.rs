use std::fs::File;
use std::io::Write;

use clap::Parser;

use compute_node::error::Result;
use compute_node::ComputeNode;
use log::{info, Log};

#[derive(Parser, Debug)]
struct Args {
    #[clap(short, long)]
    endpoints: Vec<String>,

    #[clap(short, long)]
    addr: String,
}

#[derive(Debug)]
struct LogTarget {
    stdout: bool,
    file: File,
}

impl LogTarget {
    pub fn new(stdout: bool, filename: &str) -> Self {
        let file = File::create(filename).unwrap();
        return LogTarget { stdout, file };
    }
}

impl Write for LogTarget {
    fn write(&mut self, buf: &[u8]) -> std::io::Result<usize> {
        if self.stdout {
            std::io::stdout().write(buf)?;
        }
        let n = self.file.write(buf).unwrap();
        self.file.flush().unwrap();
        Ok(n)
    }

    fn flush(&mut self) -> std::io::Result<()> {
        if self.stdout {
            std::io::stdout().flush()?;
        }
        Ok(self.file.flush().unwrap())
    }
}

#[tokio::main]
async fn main() -> Result<()> {
    env_logger::builder()
        .filter_level(log::LevelFilter::Info)
        .target(env_logger::fmt::Target::Stderr)
        .init();

    let args = Args::parse();
    info!("args={:?}", args);

    let meta = etcd_client::Client::connect(args.endpoints, None).await?;
    info!("etcd connected");

    let mut node = ComputeNode::new(meta, args.addr);
    node.register().await?;
    node.start().await?;

    Ok(())
}
