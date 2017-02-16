package main

import (
	"fmt"
	"net/http"
	"os"
)

func main() {
	http.HandleFunc("/", hello)
	http.HandleFunc("/headers", echo)
	fmt.Println("listening...")
	err := http.ListenAndServe(":"+os.Getenv("PORT"), nil)
	if err != nil {
		panic(err)
	}
}

func hello(res http.ResponseWriter, req *http.Request) {
	fmt.Fprintln(res, "go, world")
}

func echo(res http.ResponseWriter, req *http.Request) {
	fmt.Fprintln(res, fmt.Sprintf("\nRequest Headers: %s \n", req.Header))
}
