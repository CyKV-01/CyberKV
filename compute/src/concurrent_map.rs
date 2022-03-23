use crate::error::Result;
use crate::kvs::Kvs;
use crate::types::Value;
use std::collections::hash_map::DefaultHasher;
use std::hash::{Hash, Hasher};
use std::{collections::BTreeMap, sync::Arc};
use tokio::sync::RwLock;

type Bucket<K, V> = Arc<RwLock<BTreeMap<K, V>>>;
pub type MemTable = ConcurrentMap<String, Value>;

fn new_bucket() -> Bucket<String, Value> {
    Arc::new(RwLock::new(BTreeMap::new()))
}

#[derive(Debug)]
pub struct ConcurrentMap<K, V> {
    buckets: Arc<Vec<Bucket<K, V>>>,
}

impl ConcurrentMap<String, Value> {
    pub fn new() -> Self {
        let mut buckets = Vec::with_capacity(12);
        for _ in 0..buckets.capacity() {
            buckets.push(new_bucket());
        }
        Self {
            buckets: Arc::new(buckets),
        }
    }

    pub async fn get(&self, key: &str) -> Option<Value> {
        let slot = hash(key) as usize % self.buckets.len();

        let bucket = self.buckets[slot].read().await;

        match bucket.get(key) {
            Some(value) => Some(value.clone()),
            None => None,
        }
    }

    pub async fn set(&self, key: String, value: Value) -> Result<Option<Value>> {
        let slot = hash(&key) as usize % self.buckets.len();

        let mut bucket = self.buckets[slot].write().await;

        Ok(bucket.insert(key, value))
    }

    pub async fn remove(&self, key: &str) -> Result<Option<Value>> {
        let slot = hash(key) as usize % self.buckets.len();

        let mut bucket = self.buckets[slot].write().await;

        Ok(bucket.remove(key))
    }
}

fn hash(key: &str) -> u64 {
    let mut s = DefaultHasher::new();
    key.hash(&mut s);

    s.finish()
}
