use crate::error::Result;
use crate::kvs::Kvs;
use std::collections::hash_map::DefaultHasher;
use std::hash::{Hash, Hasher};
use std::{collections::BTreeMap, sync::Arc, sync::RwLock};

type Bucket<K, V> = Arc<RwLock<BTreeMap<K, V>>>;
#[derive(Debug)]
pub struct ConcurrentMap<K, V> {
    buckets: Arc<Vec<Bucket<K, V>>>,
}

impl Kvs for ConcurrentMap<String, String> {
    fn get(&self, key: &str) -> Result<Option<String>> {
        let slot = hash(key) as usize % self.buckets.len();

        let bucket = self.buckets[slot].read().unwrap();

        match bucket.get(key) {
            Some(value) => Ok(Some(value.to_owned())),
            None => Ok(None),
        }
    }

    fn set(&self, key: String, value: String) -> Result<Option<String>> {
        let slot = hash(&key) as usize % self.buckets.len();

        let mut bucket = self.buckets[slot].write().unwrap();

        Ok(bucket.insert(key, value))
    }

    fn remove(&self, key: &str) -> Result<Option<String>> {
        let slot = hash(key) as usize % self.buckets.len();

        let mut bucket = self.buckets[slot].write().unwrap();

        Ok(bucket.remove(key))
    }
}

fn hash(key: &str) -> u64 {
    let mut s = DefaultHasher::new();
    key.hash(&mut s);

    s.finish()
}
