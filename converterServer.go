package main

import (
	pb "vidConv/proto"
)

type server struct {
	serv pb.UnimplementedConverterServiceServer
	conv converter
}

func (s server) ConvertVideo(req *pb.ConvertRequest, srv pb.ConverterService_ConvertVideoServer) error {
	resp := &pb.ConvertResponse{
		Buffer: "",
		Part:   0,
	}

	if err := srv.Send(resp); err == nil {

	}
	return nil
}
