use std::collections::HashMap;

use futures::future::join_all;
use tokio::sync::RwLock;
use tonic::transport::Channel;

use crate::{
    error::Result,
    kvs::Kvs,
    proto::{
        kv::{key_value_client::*, ReadResponse},
        status::ErrorCode,
    },
    types::{calc_slot, SlotID, Status},
};

type NodeID = u64;

#[derive(Debug)]
pub struct NodeInfo {
    id: NodeID,
    addr: String,
}

#[derive(Debug)]
pub struct StorageNode {
    info: NodeInfo,
    client: Option<KeyValueClient<Channel>>,
}

impl StorageNode {
    pub async fn get(&self, key: &str) -> Result<Option<ReadResponse>> {
        todo!()
    }

    pub async fn set(&self, key: String, value: String) -> Result<Option<String>> {
        todo!()
    }

    pub async fn remove(&self, key: &str) -> Result<Option<String>> {
        todo!()
    }
}

#[derive(Debug)]
pub struct StorageLayer {
    replica_num: usize,
    write_quorum: usize,
    read_quorum: usize,
    slot_node_index: RwLock<HashMap<SlotID, Vec<StorageNode>>>,
}

impl StorageLayer {
    pub async fn get(&self, key: &str) -> Result<Option<String>> {
        let slot = calc_slot(key);
        let slot_node_index = self.slot_node_index.read().await;
        let nodes = slot_node_index.get(&slot);

        match nodes {
            Some(nodes) => {
                if nodes.len() < self.read_quorum {
                    return Err(Status(ErrorCode::IoError, "no enough storage node"));
                }

                let mut read_results = Vec::with_capacity(nodes.len());
                for node in nodes {
                    read_results.push(node.get(key));
                }

                let results = join_all(read_results).await;

                let mut read_counts = 0;
                let mut max_ts = 0;
                let mut value = String::new();
                for result in results {
                    if let Ok(Some(result)) = result {
                        if result.ts > max_ts {
                            max_ts = result.ts;
                            value = result.value;
                        }

                        read_counts += 1;
                        if read_counts >= self.read_quorum {
                            break;
                        }
                    } else {
                        continue;
                    }
                }

                if read_counts < self.read_quorum {
                    return Err(Status(
                        ErrorCode::IoError,
                        "no enough storage node to reach read quorum",
                    ));
                }

                return Ok(Some(value));
            }

            None => return Err(Status(ErrorCode::IoError, "no available storage node")),
        }
    }

    fn set(&self, key: String, value: String) -> Result<Option<String>> {
        todo!()
    }

    fn remove(&self, key: &str) -> Result<Option<String>> {
        todo!()
    }
}
