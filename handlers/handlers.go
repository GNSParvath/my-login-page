package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/fission/go-login-page/configs"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

func UserLogin(w http.ResponseWriter, r *http.Request) {

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	params := make(map[string]string)
	postBody, _ := ioutil.ReadAll(r.Body)
	json.Unmarshal(postBody, &params)
	fmt.Println(params)
	fmt.Println(getHash([]byte(params["password"])))
	w.Header().Set("Content-Type", "application/bson")

	w.Header().Set("Access-Control-Allow-Origin", "http://example.com")
	w.Header().Set("Access-Control-Max-Age", "86400")
	//coll := client.Database("admin_panel").Collection("user")

	//ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	//err := coll.FindOne(ctx, bson.M{"email": vars["email"]}).Decode(&dbUser)
	_, dbpassword, err := configs.GetUser(ctx, params["email"])
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"message":"` + err.Error() + `"}`))
		return
	}
	userPass := []byte(params["password"])
	dbPass := []byte(dbpassword)

	passErr := bcrypt.CompareHashAndPassword(dbPass, userPass)

	if passErr != nil {
		log.Println(passErr)
		w.Write([]byte(`{"response":"Wrong Password!"}`))
		return
	}
	jwtToken, err := GenerateJWT()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"message":"` + err.Error() + `"}`))
		return
	}

	w.Write([]byte(`{"token":"` + jwtToken + `"}`))

}

func GetUserDetails(w http.ResponseWriter, r *http.Request) {

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	params := mux.Vars(r)
	fmt.Println(params)
	err := TokenValid(r)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(`{"message":"` + err.Error() + `"}`))
		return
	}

	objectID, _ := primitive.ObjectIDFromHex(params["userId"])
	fmt.Println(objectID)
	resp, err := configs.GetUserByID(ctx, objectID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"message":"` + err.Error() + `"}`))
		return
	}
	w.Header().Set("Content-Type", "application/json")
	response, _ := json.Marshal(resp)
	w.WriteHeader(http.StatusOK)
	w.Write(response)
}

// func ChangePassword(w http.ResponseWriter, r *http.Request) {

// 	params := make(map[string]string)
// 	postBody, _ := ioutil.ReadAll(r.Body)
// 	var user bson.M

// 	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
// 	defer cancel()

// 	params["password"] = getHash([]byte(params["Password"]))

// 	coll := GetCollection(DB, "users")
// 	_, err = coll.UpdateOne(ctx, bson.M{"email": email}).Decode(&user)

// 	if err != nil {
// 		fmt.Println(err)
// 		params["error"] = "an error encountered"
// 		return
// 	}
// 	json.Unmarshal(postBody, &params)
// 	fmt.Println(params)
// 	fmt.Println(getHash([]byte(params["password"])))
// 	w.Header().Set("Content-Type", "application/bson")
// }
