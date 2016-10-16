package main

import (
    "log"
    "net/http"
    "fmt"


)


var logger *log.Logger
func main() {

    router := NewRouter()
    fmt.Println("Starting WebService")
    log.Fatal(http.ListenAndServe(":8080", router))


}
