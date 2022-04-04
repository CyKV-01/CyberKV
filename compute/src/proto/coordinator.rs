#[derive(Clone, PartialEq, ::prost::Message)]
pub struct AllocateSsTableRequest {}
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct AllocateSsTableResponse {
    #[prost(uint64, tag = "1")]
    pub id: u64,
}
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct ReportStatsRequest {
    #[prost(uint64, tag = "1")]
    pub id: u64,
    /// slot_id -> mem_table_size
    #[prost(map = "int32, uint64", tag = "2")]
    pub mem_table_size: ::std::collections::HashMap<i32, u64>,
    #[prost(uint64, tag = "3")]
    pub memory_usage: u64,
    #[prost(uint64, tag = "4")]
    pub memory_cap: u64,
    /// percent, full usage of dual cores is 200
    #[prost(int32, tag = "5")]
    pub cpu_usage: i32,
}
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct ReportStatsResponse {}
#[doc = r" Generated client implementations."]
pub mod coordinator_client {
    #![allow(unused_variables, dead_code, missing_docs, clippy::let_unit_value)]
    use tonic::codegen::*;
    #[derive(Debug, Clone)]
    pub struct CoordinatorClient<T> {
        inner: tonic::client::Grpc<T>,
    }
    impl CoordinatorClient<tonic::transport::Channel> {
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
    impl<T> CoordinatorClient<T>
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
        ) -> CoordinatorClient<InterceptedService<T, F>>
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
            CoordinatorClient::new(InterceptedService::new(inner, interceptor))
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
        pub async fn allocate_ss_table_id(
            &mut self,
            request: impl tonic::IntoRequest<super::AllocateSsTableRequest>,
        ) -> Result<tonic::Response<super::AllocateSsTableResponse>, tonic::Status> {
            self.inner.ready().await.map_err(|e| {
                tonic::Status::new(
                    tonic::Code::Unknown,
                    format!("Service was not ready: {}", e.into()),
                )
            })?;
            let codec = tonic::codec::ProstCodec::default();
            let path =
                http::uri::PathAndQuery::from_static("/coordinator.Coordinator/AllocateSSTableID");
            self.inner.unary(request.into_request(), path, codec).await
        }
        pub async fn report_stats(
            &mut self,
            request: impl tonic::IntoRequest<super::ReportStatsRequest>,
        ) -> Result<tonic::Response<super::ReportStatsResponse>, tonic::Status> {
            self.inner.ready().await.map_err(|e| {
                tonic::Status::new(
                    tonic::Code::Unknown,
                    format!("Service was not ready: {}", e.into()),
                )
            })?;
            let codec = tonic::codec::ProstCodec::default();
            let path = http::uri::PathAndQuery::from_static("/coordinator.Coordinator/ReportStats");
            self.inner.unary(request.into_request(), path, codec).await
        }
    }
}
#[doc = r" Generated server implementations."]
pub mod coordinator_server {
    #![allow(unused_variables, dead_code, missing_docs, clippy::let_unit_value)]
    use tonic::codegen::*;
    #[doc = "Generated trait containing gRPC methods that should be implemented for use with CoordinatorServer."]
    #[async_trait]
    pub trait Coordinator: Send + Sync + 'static {
        async fn allocate_ss_table_id(
            &self,
            request: tonic::Request<super::AllocateSsTableRequest>,
        ) -> Result<tonic::Response<super::AllocateSsTableResponse>, tonic::Status>;
        async fn report_stats(
            &self,
            request: tonic::Request<super::ReportStatsRequest>,
        ) -> Result<tonic::Response<super::ReportStatsResponse>, tonic::Status>;
    }
    #[derive(Debug)]
    pub struct CoordinatorServer<T: Coordinator> {
        inner: _Inner<T>,
        accept_compression_encodings: (),
        send_compression_encodings: (),
    }
    struct _Inner<T>(Arc<T>);
    impl<T: Coordinator> CoordinatorServer<T> {
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
    impl<T, B> tonic::codegen::Service<http::Request<B>> for CoordinatorServer<T>
    where
        T: Coordinator,
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
                "/coordinator.Coordinator/AllocateSSTableID" => {
                    #[allow(non_camel_case_types)]
                    struct AllocateSSTableIDSvc<T: Coordinator>(pub Arc<T>);
                    impl<T: Coordinator> tonic::server::UnaryService<super::AllocateSsTableRequest>
                        for AllocateSSTableIDSvc<T>
                    {
                        type Response = super::AllocateSsTableResponse;
                        type Future = BoxFuture<tonic::Response<Self::Response>, tonic::Status>;
                        fn call(
                            &mut self,
                            request: tonic::Request<super::AllocateSsTableRequest>,
                        ) -> Self::Future {
                            let inner = self.0.clone();
                            let fut = async move { (*inner).allocate_ss_table_id(request).await };
                            Box::pin(fut)
                        }
                    }
                    let accept_compression_encodings = self.accept_compression_encodings;
                    let send_compression_encodings = self.send_compression_encodings;
                    let inner = self.inner.clone();
                    let fut = async move {
                        let inner = inner.0;
                        let method = AllocateSSTableIDSvc(inner);
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
                "/coordinator.Coordinator/ReportStats" => {
                    #[allow(non_camel_case_types)]
                    struct ReportStatsSvc<T: Coordinator>(pub Arc<T>);
                    impl<T: Coordinator> tonic::server::UnaryService<super::ReportStatsRequest> for ReportStatsSvc<T> {
                        type Response = super::ReportStatsResponse;
                        type Future = BoxFuture<tonic::Response<Self::Response>, tonic::Status>;
                        fn call(
                            &mut self,
                            request: tonic::Request<super::ReportStatsRequest>,
                        ) -> Self::Future {
                            let inner = self.0.clone();
                            let fut = async move { (*inner).report_stats(request).await };
                            Box::pin(fut)
                        }
                    }
                    let accept_compression_encodings = self.accept_compression_encodings;
                    let send_compression_encodings = self.send_compression_encodings;
                    let inner = self.inner.clone();
                    let fut = async move {
                        let inner = inner.0;
                        let method = ReportStatsSvc(inner);
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
    impl<T: Coordinator> Clone for CoordinatorServer<T> {
        fn clone(&self) -> Self {
            let inner = self.inner.clone();
            Self {
                inner,
                accept_compression_encodings: self.accept_compression_encodings,
                send_compression_encodings: self.send_compression_encodings,
            }
        }
    }
    impl<T: Coordinator> Clone for _Inner<T> {
        fn clone(&self) -> Self {
            Self(self.0.clone())
        }
    }
    impl<T: std::fmt::Debug> std::fmt::Debug for _Inner<T> {
        fn fmt(&self, f: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
            write!(f, "{:?}", self.0)
        }
    }
    impl<T: Coordinator> tonic::transport::NamedService for CoordinatorServer<T> {
        const NAME: &'static str = "coordinator.Coordinator";
    }
}
