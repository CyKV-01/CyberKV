use futures::future::join_all;
use log::error;

use std::time::Duration;

use tonic::transport::Channel;
use tonic::{Response, Status};

use crate::conn_pool::ConnCache;
use crate::proto::kvs::{ReadRequest, WriteRequest, WriteResponse};
use crate::proto::node::NodeInfo;

use crate::{
    proto::{
        kvs::{key_value_client::*, ReadResponse},
        status::ErrorCode,
    },
    util::*,
};

#[derive(Debug)]
pub struct StorageNode {
    info: NodeInfo,
    client: KeyValueClient<Channel>,
}

impl StorageNode {
    pub async fn new(info: NodeInfo) -> Result<Self, tonic::transport::Error> {
        let channel = Channel::from_shared(format!("http://{}", info.addr))
            .unwrap()
            .timeout(Duration::from_secs(1))
            .connect()
            .await?;

        let client = KeyValueClient::new(channel);
        Ok(Self { info, client })
    }

    pub fn get_client(&self) -> KeyValueClient<Channel> {
        self.client.clone()
    }

    pub async fn get(&mut self, request: ReadRequest) -> Result<Response<ReadResponse>, Status> {
        Ok(self.client.get(request).await?)
    }

    pub async fn set(&mut self, request: WriteRequest) -> Result<Response<WriteResponse>, Status> {
        Ok(self.client.set(request).await?)
    }

    pub async fn remove(
        &mut self,
        _request: WriteRequest,
    ) -> Result<Response<WriteResponse>, Status> {
        todo!()
    }
}

#[derive(Debug)]
pub struct StorageLayer {
    replica_num: usize,
    read_quorum: usize,
    write_quorum: usize,
    conn_cache: ConnCache,
}

impl StorageLayer {
    pub fn new(replica_num: usize, read_quorum: usize, write_quorum: usize) -> Self {
        Self {
            replica_num,
            read_quorum,
            write_quorum,
            conn_cache: ConnCache::new(),
        }
    }

    pub async fn get(&self, request: ReadRequest) -> Result<Response<ReadResponse>, Status> {
        let mut nodes = Vec::with_capacity(request.info.len());
        for info in &request.info {
            match self.conn_cache.get_or_insert(info).await {
                Ok(client) => nodes.push(client),
                Err(err) => {
                    error!(
                        "failed to connect storage node, addr={}, code={}, msg={}",
                        info.addr, err.err_code, err.err_message
                    )
                }
            }
        }

        if nodes.len() < self.read_quorum {
            return Err(Status::failed_precondition("no enough storage node"));
        }

        let mut read_results = Vec::with_capacity(nodes.len());
        for node in nodes.iter_mut() {
            read_results.push(node.get(request.clone()));
        }

        // TODO: wait until reaching read quorum
        let results = join_all(read_results).await;

        let mut read_counts = 0;
        let mut value = String::new();
        let mut ts = 0;
        for result in results {
            if result.is_ok() {
                let read_response = result.unwrap().into_inner();
                if let Some(status) = read_response.status {
                    if status.err_code != ErrorCode::Ok as i32 {
                        error!(
                            "failed to read from storage node, code={}, msg={}",
                            status.err_code, status.err_message
                        );
                    }
                }

                if read_response.ts > ts {
                    ts = read_response.ts;
                    value = read_response.value;
                }

                read_counts += 1;
                if read_counts >= self.read_quorum {
                    break;
                }
            } else {
                error!(
                    "failed to read from storage node, err={}",
                    result.unwrap_err()
                );
                continue;
            }
        }

        if read_counts < self.read_quorum {
            return Err(Status::failed_precondition(
                "no enough storage node to reach read quorum",
            ));
        }

        return Ok(Response::new(ReadResponse {
            ts: ts,
            value: value,
            status: Some(build_status(ErrorCode::Ok, "")),
        }));
    }

    pub async fn set(&self, request: WriteRequest) -> Result<Response<WriteResponse>, Status> {
        // let slot = calc_slot(&request.key);
        let mut nodes = Vec::with_capacity(request.info.len());
        for info in &request.info {
            match self.conn_cache.get_or_insert(info).await {
                Ok(client) => nodes.push(client),
                Err(err) => {
                    error!(
                        "failed to connect storage node, addr={}, code={}, msg={}",
                        info.addr, err.err_code, err.err_message
                    )
                }
            }
        }

        if nodes.len() < self.write_quorum {
            return Err(Status::failed_precondition("no enough storage node"));
        }

        let mut write_results = Vec::with_capacity(nodes.len());
        for node in nodes.iter_mut() {
            write_results.push(node.set(request.clone()));
        }

        // TODO: wait until reaching write quorum
        let results = join_all(write_results).await;

        let mut write_counts = 0;
        for result in results {
            if result.is_ok() {
                write_counts += 1;
                if write_counts >= self.write_quorum {
                    break;
                }
            } else {
                error!(
                    "failed to write storage node, err={}",
                    result.err().unwrap()
                );
                continue;
            }
        }

        if write_counts < self.write_quorum {
            error!(
                "failed to reach write quorum, write_count={}, write_quorum={}",
                write_counts, self.write_quorum
            );
            return Err(Status::failed_precondition(
                "no enough storage node to reach write quorum",
            ));
        }

        return Ok(Response::new(WriteResponse {
            status: Some(build_status(ErrorCode::Ok, "")),
        }));
    }

    pub async fn remove(&self, _request: WriteRequest) -> Result<Response<WriteResponse>, Status> {
        todo!()
    }
}
