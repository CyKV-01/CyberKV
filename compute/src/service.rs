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
    mem_table: MemTable,
    storage_layer: StorageLayer,
}

impl KvServer {
    pub fn new() -> Self {
        Self {
            mem_table: MemTable::new(),
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

        let response = match self.mem_table.get(&request.key).await {
            Ok(value) => {
                if let Some(value) = value {
                    ReadResponse {
                        value: value.value,
                        ts: value.timestamp,
                        status: None,
                    }
                } else {
                    // TODO: Cache miss in MemTable, read it from storage layer
                    ReadResponse {
                        value: String::new(),
                        ts: 0,
                        status: Some(status::Status {
                            err_code: ErrorCode::KeyNotFound.into(),
                            err_message: "key not found".to_string(),
                        }),
                    }
                }
            }

            Err(err) => ReadResponse {
                value: String::new(),
                ts: 0,
                status: Some(err),
            },
        };

        Ok(Response::new(response))
    }

    async fn set(&self, request: Request<WriteRequest>) -> Result<Response<WriteResponse>, Status> {
        let request = request.into_inner();
        let value = Value::new(request.ts, request.value.clone());

        // write memtable
        let status = match self.mem_table.set(request.key.clone(), value).await {
            Ok(_) => build_status(ErrorCode::Ok, ""),
            Err(err) => err,
        };

        // let mut is_sync = false;
        // if let Some(option) = &request.option {
        //     is_sync = option.sync;
        // }

        self.storage_layer.set(request).await?;

        Ok(Response::new(WriteResponse {
            status: Some(status),
        }))
    }

    async fn remove(
        &self,
        request: Request<WriteRequest>,
    ) -> Result<Response<WriteResponse>, Status> {
        let request = request.into_inner();

        let response = match self.mem_table.remove(&request.key).await {
            Ok(_) => WriteResponse {
                status: Some(build_status(ErrorCode::Ok, "")),
            },
            Err(err) => WriteResponse { status: Some(err) },
        };

        Ok(Response::new(response))
    }
}
