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
            Some(value) => ReadResponse {
                value: value.value,
                ts: value.timestamp,
                status: None,
            },
            None => self.storage_layer.get(request).await?.into_inner(),
        };

        Ok(Response::new(response))
    }

    async fn set(&self, request: Request<WriteRequest>) -> Result<Response<WriteResponse>, Status> {
        let request = request.into_inner();
        let key = request.key.clone();
        let value = Value::new(request.ts, request.value.clone());

        // write WAL
        let response = self.storage_layer.set(request).await?.into_inner();
        if let Some(status) = &response.status {
            if status.err_code != ErrorCode::Ok as i32 {
                return Ok(Response::new(response));
            }
        }

        // write memtable
        let status = match self.mem_table.set(key, value).await {
            Ok(_) => build_status(ErrorCode::Ok, ""),
            Err(err) => err,
        };

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
