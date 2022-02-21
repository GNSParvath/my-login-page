package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/fission/go-login-page/handlers"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

func main() {
	fmt.Println("starting the login application")

	r := mux.NewRouter()

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://foo.com:8080"},
		AllowCredentials: false,

		// Enable Debugging for testing, consider disabling in production
		Debug: true,
	})

	handler := c.Handler(r)
	// handler = handlers.CORS(
	// 	// handlers.AllowedMethods([]string{"GET", "POST", "PUT"}),
	// 	handlers.AllowedHeaders([]string{"Accept", "Accept-Language", "Content-Type", "Content-Language", "Origin"}),
	// 	// handlers.AllowedOrigins([]string{"*"}),
	// )(handler)

	// Insert the middleware

	r.HandleFunc("/api/login", handlers.UserLogin).Methods("POST")
	r.HandleFunc("/api/user/{userId}", handlers.GetUserDetails).Methods("GET")

	// fmt.Println(getHash1("anil123"))
	log.Fatal(http.ListenAndServe("localhost:8080", handler))

}

// func getHash1(pwd string) string {
// 	hash, err := bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.MinCost)
// 	if err != nil {
// 		log.Println(err)
// 		fmt.Println(hash)

// 	}
// 	return string(hash)

// }
