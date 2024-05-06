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
	"github.com/raian621/coverdb/database"
	// oapi_middleware "github.com/oapi-codegen/nethttp-middleware"
)

func main() {
	host := flag.String("host", "127.0.0.1", "host to listen on")
	port := flag.String("port", "8000", "port to listen on")
	certfile := flag.String("certfile", "coverdb.crt", "x.509 certificate file")
	keyfile := flag.String("keyfile", "coverdb.key", "x.509 private key file")
	https := flag.Bool("https", false, "use https")
	adminUser := flag.String("admin", "admin", "initial admin username")
	password := flag.String("password", "password", "initial admin password")
	flag.Parse()

	err := database.CreateDB("database/coverdb.db", "database/schema.sql")
	if err != nil {
		panic(err)
	}
	if err := database.CreateAdminUser(*adminUser, *password); err != nil {
		panic(err)
	}
	s := CreateServer(*host, *port)

	if *https {
		log.Printf("listening on https://%s:%s\n", *host, *port)
		if err := EnsureTLSFiles(*certfile, *keyfile); err != nil {
			panic(err)
		}
		log.Fatal(s.ListenAndServeTLS(*certfile, *keyfile))
	} else {
		log.Printf("listening on http://%s:%s\n", *host, *port)
		log.Fatal(s.ListenAndServe())
	}
}

func CreateServer(host, port string) *http.Server {
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

	return &http.Server{
		Handler: r,
		Addr:    net.JoinHostPort(host, port),
	}
}
