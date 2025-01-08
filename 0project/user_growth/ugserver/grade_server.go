package ugserver

import (
	"context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"user_growth/pb"
)

type UgGrowthServer struct {
	pb.UnimplementedUserGradeServer
}

func (s *UgGrowthServer) ListGrades(ctx context.Context, in *pb.ListGradesRequest) (*pb.ListGradesReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "方法未实现")
}
func (s *UgGrowthServer) ListGradePrivileges(ctx context.Context, in *pb.ListGradePrivilegesRequest) (*pb.ListGradePrivilegesReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "方法未实现")
}
func (s *UgGrowthServer) CheckUserPrivilege(ctx context.Context, in *pb.CheckUserPrivilegeRequest) (*pb.CheckUserPrivilegeReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "方法未实现")
}
func (s *UgGrowthServer) UserGradeInfo(ctx context.Context, in *pb.UserGradeInfoRequest) (*pb.UserGradeInfoReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "方法未实现")
}
func (s *UgGrowthServer) UserGradeChange(ctx context.Context, in *pb.UserGradeChangeRequest) (*pb.UserGradeChangeReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "方法未实现")
}
