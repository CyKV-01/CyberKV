use clap::Parser;

use compute_node::error::Result;
use compute_node::ComputeNode;
use log::info;

#[derive(Parser, Debug)]
struct Args {
    #[clap(short, long)]
    endpoints: Vec<String>,

    #[clap(short, long)]
    addr: String,
}

#[tokio::main]
async fn main() -> Result<()> {
    env_logger::builder()
        .filter_level(log::LevelFilter::Info)
        .target(env_logger::fmt::Target::Stdout)
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
