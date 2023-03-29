package main

import (
	"fmt"
	"net/http"
	"os"
)

func main() {
	fmt.Println("Serving HTTP on localhost:8000")
	if err := http.ListenAndServe(":8000", http.FileServer(http.Dir("."))); err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
}
