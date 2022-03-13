#[derive(Clone, PartialEq, ::prost::Message)]
pub struct Status {
    #[prost(enumeration = "ErrorCode", tag = "1")]
    pub err_code: i32,
    #[prost(string, tag = "2")]
    pub err_message: ::prost::alloc::string::String,
}
#[derive(Clone, Copy, Debug, PartialEq, Eq, Hash, PartialOrd, Ord, ::prost::Enumeration)]
#[repr(i32)]
pub enum ErrorCode {
    Ok = 0,
    KeyNotFound = 1,
    IoError = 2,
    Corruption = 3,
    InvalidArgument = 4,
    Timeout = 5,
    RetryLater = 6,
}
