// package proxy

// import (
// 	"context"

// 	"github.com/yah01/CyberKV/proto"
// )

// type Proxy struct {
// 	proto.UnsafeKeyValueServer
// }

// func (proxy *Proxy) Get(context.Context, *proto.ReadRequest) (*proto.ReadResponse, error) {
// 	return nil, status.Errorf(codes.Unimplemented, "method Get not implemented")
// }
// func (proxy *Proxy) Set(context.Context, *proto.WriteRequest) (*proto.WriteResponse, error) {
// 	return nil, status.Errorf(codes.Unimplemented, "method Set not implemented")
// }
// func (proxy *Proxy) Remove(context.Context, *proto.WriteRequest) (*proto.WriteResponse, error) {
// 	return nil, status.Errorf(codes.Unimplemented, "method Remove not implemented")
// }
