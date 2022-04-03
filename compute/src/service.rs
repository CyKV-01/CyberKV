use dashmap::DashMap;
use log::error;
use tonic::{Request, Response, Status};

use crate::compute_node::ComputeNode;
use crate::concurrent_map::MemTable;
use crate::proto::kvs::key_value_server::KeyValue;
use crate::proto::kvs::{ReadRequest, ReadResponse, WriteRequest, WriteResponse};
use crate::proto::status::{self, ErrorCode};
use crate::storage::StorageLayer;
use crate::types::Value;
use crate::{consts, util::*};

pub struct KvServer {
    mem_table: DashMap<String, Value>,
    storage_layer: StorageLayer,
}

impl KvServer {
    pub fn new() -> Self {
        Self {
            mem_table: DashMap::new(),
            storage_layer: StorageLayer::new(
                consts::DEFAULT_REPLICA_NUM,
                consts::DEFAULT_READ_QUORUM,
                consts::DEFAULT_WRITE_QUORUM,
            ),
        }
    }
}

#[tonic::async_trait]
impl KeyValue for KvServer {
    async fn get(&self, request: Request<ReadRequest>) -> Result<Response<ReadResponse>, Status> {
        let request = request.into_inner();
        let key = request.key.clone();

        let response = match self.mem_table.get(&request.key) {
            Some(value) => ReadResponse {
                value: value.value.clone(),
                ts: value.timestamp,
                status: None,
            },
            None => {
                let response = self.storage_layer.get(request).await?.into_inner();
                let new_value = Value {
                    timestamp: response.ts,
                    value: response.value.clone(),
                };

                self.mem_table
                    .entry(key)
                    .and_modify(|value| {
                        if value.timestamp < response.ts {
                            value.timestamp = new_value.timestamp;
                            value.value = new_value.value.clone();
                        }
                    })
                    .or_insert(new_value);

                response
            }
        };

        Ok(Response::new(response))
    }

    async fn set(&self, request: Request<WriteRequest>) -> Result<Response<WriteResponse>, Status> {
        let request = request.into_inner();
        let key = request.key.clone();

        // write WAL
        let response = self.storage_layer.set(request).await?.into_inner();
        if let Some(status) = &response.status {
            if status.err_code != ErrorCode::Ok as i32 {
                error!(
                    "failed to write storage layer, code={}, msg={}",
                    status.err_code, status.err_message
                );
                return Ok(Response::new(response));
            }
        }

        // write memtable
        self.mem_table.remove(&key);

        Ok(Response::new(WriteResponse {
            status: Some(build_status(ErrorCode::Ok, "")),
        }))
    }

    async fn remove(
        &self,
        request: Request<WriteRequest>,
    ) -> Result<Response<WriteResponse>, Status> {
        let request = request.into_inner();

        self.mem_table.remove(&request.key);

        Ok(Response::new(WriteResponse {
            status: Some(build_status(ErrorCode::Ok, "")),
        }))
    }
}
