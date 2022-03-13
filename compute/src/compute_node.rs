use core::time;
use std::io::Read;
use std::path;

use crate::concurrent_map::{ConcurrentMap, MemTable};
// use crate::error::Result;
use crate::error::Result;
use crate::proto::kvs::key_value_client::*;
use crate::proto::kvs::{key_value_server::*, *};
use crate::proto::node::*;
use crate::proto::status::{self, ErrorCode};
use crate::service::{self, KvServer};
use crate::storage::StorageLayer;
use crate::types::build_status;
use crate::types::Value;
use etcd_client::{LeaseGrantOptions, PutOptions};
use log::{debug, error, info, trace, warn};
use tonic::codegen::http::header::SERVER;
use tonic::{transport::Server, Request, Response, Status};

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
                id: uuid::Uuid::new_v4().to_string(),
                addr: addr,
            },
            meta: meta,
            kv: KvServer::new(),
        }
    }

    pub async fn register(&mut self) -> Result<()> {
        let lease = self.meta.lease_grant(30, None).await?;
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

                tokio::time::sleep(time::Duration::from_secs(15)).await;
            }
        });

        info!("register...");
        self.meta
            .put(
                format!("{}/{}", COMPUTE_SERVICE_PREFIX, self.info.id),
                self.info.addr.as_str(),
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
