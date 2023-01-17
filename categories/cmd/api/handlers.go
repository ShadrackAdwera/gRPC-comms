package main

import (
	"categories/protobufs"
	"categories/repo"
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type JsonResponse struct {
	Error   bool   `json:"error"`
	Message string `json:"message"`
	Data    any    `json:"data"`
}

type JsonRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	CreatedBy   string `json:"createdy"`
}

func (app *Config) FetchCategories(w http.ResponseWriter, r *http.Request) {
	data, err := app.Models.Category.GetAll()

	if err != nil {
		app.errJSON(w, err)
		return
	}

	response := jsonResponse{
		Message: "Found Found",
		Error:   false,
		Data:    data,
	}
	app.writeJSON(w, http.StatusOK, response)
}

func (app *Config) AddCategory(w http.ResponseWriter, r *http.Request) {
	var newCategory JsonRequest

	err := app.readJSON(w, r, &newCategory)

	if err != nil {
		app.errJSON(w, err)
		return
	}

	categoryData := repo.Category{
		Name:        newCategory.Name,
		Description: newCategory.Description,
		CreatedBy:   newCategory.CreatedBy,
	}

	conn, err := grpc.Dial("products-service:8000", grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithBlock())

	if err != nil {
		app.errJSON(w, err)
		return
	}

	defer conn.Close()

	c := protobufs.NewCategoryServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)

	defer cancel()

	_, err = c.WriteCategory(ctx, &protobufs.CategoryRequest{
		CategoryEntry: &protobufs.Category{
			Name:        newCategory.Name,
			Description: newCategory.Description,
			CreatedBy:   newCategory.CreatedBy,
		},
	})

	if err != nil {
		app.errJSON(w, err)
		return
	}

	log.Printf("Data sent via gRPC. . . ")

	res, err := app.Models.Category.Insert(categoryData)

	response := jsonResponse{
		Error:   false,
		Message: fmt.Sprintf("Category with ID: %v has been added", res),
	}

	app.writeJSON(w, http.StatusCreated, response)
	// send to products service via gRPC
}

func (app *Config) PatchCategory() {
	// send to products service via gRPC
}

func (app *Config) DeleteCategory() {
	// send to products service via gRPC
}
