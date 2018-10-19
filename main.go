package main

import (
    "encoding/json"
    "fmt"
    "log"
    "net/http"
    "github.com/dgrijalva/jwt-go"
    "github.com/gorilla/mux"
    "github.com/mitchellh/mapstructure"
    "github.com/tidwall/buntdb"
    "time"
)

type User struct {
    Username string `json:"username"`
    Password string `json:"password"`
}

type JwtToken struct {
	Username string 
    Token string `json:"token"`
    RefreshToken string `json:"refreshToken"`
}

type Exception struct{
    Message string `json:"message"`
}

// global database
var db *buntdb.DB

func CreateTokenEndpoint(w http.ResponseWriter, req *http.Request) {
    var user User
    _ = json.NewDecoder(req.Body).Decode(&user)
    
    // validate user credentials
    user_valid := false
    err := db.View(func(tx *buntdb.Tx) error {
		val, err := tx.Get(user.Username)	
		if (err != nil){ 
			return err }
		if (val == user.Password){ 
			user_valid = true }
		return nil
	})
    if err != nil { fmt.Println(err) }
    
    if user_valid == true {
	    accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
	        "username":  user.Username,
	        "exp": 		 time.Now().Local().Add(time.Hour * time.Duration(24)).Unix(),
	    })
	    
	    refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
	        "username":  user.Username,
	        "exp": 		 time.Now().Local().Add(time.Hour * time.Duration(24) * 7).Unix(),
	    })
	    
	    accessTokenString, error := accessToken.SignedString([]byte("secret"))
	    if (error != nil) { fmt.Println(error) }
	    
	    refreshTokenString, error2 := refreshToken.SignedString([]byte("refreshSecret"))
	    if (error2 != nil) { fmt.Println(error)}
	   
	    json.NewEncoder(w).Encode(JwtToken{Username: user.Username, Token: accessTokenString, RefreshToken: refreshTokenString})
	}
}

func ProtectedEndpoint(w http.ResponseWriter, req *http.Request) {
    params := req.URL.Query()
    
    token, _ := jwt.Parse(params["token"][0], func(token *jwt.Token) (interface{}, error) {
        if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
            return nil, fmt.Errorf("There was an error")
        }
        return []byte("secret"), nil
    })
    
    if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
        var user User
        mapstructure.Decode(claims, &user)
        json.NewEncoder(w).Encode("'" + user.Username + "' you made it in the Protected Zone... YAY!")
    } else {
        json.NewEncoder(w).Encode(Exception{Message: "Invalid authorization token"})
    }
}


func main() {
	// Create in-memory database
	db, _ = buntdb.Open(":memory:")
	
	// Add username and passwords
	db.Update(func(tx *buntdb.Tx) error {
		tx.Set("brett.lee", "W3lcome!", nil)
		tx.Set("michael.jordan", "P@ssword1", nil)
		tx.Set("abraham.lincoln", "Password!@#", nil)
		return nil
	})
	
    router := mux.NewRouter()
    fmt.Println("Starting the application...")
    router.HandleFunc("/authenticate", CreateTokenEndpoint).Methods("POST")
    router.HandleFunc("/protected", ProtectedEndpoint).Methods("GET")
    log.Fatal(http.ListenAndServe(":8080", router))
}
