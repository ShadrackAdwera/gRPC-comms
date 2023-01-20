package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"products/protobufs"
	"products/repo"

	"google.golang.org/grpc"
)

type CategoryServer struct {
	protobufs.UnimplementedCategoryServiceServer
	Models repo.Models
}

func (c *CategoryServer) WriteCategory(ctx context.Context, req *protobufs.CategoryRequest) (*protobufs.CategoryResponse, error) {
	input := req.GetCategoryEntry()
	category := repo.CategoryEntry{
		Name:        input.Name,
		Description: input.Description,
		CreatedBy:   input.CreatedBy,
	}

	err := c.Models.CategoryEntry.AddCategory(category)

	if err != nil {
		res := &protobufs.CategoryResponse{Result: "error inserting category into the DB"}
		return res, err
	}
	res := &protobufs.CategoryResponse{Result: "category successfully added into the DB"}
	return res, nil
}

func (app *Config) grpcListen() {
	l, err := net.Listen("tcp", fmt.Sprintf(":%s", gRPCport))

	if err != nil {
		log.Fatalf("Fail to listen to gRPC connection: %v", err)
	}

	s := grpc.NewServer()

	protobufs.RegisterCategoryServiceServer(s, &CategoryServer{
		Models: app.Models,
	})

	log.Printf("gRPC server started on PORT: %s", gRPCport)

	if err := s.Serve(l); err != nil {
		log.Fatalf("Fail to listen to gRPC connection: %v", err)
	}

}
