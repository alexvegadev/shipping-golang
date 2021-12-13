package routes

import (
	"encoding/json"
	"net/http"
	"shipping/controller"
	"shipping/core"
	"shipping/repository"
	"shipping/template"

	"github.com/gorilla/mux"
	"gopkg.in/validator.v2"
)

type Routes struct {
}

func (p *Routes) SetupRoutes(router *mux.Router) {
	enableCORS(router)
	//products POST ENDPOINT
	router.HandleFunc("/products", func(w http.ResponseWriter, r *http.Request) {
		var request core.ProductStock
		err := json.NewDecoder(r.Body).Decode(&request)
		var template template.MLTemplate
		item := template.GetProduct(request.ProductId)
		if item.Id == "" {
			respondWithError(core.ErrorMessage{StatusCode: 404, Message: "El producto no fue encontrado"}, w)
		} else if err != nil {
			respondWithError(err, w)
		} else if item.Shipping.LogisticType != "fulfillment" {
			respondWithError(core.ErrorMessage{StatusCode: 404, Message: "El producto no es de tipo fulfillment"}, w)
		} else {
			errs := validator.Validate(request)
			if errs != nil {
				respondWithError(core.ErrorMessage{StatusCode: 404, Message: "Algunos de los campos no cumplen con los patrones establecidos."}, w)
			} else {
				var productStock controller.ProductStockController
				var repo = repository.PostgreRepository{}
				limitReached := productStock.IsLimitReachedDeposit(int(request.Quantity), request.Location, repo)
				limitProduct := productStock.IsLimitProductCountReached(request.ProductId, request.Location, request.Deposit, repo)

				if limitReached {
					respondWithError(core.ErrorMessage{StatusCode: 404, Message: "Se ha alcanzado el límite de deposito"}, w)
				} else if limitProduct {
					respondWithError(core.ErrorMessage{StatusCode: 404, Message: "Se ha alcanzado el límite por producto"}, w)
				} else {
					productStock.Save(request, repo)
					respondWithOk(core.ProductStockResponse{StatusCode: 200, Message: "Se ingresó el producto al depósito", Data: request}, w)
				}
			}
		}
	}).Methods(http.MethodPost)
	//products PUT ENDPOINT
	router.HandleFunc("/products", func(w http.ResponseWriter, r *http.Request) {
		var request core.ProductStock
		err := json.NewDecoder(r.Body).Decode(&request)
		if err != nil {
			respondWithError(core.ErrorMessage{StatusCode: 404, Message: "Es necesario agregar el modelo del producto."}, w)
		} else {
			var productStock controller.ProductStockController
			var repo = repository.PostgreRepository{}
			ps, errs := productStock.DrawOutProduct(request, repo)
			if errs != nil {
				respondWithError(core.ErrorMessage{StatusCode: 404, Message: "El producto no existe en el depósito."}, w)
			} else {
				respondWithOk(ps, w)
			}
		}
	}).Methods(http.MethodPut)
	//products GET ENDPOINT
	router.HandleFunc("/products", func(w http.ResponseWriter, r *http.Request) {
		params := r.URL.Query()
		deposit := params.Get("deposit")
		location := params.Get("location")
		if deposit != "" && location != "" {
			var productStock controller.ProductStockController
			var repo = repository.PostgreRepository{}
			ps := productStock.GetProductsByDepositAndLocation(deposit, location, repo)
			respondWithOk(ps, w)
		} else {
			respondWithError(core.ErrorMessage{StatusCode: 404, Message: "Es necesario ingresar los parámetros deposit y location."}, w)
		}
	}).Methods(http.MethodGet)
	//products/find GET ENDPOINT
	router.HandleFunc("/products/find", func(w http.ResponseWriter, r *http.Request) {
		params := r.URL.Query()
		deposit := params.Get("deposit")
		productId := params.Get("productId")
		if deposit != "" && productId != "" {
			var productStock controller.ProductStockController
			var repo = repository.PostgreRepository{}
			ps := productStock.GetProductsByDepositAndProduct(deposit, productId, repo)
			respondWithOk(ps, w)
		} else {
			respondWithError(core.ErrorMessage{StatusCode: 404, Message: "Es necesario ingresar los parámetros deposit y productId."}, w)
		}
	}).Methods(http.MethodGet)

}

func enableCORS(router *mux.Router) {
	router.PathPrefix("/").HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost")
	}).Methods(http.MethodOptions)
	router.Use(middlewareCors)
}

func middlewareCors(next http.Handler) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, req *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", "http://localhost")
			w.Header().Set("Access-Control-Allow-Credentials", "true")
			w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
			w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
			w.Header().Set("Content-Type", "application/json")
			w.Header().Set("Accept", "application/json")
			next.ServeHTTP(w, req)
		})
}

func respondWithError(data interface{}, w http.ResponseWriter) {
	w.WriteHeader(http.StatusInternalServerError)
	json.NewEncoder(w).Encode(data)
}

func respondWithForbidden(data interface{}, w http.ResponseWriter) {
	w.WriteHeader(http.StatusForbidden)
	json.NewEncoder(w).Encode(data)
}

func respondWithOk(data interface{}, w http.ResponseWriter) {
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(data)
}
