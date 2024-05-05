package main

import (
	"flag"
	"log"
	"net"
	"net/http"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/raian621/coverdb/api"
	// oapi_middleware "github.com/oapi-codegen/nethttp-middleware"
)

func main() {
	host := flag.String("host", "127.0.0.1", "host to listen on")
	port := flag.String("port", "8000", "port to listen on")
	flag.Parse()

	swagger, err := api.GetSwagger()
	if err != nil {
		panic(err)
	}

	swagger.Servers = nil
	updatedPaths := openapi3.Paths{}
	for key, value := range swagger.Paths.Map() {
		updatedPaths.Set("/api/v1"+key, value)
	}
	swagger.Paths = &updatedPaths

	server := api.Server{}
	apiRouter := chi.NewRouter()
	// r.Use(oapi_middleware.OapiRequestValidator(swagger))
	apiRouter.Use(middleware.RequestID)
	apiRouter.Use(middleware.RealIP)
	apiRouter.Use(middleware.Logger)
	apiRouter.Use(middleware.Recoverer)

	r := chi.NewRouter()
	r.Mount("/api/v1", apiRouter)
	r.Handle("/*", http.FileServer(http.Dir("client/dist")))
	api.HandlerFromMux(server, apiRouter)

	s := &http.Server{
		Handler: r,
		Addr:    net.JoinHostPort(*host, *port),
	}

	log.Fatal(s.ListenAndServe())
}
