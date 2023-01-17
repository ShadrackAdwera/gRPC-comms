package main

import (
	"fmt"
	"net/http"
	"products/repo"
)

func (app *Config) GetProducts(w http.ResponseWriter, r *http.Request) {
	res, err := app.Models.ProductEntry.All()

	if err != nil {
		app.errJSON(w, err)
		return
	}

	response := jsonResponse{
		Error:   false,
		Message: fmt.Sprintf("Found Products: %v", len(res)),
		Data:    res,
	}
	app.writeJSON(w, http.StatusOK, response)
}

func (app *Config) AddProduct(w http.ResponseWriter, r *http.Request) {
	var reqBody repo.ProductEntry

	err := app.readJSON(w, r, &reqBody)

	if err != nil {
		app.errJSON(w, err)
		return
	}
	err = app.Models.ProductEntry.Insert(reqBody)
	if err != nil {
		app.errJSON(w, err)
		return
	}
	response := jsonResponse{
		Error:   false,
		Message: "Product created successfully",
	}
	app.writeJSON(w, http.StatusOK, response)
}

func (app *Config) FetchGrpcCategories(w http.ResponseWriter, r *http.Request) {
	res, err := app.Models.CategoryEntry.GetCategories()

	if err != nil {
		app.errJSON(w, err)
		return
	}

	response := jsonResponse{
		Message: fmt.Sprintf("Found Categories: %v", len(res)),
		Error:   false,
		Data:    res,
	}

	app.writeJSON(w, http.StatusOK, response)

}
