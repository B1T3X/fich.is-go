package main

import (
	"crypto/rand"
	"crypto/tls"
	"encoding/base64"
	"fmt"
	"log"
	"math/big"
	"net"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
	"github.com/gorilla/mux"
)

var httpsPort string = os.Getenv("FICHIS_HTTPS_PORT")
var httpPort string = os.Getenv("FICHIS_HTTP_PORT")

var certFile string = os.Getenv("FICHIS_CERTIFICATE_FILE_PATH")
var keyFile string = os.Getenv("FICHIS_KEY_FILE_PATH")

var fichisTlsOn string = os.Getenv("FICHIS_TLS_ON")
var fichisApiValidationOn string = strings.ToLower(os.Getenv("FICHIS_API_VALIDATION_ON"))
var fichisApiKey string = os.Getenv("FICHIS_API_KEY")

var probePath string = os.Getenv("FICHIS_HEALTH_PROBE_PATH")

// Generates random Base64 IDs for apiAutoAddLinkHandler
const letters string = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz-"

func GenerateRandomShortId(n int) (string, error) {

	ret := make([]byte, n)
	for i := 0; i < n; i++ {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(letters))))
		if err != nil {
			return "", err
		}
		ret[i] = letters[num.Int64()]
	}
	return base64.URLEncoding.EncodeToString(ret), nil
}

// URL Validation
func IsUrl(str string) bool {
	u, err := url.Parse(str)
	return err == nil && u.Scheme != "" && u.Host != ""
}

func validateAPIKey(key string) (valid bool) {
	if (key == fichisApiKey && fichisApiValidationOn == "yes") || fichisApiValidationOn != "yes" {
		valid = true
	} else {
		valid = false
	}
	return
}

// This function is needed in order to bypass Go only listening on IPv6 by default
func listenOnIPv4(portToListenTo string) (router *mux.Router, server *http.Server, listener net.Listener, err error) {
	router = mux.NewRouter()
	address := fmt.Sprintf("0.0.0.0:%v", portToListenTo)
	log.Printf("Going to listen on %v", address)
	server = &http.Server{
		Handler: router,
		Addr:    address,
		TLSConfig: &tls.Config{
			MaxVersion: tls.VersionTLS13,
			MinVersion: tls.VersionTLS12,
		},
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	listener, err = net.Listen("tcp4", address)
	return
}

// Get link, do not redirect
func apiGetLinkHandler(w http.ResponseWriter, r *http.Request) {
	if !(validateAPIKey(r.URL.Query().Get("api_key"))) {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("Invalid API key"))
		return
	}
	id := r.URL.Query().Get("id")
	url, err := getLink(id)
	fmt.Println(url)

	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	w.Write([]byte(url))
}

// Delete link
func apiDeleteLinkHandler(w http.ResponseWriter, r *http.Request) {
	if !(validateAPIKey(r.URL.Query().Get("api_key"))) {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("Invalid API key"))
		return
	}
	id := r.URL.Query().Get("id")
	err := deleteLink(id)

	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
}

// Add link by specified id
func apiAddLinkHandler(w http.ResponseWriter, r *http.Request) {
	if !(validateAPIKey(r.URL.Query().Get("api_key"))) {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("Invalid API key"))
		return
	}
	id := r.URL.Query().Get("id")
	url := r.URL.Query().Get("url")

	if id == "" || url == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if !IsUrl(url) {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("URL is invalid"))
		return
	}

	generatedLink, err := addLink(id, url)

	if err != nil {
		w.WriteHeader(http.StatusConflict)
		w.Write([]byte("Link already exists"))
		return
	}

	w.Write([]byte(generatedLink))
}

// Add link by randomly generated Base64 id
func apiAutoAddLinkHandler(w http.ResponseWriter, r *http.Request) {
	if !(validateAPIKey(r.URL.Query().Get("api_key"))) {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("Invalid API key"))
		return
	}
	id, err := GenerateRandomShortId(6)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Something went wrong with automatic id creation"))
	}
	url := r.URL.Query().Get("url")

	if id == "" || url == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("URL not supplied"))
		return
	}

	if !IsUrl(url) {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("URL is invalid"))
		return
	}

	generatedLink, err := addLink(id, url)

	if err != nil {
		w.WriteHeader(http.StatusConflict)
		w.Write([]byte("Link already exists"))
		return
	}

	w.Write([]byte(generatedLink))
}

// Redirect to URL
func redirectLinkHandler(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	url, err := getLink(id)

	// TODO: Reimplement check if domain exists
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Link does not exist"))
		fmt.Println("Null url")
		return
	}

	log.Printf("Redirecting from %v to %v\n", string(domainName+id), url)

	w.Header().Set("location", url)
	http.Redirect(w, r, url, http.StatusMovedPermanently)
}

func sayHello(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("I'm here!"))
}

func main() {
	log.Println("Starting listener...")
	var portToListenTo *string

	if fichisTlsOn == "yes" {
		portToListenTo = &httpsPort
	} else {
		portToListenTo = &httpPort
	}
	r, srv, listener, err := listenOnIPv4(*portToListenTo)
	if err != nil {
		panic(err)
	}

	log.Println("Listener started successfully!")

	log.Println("Configuring handlers...")
	r.HandleFunc("/api/get/ShortenedLink", apiGetLinkHandler).Methods("GET")
	r.HandleFunc("/api/delete/ShortenedLink", apiDeleteLinkHandler).Methods("DELETE")
	r.HandleFunc("/api/create/ShortenedLink", apiAddLinkHandler).Methods("POST")
	r.HandleFunc("/api/create/AutoShortenedLink", apiAutoAddLinkHandler).Methods("POST")
	r.HandleFunc("/{id}", redirectLinkHandler).Methods("GET")
	r.HandleFunc(probePath, sayHello).Methods("GET")

	log.Println("Done!\nRunning.")

	if fichisTlsOn == "yes" {
		err = srv.ServeTLS(listener, certFile, keyFile)
	} else {
		err = srv.Serve(listener)
	}

	if err != nil {
		panic(err)
	}
}
