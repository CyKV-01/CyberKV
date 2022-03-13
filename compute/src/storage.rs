use std::collections::HashMap;

use futures::future::join_all;
use tokio::sync::RwLock;
use tonic::transport::Channel;

use crate::proto::node::NodeInfo;
use crate::types::Value;
use crate::{
    error::Result,
    kvs::Kvs,
    proto::{
        kvs::{key_value_client::*, ReadResponse},
        status::ErrorCode,
    },
    types::{build_status, calc_slot, SlotID},
};
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
    read_quorum: usize,
    write_quorum: usize,
    slot_node_index: RwLock<HashMap<SlotID, Vec<StorageNode>>>,
}

impl StorageLayer {
    pub fn new(replica_num: usize, read_quorum: usize, write_quorum: usize) -> Self {
        Self {
            replica_num,
            read_quorum,
            write_quorum,
            slot_node_index: RwLock::new(HashMap::new()),
        }
    }

    pub async fn get(&self, key: &str) -> Result<Option<Value>> {
        let slot = calc_slot(key);
        let slot_node_index = self.slot_node_index.read().await;
        let nodes = slot_node_index.get(&slot);

        match nodes {
            Some(nodes) => {
                if nodes.len() < self.read_quorum {
                    return Err(build_status(ErrorCode::IoError, "no enough storage node"));
                }

                let mut read_results = Vec::with_capacity(nodes.len());
                for node in nodes {
                    read_results.push(node.get(key));
                }

                // TODO: wait until reaching read quorum
                let results = join_all(read_results).await;

                let mut read_counts = 0;
                let mut value = Value {
                    timestamp: 0,
                    value: String::new(),
                };
                for result in results {
                    if let Ok(Some(result)) = result {
                        if result.ts > value.timestamp {
                            value.timestamp = result.ts;
                            value.value = result.value;
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
                    return Err(build_status(
                        ErrorCode::IoError,
                        "no enough storage node to reach read quorum",
                    ));
                }

                return Ok(Some(value));
            }

            None => {
                return Err(build_status(
                    ErrorCode::IoError,
                    "no available storage node",
                ))
            }
        }
    }

    pub async fn set(&self, key: String, value: Value) -> Result<Option<Value>> {
        todo!()
    }

    pub async fn remove(&self, key: &str) -> Result<Option<Value>> {
        todo!()
    }
}
