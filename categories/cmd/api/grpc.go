package main

import (
	"categories/protobufs"
	"categories/repo"
	"context"
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc"
)

type CategoryServer struct {
	protobufs.UnimplementedCategoryServiceServer
	Models repo.Models
}

func (c *CategoryServer) WriteCategory(ctx context.Context, req *protobufs.CategoryRequest) (*protobufs.CategoryResponse, error) {

	// get an input
	input := req.GetCategoryEntry()

	// write an input
	catData := repo.Category{
		Name:        input.Name,
		Description: input.Description,
		CreatedBy:   input.CreatedBy,
	}

	_, err := c.Models.Category.Insert(catData)

	if err != nil {
		res := &protobufs.CategoryResponse{Result: "failed to insert data into the db"}
		return res, err
	}

	res := &protobufs.CategoryResponse{Result: "Data saved successfully"}
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
