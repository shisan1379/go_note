package ugserver

import (
	"context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"user_growth/pb"
)

type UgCoinServer struct {
	pb.UnimplementedUserCoinServer
}

func (s *UgCoinServer) ListTask(ctx context.Context, in *pb.ListTaskRequest) (*pb.ListTaskReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "方法未实现")
}
func (s *UgCoinServer) UserCoinInfo(ctx context.Context, in *pb.UserCoinInfoRequest) (*pb.UserCoinInfoReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "方法未实现")
}
func (s *UgCoinServer) UserDetails(ctx context.Context, in *pb.UserDetailsRequest) (*pb.UserDetailsReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "方法未实现")
}
func (s *UgCoinServer) UserCoinChange(ctx context.Context, in *pb.UserCoinChangeRequest) (*pb.UserCoinChangeReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "方法未实现")
}
