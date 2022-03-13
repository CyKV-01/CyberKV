// enum NodeType {
//     Invalid = 0;
//     ComputeNode = 1;
//     StorageNode = 2;
// }

#[derive(Clone, PartialEq, ::prost::Message)]
pub struct NodeInfo {
    #[prost(string, tag = "1")]
    pub id: ::prost::alloc::string::String,
    /// NodeType type = 2;
    #[prost(string, tag = "2")]
    pub addr: ::prost::alloc::string::String,
}
