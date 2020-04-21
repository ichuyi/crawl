package main

import (
	"fmt"
	"net/http"
)

func test(w http.ResponseWriter,r *http.Request)  {
	fmt.Println(r.Header)
	fmt.Println(r.Host)
	w.Write([]byte("hello"))
}
func main()  {
	http.HandleFunc("/",test)
	http.ListenAndServe(":1234",nil)
}
