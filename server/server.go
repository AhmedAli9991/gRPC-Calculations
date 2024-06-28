package server

import (
	"context"
	"log"

	"github.com/ahmedalialphasquad123/calculationService/handlers"
	pb "github.com/ahmedalialphasquad123/calculationService/proto"
	"google.golang.org/protobuf/types/known/structpb"
)

type Server struct {
	pb.UnimplementedCalculationServiceServer
}

func (s *Server) AggregationObject(ctx context.Context, req *pb.AggregationObjectRequest) (*pb.CalculationResponse, error) {
	response, err := handlers.HandleAggregationObject(req)
	if err != nil {
		log.Printf("Error handling AggregationObject: %v", err)
		return &pb.CalculationResponse{
			Data: &structpb.Struct{
				Fields: map[string]*structpb.Value{
					"message": structpb.NewStringValue(err.Error()),
				},
			},
		}, err
	}
	return response, nil
}

func (s *Server) AggregationMaterialData(ctx context.Context, req *pb.CalculationRequest) (*pb.CalculationResponse, error) {
	return handlers.HandleRequest(req)
}
