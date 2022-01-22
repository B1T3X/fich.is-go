package main

import (
	"crypto/rand"
	"crypto/tls"
	"encoding/base64"
	"fmt"
	"log"
	"math/big"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/gorilla/mux"
)

var httpsPort string = os.Getenv("FICHIS_HTTPS_PORT")
var certFile string = os.Getenv("FICHIS_CERTIFICATE_FILE_PATH")
var keyFile string = os.Getenv("FICHIS_CERTIFICATE_KEY_PATH")

var redisHost string = os.Getenv("REDIS_HOST")
var redisPort string = os.Getenv("REDIS_PORT")

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

// func getAPIKey(filename string) (key string, err error) {
// 	content, err := ioutil.ReadFile(filename)

// 	key = string(content)

// 	return key, err
// }

// This function is needed in order to bypass Go only listening on IPv6 by default
func listenOnIPv4() (router *mux.Router, server *http.Server, err error) {
	router = mux.NewRouter()
	address := fmt.Sprintf("0.0.0.0:%v", httpsPort)
	log.Printf("Going to listen of %v", address)
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
	return
}

// Get link, do not redirect
func apiGetLinkHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	url, err := getLink(id)
	fmt.Println(url)

	if url == "" || err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	w.Write([]byte(url))
}

// Delete link
func apiDeleteLinkHandler(w http.ResponseWriter, r *http.Request) {
	shortId := r.URL.Query().Get("id")
	err := deleteLink(shortId)

	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
}

// Add link by specified id
func apiAddLinkHandler(w http.ResponseWriter, r *http.Request) {
	shortId := r.URL.Query().Get("shortId")
	url := r.URL.Query().Get("url")

	if shortId == "" || url == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if !IsUrl(url) {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("URL is invalid"))
		return
	}

	generatedLink, err := addLink(shortId, url)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.Write([]byte(generatedLink))
}

// Add link by randomly generated Base64 id
func apiAutoAddLinkHandler(w http.ResponseWriter, r *http.Request) {
	shortId, err := GenerateRandomShortId(6)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Something went wrong with automatic shortId creation"))
	}
	url := r.URL.Query().Get("url")

	if shortId == "" || url == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if !IsUrl(url) {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("URL is invalid"))
		return
	}

	generatedLink, err := addLink(shortId, url)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.Write([]byte(generatedLink))
}

// Redirect to URL
func redirectLinkHandler(w http.ResponseWriter, r *http.Request) {
	shortId := mux.Vars(r)["shortId"]

	url, _ := getLink(shortId)
	// if url == "" {
	// 	w.WriteHeader(http.StatusNotFound)
	// 	w.Write([]byte("Link does not exist"))
	// 	fmt.Println("Null url")
	// 	return
	// }

	log.Printf("Redirecting from %v to %v\n", string(domainName+shortId), url)

	w.Header().Set("location", url)
	http.Redirect(w, r, url, http.StatusMovedPermanently)
}

func main() {
	log.Println("Starting listener...")
	r, srv, err := listenOnIPv4()
	if err != nil {
		panic(err)
	}

	log.Println("Listener started successfully!")

	log.Println("Configuring handlers...")
	r.HandleFunc("/api/get/linkByShortId", apiGetLinkHandler).Methods("GET")
	r.HandleFunc("/api/delete/linkByShortId", apiDeleteLinkHandler).Methods("DELETE")
	r.HandleFunc("/api/create/ShortenedLink", apiAddLinkHandler).Methods("POST")
	r.HandleFunc("/api/create/AutoShortenedLink", apiAutoAddLinkHandler).Methods("POST")
	r.HandleFunc("/{shortId}", redirectLinkHandler).Methods("GET")

	log.Println("Done!\nRunning.")

	err = srv.ListenAndServeTLS(certFile, keyFile)

	if err != nil {
		panic(err)
	}
}
