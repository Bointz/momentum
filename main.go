package main

import (
    "encoding/json"
    "fmt"
    "log"
    "net/http"
    "github.com/dgrijalva/jwt-go"
    "github.com/gorilla/mux"
    "github.com/mitchellh/mapstructure"
    "time"
)

type Claim struct {
    Username string `json:"username"`
    Password string `json:"password"`
}

type JwtToken struct {
    Token string `json:"token"`
    ExpireToken time.Time `json:"expiry"`
}

type Exception struct{
    Message string `json:"message"`
}

func CreateTokenEndpoint(w http.ResponseWriter, req *http.Request) {
    var claim Claim
    _ = json.NewDecoder(req.Body).Decode(&claim)

    
    expiry := time.Now().Local().Add(time.Hour * time.Duration(24))
    
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
        "username":  claim.Username,
        "password":  claim.Password,
        "expire_at": expiry,
    })
    tokenString, error := token.SignedString([]byte("secret"))
    if error != nil {
        fmt.Println(error)
    }
   
    json.NewEncoder(w).Encode(JwtToken{Token: tokenString })
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
        var claim Claim
        mapstructure.Decode(claims, &claim)
        json.NewEncoder(w).Encode("'" + claim.Username + "' you made it in the Protected Zone... YAY!")
    } else {
        json.NewEncoder(w).Encode(Exception{Message: "Invalid authorization token"})
    }
}


func main() {
    router := mux.NewRouter()
    fmt.Println("Starting the application...")
    router.HandleFunc("/authenticate", CreateTokenEndpoint).Methods("POST")
    router.HandleFunc("/protected", ProtectedEndpoint).Methods("GET")
    log.Fatal(http.ListenAndServe(":8080", router))
}
