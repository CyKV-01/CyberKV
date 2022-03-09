use crate::concurrent_map::ConcurrentMap;
use crate::kvs::Kvs;
use crate::proto::kv::key_value_client::*;
use crate::proto::kv::{key_value_server::*, *};
use crate::proto::status::{self, ErrorCode};
use crate::storage::StorageLayer;
use tonic::{transport::Server, Request, Response, Status};

#[derive(Debug)]
pub struct ComputeNode<MemTable: Kvs> {
    kvs: MemTable,
    
    storage_layer: StorageLayer,
}

#[tonic::async_trait]
impl<MemTable: Kvs> KeyValue for ComputeNode<MemTable> {
    async fn get(&self, request: Request<ReadRequest>) -> Result<Response<ReadResponse>, Status> {
        let request = request.into_inner();
        
        let response = match self.kvs.get(&request.key) {
            Ok(value) => {
                if let Some(value) = value {
                    ReadResponse {
                        value,
                        status: None,
                    }
                } else {
                    ReadResponse {
                        value: String::new(),
                        status: Some(status::Status {
                            err_code: ErrorCode::KeyNotFound.into(),
                            err_message: "key not found".to_string(),
                        }),
                    }
                }
            }

            Err(err) => ReadResponse {
                value: String::new(),
                status: Some(err),
            },
        };

        Ok(Response::new(response))
    }

    async fn set(&self, request: Request<WriteRequest>) -> Result<Response<WriteResponse>, Status> {
        let request = request.into_inner();
    }

    async fn remove(
        &self,
        request: Request<WriteRequest>,
    ) -> Result<Response<WriteResponse>, Status> {
        todo!()
    }
}
