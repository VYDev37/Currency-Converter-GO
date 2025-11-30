package routes

import (
	pb "curr-conv-api/proto"
	"encoding/json"
	"fmt"
	"net/http"
)

type HTTPFunc interface {
	RegisterRoute(method, endpoint string, function func(http.ResponseWriter, *http.Request))
	Listen(port uint16) error
}

type HTTPServer struct {
	Client pb.ConverterServiceClient
	mux    *http.ServeMux
}

func CreateHTTPServer(client pb.ConverterServiceClient) *HTTPServer {
	return &HTTPServer{
		Client: client,
		mux:    http.NewServeMux(),
	}
}

func WriteJSON(res http.ResponseWriter, status int, v any) error {
	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(status)
	return json.NewEncoder(res).Encode(v) // json structure for message
}

func SendMessage(res http.ResponseWriter, message string, status int) error {
	return WriteJSON(res, status, map[string]string{"message": message}) // json structure for message
}

func AllowCORS(next http.Handler) http.Handler {
	allowedOrigins := map[string]bool{ // domain_name: can access / not
		"http://localhost:5173": true,
	}

	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		origin := req.Header.Get("Origin")
		if allowedOrigins[origin] {
			res.Header().Set("Access-Control-Allow-Origin", origin)
		}
		res.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
		res.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if req.Method == "OPTIONS" {
			res.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(res, req)
	})
}

func (server *HTTPServer) RegisterRoute(method, endpoint string, function func(http.ResponseWriter, *http.Request)) {
	server.mux.HandleFunc(fmt.Sprintf("%s %s", method, endpoint), function)
}

func (server *HTTPServer) RegisterRoutes() {
	server.RegisterRoute("GET", "/", func(res http.ResponseWriter, req *http.Request) {
		fmt.Fprintf(res, "Hello world!")
	})
	server.RegisterRoute("GET", "/get-currencies", server.HandleGetCurrencies)
	server.RegisterRoute("POST", "/convert", server.HandleConversion)
}

func (server *HTTPServer) Listen(port uint16) error {
	fmt.Println("Registered routes.")
	fmt.Printf("Server is running on port %d.\n", port)

	server.RegisterRoutes()

	if err := http.ListenAndServe(fmt.Sprintf(":%d", port), AllowCORS(server.mux)); err != nil {
		return err
	}

	return nil
}
