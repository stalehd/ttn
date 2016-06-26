package main

import (
    "fmt"
    "crypto/rand"
    "net/http"
    "encoding/json"
    "encoding/hex"
    "io/ioutil"
    "os"
)

// OAuthToken contains a mock token.
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

// ApplicationEntry contains the application meta-data.
type ApplicationEntry struct {
    eui        []byte
    EUI        string
    Name       string
    Owner      string
    AccessKeys []string
    Valid      bool
}

var applications []ApplicationEntry

func rootHandler(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "Hello, hello!")
}

func genRandHex(len int) []byte {
    bytes := make([]byte, len)
    _, err := rand.Read(bytes)
    if (err != nil) {
        fmt.Printf("Error generating random number: %s\n", err)
        return nil
    }
    return bytes
}

func (app *ApplicationEntry) GenerateEUI() {
    //TODO: Ensure this is unique
    app.eui = genRandHex(8)
    app.EUI = hex.EncodeToString(app.eui)
}

func (app *ApplicationEntry) GenerateAccessKey() {
    newAccessKey := hex.EncodeToString(genRandHex(16))
    app.AccessKeys = append(app.AccessKeys, newAccessKey)
}

const handlerEndpoint = "localhost:1882"

func applicationsHandler(w http.ResponseWriter, r *http.Request) {
    statusCode := http.StatusOK

    if (r.Method == "POST") {
        r.ParseForm()
        name := r.Form["name"]
        app := ApplicationEntry{
            Name: name[0],
            Owner: "lora@telenordigital.com",
            Valid: true }
        app.GenerateEUI()
        app.GenerateAccessKey()
        applications = append(applications, app)
        statusCode = http.StatusCreated
        writeApplicationsFile()
    }
    jsonBytes, err := json.Marshal(applications)
    if (err != nil) {
        fmt.Fprintf(w, "Error getting apps: %s\n", err)
    }
    w.WriteHeader(statusCode)
    fmt.Fprintf(w, string(jsonBytes))
}

// FILENAME contains the name (and path) of the file used to persist the
// application data.
const FILENAME = "applications.json"

func readApplicationsFile() {
    bytes, err := ioutil.ReadFile(FILENAME)
    if err != nil {
        fmt.Println("Got error reading",FILENAME,":", err)
        return
    }
    err = json.Unmarshal(bytes, &applications)
    if err != nil {
        fmt.Println("Got error unmarshalling json:",err)
        return
    }
}

func writeApplicationsFile() {
    bytes, err := json.Marshal(applications)
    if err != nil {
        fmt.Println("Couldn't marshal JSON:", err)
        return
    }
    ioutil.WriteFile(FILENAME, bytes, os.ModeAppend|0700)
}

// PORT is the port that the mockend server will listen to.
const PORT = 8080

func main() {
    readApplicationsFile()

    fmt.Printf("TTN Mockend Server @ port %d\n", PORT)

    http.HandleFunc("/", rootHandler)
    http.HandleFunc("/key", tokenHandler)
    http.HandleFunc("/users/token", tokenHandler)
    http.HandleFunc("/applications", applicationsHandler)

    http.ListenAndServe(fmt.Sprintf(":%d", PORT), nil)
}

