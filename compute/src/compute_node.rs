use core::time;

use crate::error::Result;
use crate::id_generator::next_id;
use crate::proto::kvs::key_value_server::*;
use crate::proto::node::*;
use crate::service::KvServer;

use etcd_client::PutOptions;
use log::{debug, info};

use tonic::transport::Server;

const SERVICE_PREFIX: &str = "services";
const COMPUTE_SERVICE_PREFIX: &str = "services/compute";

pub struct ComputeNode {
    // data
    info: NodeInfo,
    meta: etcd_client::Client,

    // services
    kv: KvServer,
}

impl ComputeNode {
    pub fn new(meta: etcd_client::Client, addr: String) -> Self {
        Self {
            info: NodeInfo {
                id: next_id(),
                addr: addr,
            },
            meta: meta,

            kv: KvServer::new(),
        }
    }

    pub async fn register(&mut self) -> Result<()> {
        let lease = self.meta.lease_grant(10, None).await?;
        let (mut keeper, mut stream) = self.meta.lease_keep_alive(lease.id()).await?;

        tokio::spawn(async move {
            loop {
                keeper.keep_alive().await.unwrap();
                if let Some(resp) = stream.message().await.unwrap() {
                    debug!(
                        "keep alive stream message: id={} ttl={}",
                        resp.id(),
                        resp.ttl()
                    )
                }

                tokio::time::sleep(time::Duration::from_secs(3)).await;
            }
        });

        info!("register...");
        let info = serde_json::to_string(&self.info).unwrap();
        self.meta
            .put(
                format!("{}/{}", COMPUTE_SERVICE_PREFIX, self.info.id),
                info,
                Some(PutOptions::new().with_lease(lease.id())),
            )
            .await?;

        Ok(())
    }

    pub async fn start(self) -> Result<()> {
        info!("start service...");
        Server::builder()
            .add_service(KeyValueServer::new(self.kv))
            .serve(self.info.addr.parse()?)
            .await?;

        Ok(())
    }
}
