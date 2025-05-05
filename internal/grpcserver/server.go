package grpcserver

import (
	"io"
	"log"

	"nova-panel/internal/store"
	"nova-panel/pb"
)

type Server struct {
	pb.UnimplementedVpsServer
}

func (s *Server) ReportStatus(stream pb.Vps_ReportStatusServer) error {
	for {
		status, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Printf("Agent 状态接收错误: %v", err)
			break
		}
		store.UpdateStatus(status)
	}
	return stream.SendAndClose(&pb.StatusAck{Message: "接收完毕"})
}
