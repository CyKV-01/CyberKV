use crate::error::Result;
use crate::proto::kvs::key_value_client::KeyValueClient;
use crate::proto::node::NodeInfo;
use crate::storage::{StorageLayer, StorageNode};
use dashmap::DashMap;
use tonic::transport::Channel;

#[derive(Debug)]
pub struct ConnCache {
    cache: DashMap<String, StorageNode>,
}

impl ConnCache {
    pub fn new() -> Self {
        Self {
            cache: DashMap::new(),
        }
    }

    pub fn get(&self, node_id: &str) -> Option<KeyValueClient<Channel>> {
        match self.cache.get(node_id) {
            Some(entry) => Some(entry.value().get_client()),
            None => None,
        }
    }

    pub async fn get_or_insert(&self, info: &NodeInfo) -> Result<KeyValueClient<Channel>> {
        match self.get(&info.id) {
            Some(client) => Ok(client),
            None => self.insert(info.clone()).await,
        }
    }

    pub async fn insert(&self, info: NodeInfo) -> Result<KeyValueClient<Channel>> {
        let id = info.id.clone();
        let node = StorageNode::new(info).await?;
        let client = node.get_client();
        self.cache.insert(id, node);
        Ok(client)
    }
}
