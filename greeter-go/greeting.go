/*
Copyright 2018 The Knative Authors

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"

	"flag"
)

var resource string

func main() {
	flag.Parse()
	router := mux.NewRouter().StrictSlash(true)

	resource = "greeting"
	root := "/api/" + resource

	router.HandleFunc("/", Index)
	router.HandleFunc(root, StockIndex)

	log.Fatal(http.ListenAndServe(":8080", router))
}

func Index(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	fmt.Fprintf(w, "{\"message\": \"Welcome to the %s app!\" }", resource)
}

func StockIndex(w http.ResponseWriter, r *http.Request) {
	envname := os.Getenv("ENV_NAME")
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	fmt.Fprintf(w, "{\"content\": \"%s\" }", envname)
}

