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
		log.Printf("[Agent %s] 上报状态. CPU:%.2f 内存:%.2f 上传网络:%.2f 下载网速:%.2f, CPU:%s 硬盘:%.2f",
			status.AgentId, status.CpuPercent, status.MemoryPercent, status.UploadSpeedKbps, status.DownloadSpeedKbps,
			status.CpuModel, status.DiskPercent,
		)
		store.UpdateStatus(&store.AgentStatus{
			AgentId:      status.AgentId,
			CpuPercent:   status.CpuPercent,
			MemPercent:   status.MemoryPercent,
			UploadKbps:   status.UploadSpeedKbps,
			DownloadKbps: status.DownloadSpeedKbps,
		})
	}
	return stream.SendAndClose(&pb.StatusAck{Message: "接收完毕"})
}
