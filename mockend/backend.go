package main

import (
    "fmt"
    "crypto/rand"
    "net/http"
    "encoding/json"
    "encoding/hex"
)

type OAuthToken struct {
    GrantType    string  `json:grant_type`
    AccessToken  string  `json:access_token`
    RefreshToken string  `json:refresh_token`
    Expires      int     `json:expires`
    ExpiresIn    int     `json:expires_in`
    ClientId     string  `json:client_id`
}

func getDummyToken() OAuthToken {
    return OAuthToken{
        GrantType:    "imahara",
        AccessToken:  "microsoft",
        RefreshToken: "60hz",
        Expires:      9999999999,
        ExpiresIn:    9999999999,
        ClientId:     "swell"}
}

func tokenHandler(w http.ResponseWriter, r *http.Request) {
    jsonBytes, err := json.Marshal(getDummyToken())
    if (err != nil) {
        fmt.Fprintf(w, "Error marshalling JSON")
        return
    }
    fmt.Fprintf(w, string(jsonBytes))
}


type ApplicationEntry struct {
    EUI        string
    Name       string
    Owner      string
    AccessKeys []string
    Valid      bool
    AppKey     string
}

var applications []ApplicationEntry

func rootHandler(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "Hello, hello!")
}

func genRandHex(len int) string {
    bytes := make([]byte, len)
    _, err := rand.Read(bytes)
    if (err != nil) {
        fmt.Printf("Error generating random number: %s\n", err)
        return ""
    }
    return hex.EncodeToString(bytes)
}

func (app *ApplicationEntry) GenerateEUI() {
    //TODO: Ensure this is unique
    app.EUI = genRandHex(8)
}

func (app *ApplicationEntry) GenerateAppKey() {
    app.AppKey = genRandHex(16)
}

func (app *ApplicationEntry) GenerateAccessKey() {
    app.AccessKeys = append(app.AccessKeys, genRandHex(32))
}

func applicationsHandler(w http.ResponseWriter, r *http.Request) {
    if (r.Method == "POST") {
        r.ParseForm()
        name := r.Form["name"]
        app := ApplicationEntry{
            Name: name[0],
            Owner: "lora@telenordigital.com",
            Valid: true }
        app.GenerateEUI()
        app.GenerateAppKey()
        app.GenerateAccessKey()
        applications = append(applications, app)
    }
    jsonBytes, err := json.Marshal(applications)
    if (err != nil) {
        fmt.Fprintf(w, "Error getting apps: %s\n", err)
    }
    fmt.Fprintf(w, string(jsonBytes))
}

const PORT = 8080

func main() {

    fmt.Printf("TTN Mockend Server @ port %d\n", PORT)

    http.HandleFunc("/", rootHandler)
    http.HandleFunc("/key", tokenHandler)
    http.HandleFunc("/users/token", tokenHandler)
    http.HandleFunc("/applications", applicationsHandler)

    http.ListenAndServe(fmt.Sprintf(":%d", PORT), nil)
}

