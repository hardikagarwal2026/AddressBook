package main

import (
	"fmt"
	"log"
	"net/http"
	"addressbook/handlers"
)

func main(){
	//Handlefunc is use to register a handler function for a specific URL Path
	//syntax - http.HandleFunc(pattern string, handler func(http.ResponseWriter, *http.Request))
	//ye basically hellohandler ko register kr dega for all requests for "/" 
	http.HandleFunc("/",handlers.HelloHandler) // esentially a routing of request to its handler funtions
	http.HandleFunc("/contacts", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet :
			handlers.GetcontactsHandler(w, r)
		case http.MethodPost :
			handlers.CreateContactHandler(w, r)
		default:
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}
	})

	http.HandleFunc("/contacts/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet :
				handlers.GetContactByIDHandler(w, r)
		case http.MethodPut :
			  	handlers.UpdateContactHandler(w, r)
		case http.MethodDelete :
				handlers.DeleteContactHandler(w, r)
		default :
		 	http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}
	})

	http.HandleFunc("/contacts/search", handlers.SearchContactHandler)

	fmt.Println("Server satarting on:8080")

	// syntax - http.ListenAndServe(addr string, handler http.Handler) error
	// it takes the address and port to listen on, and http.handler ususally nil, uses defaultservermux registered via http.HandleFunc or http.Handle
	// hindi mae bolu to , http server port 8080 pr chla dega
	// nil matlab default multiplexer, go ka default router

	handler1 := handlers.MiddlewareLogger(http.DefaultServeMux)
	handler2 := handlers.MiddlewareCORS(handler1)
	log.Fatal(http.ListenAndServe(":8080", handler2))

}

