package main

import (
	"net/http"
)

//go:generate go run ./cmd/assemble/main.go

func main() {
	if err := http.ListenAndServe(":8080", http.FileServer(http.Dir("./dist"))); err != nil {
		panic(err)
	}
}
