package main

import (
	"context"
	"database/sql"
	"log"
	"net"

	"google.golang.org/grpc"
	pb "xmlservice/pb"
)

type grpcServer struct {
	pb.UnimplementedBIServiceServer
	db *sql.DB
}

func (s *grpcServer) SalesByCategory(ctx context.Context, req *pb.CategoryRequest) (*pb.CategoryResponse, error) {
	q := `
	WITH last_doc AS (
	  SELECT xml_documento
	  FROM relatorio
	  ORDER BY id DESC
	  LIMIT 1
	),
	vals AS (
	  SELECT unnest(xpath('//Item[@Categoria=' || quote_literal($1) || ']/ValorVenda/text()', xml_documento))::text AS v
	  FROM last_doc
	)
	SELECT COALESCE(SUM((v)::numeric), 0) FROM vals;
	`

	var total float64
	if err := s.db.QueryRow(q, req.Category).Scan(&total); err != nil {
		return nil, err
	}

	return &pb.CategoryResponse{
		Category: req.Category,
		Total:    total,
	}, nil
}

func startGRPC(db *sql.DB) error {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		return err
	}

	grpcSrv := grpc.NewServer()

	pb.RegisterBIServiceServer(grpcSrv, &grpcServer{db: db})
	
	log.Println("gRPC server escutando em :50051")
	return grpcSrv.Serve(lis)
}
