// Copyright 2017 tsuru authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/tsuru/tsuru-usage/exporter"
)

func main() {
	port := os.Getenv("PORT")
	tsuruEndpoint := os.Getenv("TSURU_HOST")
	tsuruToken := os.Getenv("USAGE_USER_TOKEN")
	tsuruServicesStr := os.Getenv("USAGE_SERVICES")
	if port == "" {
		port = "8080"
	}
	if tsuruEndpoint == "" {
		log.Fatal("Must set tsuru endpoint with TSURU_HOST env")
	}
	if tsuruToken == "" {
		log.Fatal("Must set tsuru token with USAGE_USER_TOKEN env")
	}
	var services []string
	if tsuruServicesStr != "" {
		services = strings.Split(tsuruServicesStr, ",")
	}
	exporter.Register(tsuruEndpoint, tsuruToken, services)
	runServer(port)
}

func router() http.Handler {
	r := mux.NewRouter()
	r.Handle("/metrics", promhttp.Handler())
	http.Handle("/metrics", promhttp.Handler())
	return r
}

func runServer(port string) {
	http.Handle("/", router())
	log.Printf("HTTP server listening at :%s...\n", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
