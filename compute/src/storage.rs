use futures::future::join_all;
use std::cell::RefCell;
use std::collections::HashMap;
use std::sync::Arc;
use tokio::sync::RwLock;
use tonic::transport::Channel;
use tonic::{Response, Status};

use crate::proto::kvs::{ReadRequest, WriteRequest, WriteResponse};
use crate::proto::node::NodeInfo;
use crate::types::Value;
use crate::{
    proto::{
        kvs::{key_value_client::*, ReadResponse},
        status::ErrorCode,
    },
    types::SlotID,
    util::*,
};
#[derive(Debug)]
pub struct StorageNode {
    info: NodeInfo,
    client: KeyValueClient<Channel>,
}

impl StorageNode {
    pub async fn new(info: NodeInfo) -> Result<Self, tonic::transport::Error> {
        let client = KeyValueClient::connect(format!("http://{}", info.addr)).await?;
        Ok(Self { info, client })
    }

    pub async fn get(&mut self, request: ReadRequest) -> Result<Response<ReadResponse>, Status> {
        todo!()
    }

    pub async fn set(&mut self, request: WriteRequest) -> Result<Response<WriteResponse>, Status> {
        Ok(self.client.set(request).await?)
    }

    pub async fn remove(
        &mut self,
        request: WriteRequest,
    ) -> Result<Response<WriteResponse>, Status> {
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

    pub async fn get(&self, request: ReadRequest) -> Result<Response<ReadResponse>, Status> {
        // let slot = calc_slot(&request.key);
        // let slot_node_index = self.slot_node_index.read().await;
        // let nodes = slot_node_index.get(&slot);

        // if nodes.is_none() {
        //     return Err(Status::failed_precondition("no storage node to serve"));
        // }

        // let nodes = nodes.unwrap();

        // if nodes.len() < self.read_quorum {
        //     return Err(Status::failed_precondition("no enough storage node"));
        // }

        // let mut read_results = Vec::with_capacity(nodes.len());
        // for node in nodes {
        //     read_results.push(node.get(request.clone()));
        // }

        // // TODO: wait until reaching read quorum
        // let results = join_all(read_results).await;

        // let mut read_counts = 0;
        // let mut value = Value {
        //     timestamp: 0,
        //     value: String::new(),
        // };
        // for result in results {
        //     if let Ok(result) = result {
        //         let respones = result.into_inner();
        //         if respones.ts > value.timestamp {
        //             value.timestamp = respones.ts;
        //             value.value = respones.value;
        //         }

        //         read_counts += 1;
        //         if read_counts >= self.read_quorum {
        //             break;
        //         }
        //     } else {
        //         continue;
        //     }
        // }

        // if read_counts < self.read_quorum {
        //     return Err(Status::failed_precondition(
        //         "no enough storage node to reach read quorum",
        //     ));
        // }

        // return Ok(Response::new(ReadResponse {
        //     value: value.value,
        //     ts: value.timestamp,
        //     status: None,
        // }));
        todo!()
    }

    pub async fn set(&self, request: WriteRequest) -> Result<Response<WriteResponse>, Status> {
        // let slot = calc_slot(&request.key);
        let mut nodes = Vec::with_capacity(request.info.len());
        for info in &request.info {
            match StorageNode::new(info.clone()).await {
                Ok(node) => nodes.push(node),
                Err(err) => {
                    return Err(Status::failed_precondition(format!(
                        "failed to connect storage node, err={}",
                        err
                    )))
                }
            }
        }

        // let mut slot_node_index = self.slot_node_index.write().await;
        // slot_node_index.insert(slot, nodes);
        // let slot_node_index = self.slot_node_index.read().await;
        // let nodes = slot_node_index.get_mut(&slot);

        // if nodes.is_none() {
        //     return Err(Status::failed_precondition("no available storage node"));
        // }

        // let nodes = nodes.unwrap();
        if nodes.len() < self.write_quorum {
            return Err(Status::failed_precondition("no enough storage node"));
        }

        let mut write_results = Vec::with_capacity(nodes.len());
        for node in nodes.iter_mut() {
            write_results.push(node.set(request.clone()));
        }

        // TODO: wait until reaching read quorum
        let results = join_all(write_results).await;

        let mut write_counts = 0;
        for result in results {
            if result.is_ok() {
                write_counts += 1;
                if write_counts >= self.write_quorum {
                    break;
                }
            } else {
                continue;
            }
        }

        if write_counts < self.write_quorum {
            return Err(Status::failed_precondition(
                "no enough storage node to reach write quorum",
            ));
        }

        return Ok(Response::new(WriteResponse {
            status: Some(build_status(ErrorCode::Ok, "")),
        }));
    }

    pub async fn remove(&self, request: WriteRequest) -> Result<Response<WriteResponse>, Status> {
        todo!()
    }
}
