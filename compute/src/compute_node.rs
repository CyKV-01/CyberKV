use crate::concurrent_map::{ConcurrentMap, MemTable};
use crate::kvs::Kvs;
use crate::proto::kv::key_value_client::*;
use crate::proto::kv::{key_value_server::*, *};
use crate::proto::status::{self, ErrorCode};
use crate::storage::StorageLayer;
use crate::types::build_status;
use crate::types::Value;
use tonic::{transport::Server, Request, Response, Status};

#[derive(Debug)]
pub struct ComputeNode {
    mem_table: MemTable,

    storage_layer: StorageLayer,
}

#[tonic::async_trait]
impl KeyValue for ComputeNode {
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

        let response = match self
            .mem_table
            .set(
                request.key,
                Value {
                    timestamp: request.ts,
                    value: request.value,
                },
            )
            .await
        {
            Ok(_) => WriteResponse {
                status: Some(build_status(ErrorCode::Ok, "")),
            },
            Err(err) => WriteResponse { status: Some(err) },
        };

        Ok(Response::new(response))
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
