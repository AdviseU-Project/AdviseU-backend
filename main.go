package main

import (
    "fmt"
    "net/http"
)

func main() {
    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        fmt.Fprintln(w, "Hello from AdviseU Backend!")
    })
    http.ListenAndServe(":8080", nil)
}
