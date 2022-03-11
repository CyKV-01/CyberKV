use crate::{error::Result, types::Value};

pub trait Kvs: Sync + Send + 'static {
    fn get(&self, key: &str) -> Result<Option<&Value>>;
    fn set(&self, key: String, value: Value) -> Result<Option<Value>>;
    fn remove(&self, key: &str) -> Result<Option<Value>>;
}
