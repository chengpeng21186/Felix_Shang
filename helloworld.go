package main

import "net/http"

func helloWorld(w http.ResponseWriter, r *http.Request)  {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	w.Write([]byte("Hello World!!!\n"))
}

func welcomeTo(v http.ResponseWriter, r *http.Request) {
	v.Write([]byte("Welcome to you!!!\n"))
}
func main() {
	http.HandleFunc("/", helloWorld)
	http.HandleFunc("/wel", welcomeTo)
	http.ListenAndServe(":8000", nil)
}
