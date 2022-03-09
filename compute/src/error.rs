use crate::proto::status;

pub type Result<T> = std::result::Result<T, status::Status>;
