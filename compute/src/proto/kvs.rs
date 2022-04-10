#[derive(Clone, PartialEq, ::prost::Message)]
pub struct AssignSlotRequest {
    #[prost(int32, tag = "1")]
    pub slot_id: i32,
}
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct AssignSlotResponse {}
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct ReadOption {
    #[prost(int32, tag = "1")]
    pub consistency_level: i32,
}
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct ReadRequest {
    #[prost(string, tag = "1")]
    pub key: ::prost::alloc::string::String,
    #[prost(uint64, tag = "2")]
    pub ts: u64,
    #[prost(message, optional, tag = "3")]
    pub option: ::core::option::Option<ReadOption>,
    #[prost(message, repeated, tag = "5")]
    pub info: ::prost::alloc::vec::Vec<super::node::NodeInfo>,
}
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct ReadResponse {
    #[prost(string, tag = "1")]
    pub value: ::prost::alloc::string::String,
    #[prost(uint64, tag = "2")]
    pub ts: u64,
    #[prost(message, optional, tag = "3")]
    pub status: ::core::option::Option<super::status::Status>,
}
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct WriteOption {
    #[prost(bool, tag = "1")]
    pub sync: bool,
}
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct WriteRequest {
    #[prost(string, tag = "1")]
    pub key: ::prost::alloc::string::String,
    #[prost(string, tag = "2")]
    pub value: ::prost::alloc::string::String,
    #[prost(bool, tag = "3")]
    pub is_remove: bool,
    #[prost(uint64, tag = "4")]
    pub ts: u64,
    #[prost(message, optional, tag = "5")]
    pub option: ::core::option::Option<WriteOption>,
    #[prost(message, repeated, tag = "7")]
    pub info: ::prost::alloc::vec::Vec<super::node::NodeInfo>,
}
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct WriteResponse {
    #[prost(message, optional, tag = "1")]
    pub status: ::core::option::Option<super::status::Status>,
}
#[doc = r" Generated client implementations."]
pub mod key_value_client {
    #![allow(unused_variables, dead_code, missing_docs, clippy::let_unit_value)]
    use tonic::codegen::*;
    #[derive(Debug, Clone)]
    pub struct KeyValueClient<T> {
        inner: tonic::client::Grpc<T>,
    }
    impl KeyValueClient<tonic::transport::Channel> {
        #[doc = r" Attempt to create a new client by connecting to a given endpoint."]
        pub async fn connect<D>(dst: D) -> Result<Self, tonic::transport::Error>
        where
            D: std::convert::TryInto<tonic::transport::Endpoint>,
            D::Error: Into<StdError>,
        {
            let conn = tonic::transport::Endpoint::new(dst)?.connect().await?;
            Ok(Self::new(conn))
        }
    }
    impl<T> KeyValueClient<T>
    where
        T: tonic::client::GrpcService<tonic::body::BoxBody>,
        T::ResponseBody: Body + Send + 'static,
        T::Error: Into<StdError>,
        <T::ResponseBody as Body>::Error: Into<StdError> + Send,
    {
        pub fn new(inner: T) -> Self {
            let inner = tonic::client::Grpc::new(inner);
            Self { inner }
        }
        pub fn with_interceptor<F>(
            inner: T,
            interceptor: F,
        ) -> KeyValueClient<InterceptedService<T, F>>
        where
            F: tonic::service::Interceptor,
            T: tonic::codegen::Service<
                http::Request<tonic::body::BoxBody>,
                Response = http::Response<
                    <T as tonic::client::GrpcService<tonic::body::BoxBody>>::ResponseBody,
                >,
            >,
            <T as tonic::codegen::Service<http::Request<tonic::body::BoxBody>>>::Error:
                Into<StdError> + Send + Sync,
        {
            KeyValueClient::new(InterceptedService::new(inner, interceptor))
        }
        #[doc = r" Compress requests with `gzip`."]
        #[doc = r""]
        #[doc = r" This requires the server to support it otherwise it might respond with an"]
        #[doc = r" error."]
        pub fn send_gzip(mut self) -> Self {
            self.inner = self.inner.send_gzip();
            self
        }
        #[doc = r" Enable decompressing responses with `gzip`."]
        pub fn accept_gzip(mut self) -> Self {
            self.inner = self.inner.accept_gzip();
            self
        }
        pub async fn assign_slot(
            &mut self,
            request: impl tonic::IntoRequest<super::AssignSlotRequest>,
        ) -> Result<tonic::Response<super::AssignSlotResponse>, tonic::Status> {
            self.inner.ready().await.map_err(|e| {
                tonic::Status::new(
                    tonic::Code::Unknown,
                    format!("Service was not ready: {}", e.into()),
                )
            })?;
            let codec = tonic::codec::ProstCodec::default();
            let path = http::uri::PathAndQuery::from_static("/kvs.KeyValue/AssignSlot");
            self.inner.unary(request.into_request(), path, codec).await
        }
        pub async fn get(
            &mut self,
            request: impl tonic::IntoRequest<super::ReadRequest>,
        ) -> Result<tonic::Response<super::ReadResponse>, tonic::Status> {
            self.inner.ready().await.map_err(|e| {
                tonic::Status::new(
                    tonic::Code::Unknown,
                    format!("Service was not ready: {}", e.into()),
                )
            })?;
            let codec = tonic::codec::ProstCodec::default();
            let path = http::uri::PathAndQuery::from_static("/kvs.KeyValue/Get");
            self.inner.unary(request.into_request(), path, codec).await
        }
        pub async fn set(
            &mut self,
            request: impl tonic::IntoRequest<super::WriteRequest>,
        ) -> Result<tonic::Response<super::WriteResponse>, tonic::Status> {
            self.inner.ready().await.map_err(|e| {
                tonic::Status::new(
                    tonic::Code::Unknown,
                    format!("Service was not ready: {}", e.into()),
                )
            })?;
            let codec = tonic::codec::ProstCodec::default();
            let path = http::uri::PathAndQuery::from_static("/kvs.KeyValue/Set");
            self.inner.unary(request.into_request(), path, codec).await
        }
        pub async fn remove(
            &mut self,
            request: impl tonic::IntoRequest<super::WriteRequest>,
        ) -> Result<tonic::Response<super::WriteResponse>, tonic::Status> {
            self.inner.ready().await.map_err(|e| {
                tonic::Status::new(
                    tonic::Code::Unknown,
                    format!("Service was not ready: {}", e.into()),
                )
            })?;
            let codec = tonic::codec::ProstCodec::default();
            let path = http::uri::PathAndQuery::from_static("/kvs.KeyValue/Remove");
            self.inner.unary(request.into_request(), path, codec).await
        }
    }
}
#[doc = r" Generated server implementations."]
pub mod key_value_server {
    #![allow(unused_variables, dead_code, missing_docs, clippy::let_unit_value)]
    use tonic::codegen::*;
    #[doc = "Generated trait containing gRPC methods that should be implemented for use with KeyValueServer."]
    #[async_trait]
    pub trait KeyValue: Send + Sync + 'static {
        async fn assign_slot(
            &self,
            request: tonic::Request<super::AssignSlotRequest>,
        ) -> Result<tonic::Response<super::AssignSlotResponse>, tonic::Status>;
        async fn get(
            &self,
            request: tonic::Request<super::ReadRequest>,
        ) -> Result<tonic::Response<super::ReadResponse>, tonic::Status>;
        async fn set(
            &self,
            request: tonic::Request<super::WriteRequest>,
        ) -> Result<tonic::Response<super::WriteResponse>, tonic::Status>;
        async fn remove(
            &self,
            request: tonic::Request<super::WriteRequest>,
        ) -> Result<tonic::Response<super::WriteResponse>, tonic::Status>;
    }
    #[derive(Debug)]
    pub struct KeyValueServer<T: KeyValue> {
        inner: _Inner<T>,
        accept_compression_encodings: (),
        send_compression_encodings: (),
    }
    struct _Inner<T>(Arc<T>);
    impl<T: KeyValue> KeyValueServer<T> {
        pub fn new(inner: T) -> Self {
            let inner = Arc::new(inner);
            let inner = _Inner(inner);
            Self {
                inner,
                accept_compression_encodings: Default::default(),
                send_compression_encodings: Default::default(),
            }
        }
        pub fn with_interceptor<F>(inner: T, interceptor: F) -> InterceptedService<Self, F>
        where
            F: tonic::service::Interceptor,
        {
            InterceptedService::new(Self::new(inner), interceptor)
        }
    }
    impl<T, B> tonic::codegen::Service<http::Request<B>> for KeyValueServer<T>
    where
        T: KeyValue,
        B: Body + Send + 'static,
        B::Error: Into<StdError> + Send + 'static,
    {
        type Response = http::Response<tonic::body::BoxBody>;
        type Error = Never;
        type Future = BoxFuture<Self::Response, Self::Error>;
        fn poll_ready(&mut self, _cx: &mut Context<'_>) -> Poll<Result<(), Self::Error>> {
            Poll::Ready(Ok(()))
        }
        fn call(&mut self, req: http::Request<B>) -> Self::Future {
            let inner = self.inner.clone();
            match req.uri().path() {
                "/kvs.KeyValue/AssignSlot" => {
                    #[allow(non_camel_case_types)]
                    struct AssignSlotSvc<T: KeyValue>(pub Arc<T>);
                    impl<T: KeyValue> tonic::server::UnaryService<super::AssignSlotRequest> for AssignSlotSvc<T> {
                        type Response = super::AssignSlotResponse;
                        type Future = BoxFuture<tonic::Response<Self::Response>, tonic::Status>;
                        fn call(
                            &mut self,
                            request: tonic::Request<super::AssignSlotRequest>,
                        ) -> Self::Future {
                            let inner = self.0.clone();
                            let fut = async move { (*inner).assign_slot(request).await };
                            Box::pin(fut)
                        }
                    }
                    let accept_compression_encodings = self.accept_compression_encodings;
                    let send_compression_encodings = self.send_compression_encodings;
                    let inner = self.inner.clone();
                    let fut = async move {
                        let inner = inner.0;
                        let method = AssignSlotSvc(inner);
                        let codec = tonic::codec::ProstCodec::default();
                        let mut grpc = tonic::server::Grpc::new(codec).apply_compression_config(
                            accept_compression_encodings,
                            send_compression_encodings,
                        );
                        let res = grpc.unary(method, req).await;
                        Ok(res)
                    };
                    Box::pin(fut)
                }
                "/kvs.KeyValue/Get" => {
                    #[allow(non_camel_case_types)]
                    struct GetSvc<T: KeyValue>(pub Arc<T>);
                    impl<T: KeyValue> tonic::server::UnaryService<super::ReadRequest> for GetSvc<T> {
                        type Response = super::ReadResponse;
                        type Future = BoxFuture<tonic::Response<Self::Response>, tonic::Status>;
                        fn call(
                            &mut self,
                            request: tonic::Request<super::ReadRequest>,
                        ) -> Self::Future {
                            let inner = self.0.clone();
                            let fut = async move { (*inner).get(request).await };
                            Box::pin(fut)
                        }
                    }
                    let accept_compression_encodings = self.accept_compression_encodings;
                    let send_compression_encodings = self.send_compression_encodings;
                    let inner = self.inner.clone();
                    let fut = async move {
                        let inner = inner.0;
                        let method = GetSvc(inner);
                        let codec = tonic::codec::ProstCodec::default();
                        let mut grpc = tonic::server::Grpc::new(codec).apply_compression_config(
                            accept_compression_encodings,
                            send_compression_encodings,
                        );
                        let res = grpc.unary(method, req).await;
                        Ok(res)
                    };
                    Box::pin(fut)
                }
                "/kvs.KeyValue/Set" => {
                    #[allow(non_camel_case_types)]
                    struct SetSvc<T: KeyValue>(pub Arc<T>);
                    impl<T: KeyValue> tonic::server::UnaryService<super::WriteRequest> for SetSvc<T> {
                        type Response = super::WriteResponse;
                        type Future = BoxFuture<tonic::Response<Self::Response>, tonic::Status>;
                        fn call(
                            &mut self,
                            request: tonic::Request<super::WriteRequest>,
                        ) -> Self::Future {
                            let inner = self.0.clone();
                            let fut = async move { (*inner).set(request).await };
                            Box::pin(fut)
                        }
                    }
                    let accept_compression_encodings = self.accept_compression_encodings;
                    let send_compression_encodings = self.send_compression_encodings;
                    let inner = self.inner.clone();
                    let fut = async move {
                        let inner = inner.0;
                        let method = SetSvc(inner);
                        let codec = tonic::codec::ProstCodec::default();
                        let mut grpc = tonic::server::Grpc::new(codec).apply_compression_config(
                            accept_compression_encodings,
                            send_compression_encodings,
                        );
                        let res = grpc.unary(method, req).await;
                        Ok(res)
                    };
                    Box::pin(fut)
                }
                "/kvs.KeyValue/Remove" => {
                    #[allow(non_camel_case_types)]
                    struct RemoveSvc<T: KeyValue>(pub Arc<T>);
                    impl<T: KeyValue> tonic::server::UnaryService<super::WriteRequest> for RemoveSvc<T> {
                        type Response = super::WriteResponse;
                        type Future = BoxFuture<tonic::Response<Self::Response>, tonic::Status>;
                        fn call(
                            &mut self,
                            request: tonic::Request<super::WriteRequest>,
                        ) -> Self::Future {
                            let inner = self.0.clone();
                            let fut = async move { (*inner).remove(request).await };
                            Box::pin(fut)
                        }
                    }
                    let accept_compression_encodings = self.accept_compression_encodings;
                    let send_compression_encodings = self.send_compression_encodings;
                    let inner = self.inner.clone();
                    let fut = async move {
                        let inner = inner.0;
                        let method = RemoveSvc(inner);
                        let codec = tonic::codec::ProstCodec::default();
                        let mut grpc = tonic::server::Grpc::new(codec).apply_compression_config(
                            accept_compression_encodings,
                            send_compression_encodings,
                        );
                        let res = grpc.unary(method, req).await;
                        Ok(res)
                    };
                    Box::pin(fut)
                }
                _ => Box::pin(async move {
                    Ok(http::Response::builder()
                        .status(200)
                        .header("grpc-status", "12")
                        .header("content-type", "application/grpc")
                        .body(empty_body())
                        .unwrap())
                }),
            }
        }
    }
    impl<T: KeyValue> Clone for KeyValueServer<T> {
        fn clone(&self) -> Self {
            let inner = self.inner.clone();
            Self {
                inner,
                accept_compression_encodings: self.accept_compression_encodings,
                send_compression_encodings: self.send_compression_encodings,
            }
        }
    }
    impl<T: KeyValue> Clone for _Inner<T> {
        fn clone(&self) -> Self {
            Self(self.0.clone())
        }
    }
    impl<T: std::fmt::Debug> std::fmt::Debug for _Inner<T> {
        fn fmt(&self, f: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
            write!(f, "{:?}", self.0)
        }
    }
    impl<T: KeyValue> tonic::transport::NamedService for KeyValueServer<T> {
        const NAME: &'static str = "kvs.KeyValue";
    }
}
