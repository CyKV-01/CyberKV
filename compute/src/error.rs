use std::net::AddrParseError;

use crate::proto::status::{self, ErrorCode};
use crate::types::build_status;

pub type Result<T> = std::result::Result<T, status::Status>;

impl From<etcd_client::Error> for status::Status {
    fn from(err: etcd_client::Error) -> Self {
        build_status(
            ErrorCode::IoError,
            format!("operation on etcd failed, err={}", err).as_str(),
        )
    }
}

impl From<AddrParseError> for status::Status {
    fn from(err: AddrParseError) -> Self {
        build_status(
            ErrorCode::InvalidArgument,
            format!("failed to parse addr, err={}", err).as_str(),
        )
    }
}

impl From<tonic::transport::Error> for status::Status {
    fn from(err: tonic::transport::Error) -> Self {
        build_status(
            ErrorCode::IoError,
            format!("rpc failed, err={}", err).as_str(),
        )
    }
}
