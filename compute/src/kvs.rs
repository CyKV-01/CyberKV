use crate::error::Result;

pub trait Kvs: Sync + Send + 'static {
    fn get(&self, key: &str) -> Result<Option<String>>;
    fn set(&self, key: String, value: String) -> Result<Option<String>>;
    fn remove(&self, key: &str) -> Result<Option<String>>;
}
