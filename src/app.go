package main

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"math/big"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/gorilla/mux"
)

var httpsPort string = os.Getenv("FICHIS_HTTPS_PORT")
var certificateFilePath string = os.Getenv("FICHIS_CERTIFICATE_FILE_PATH")
var privateKeyFilePath string = os.Getenv("FICHIS_PRIVATE_KEY_FILE_PATH")

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
	server = &http.Server{
		Handler:      router,
		Addr:         fmt.Sprintf("127.0.0.1:%d", httpsPort),
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

	fmt.Printf("Redirecting from %v to %v\n", string(domainName+shortId), url)

	w.Header().Set("location", url)
	http.Redirect(w, r, url, http.StatusFound)
}

func main() {

	r, srv, err := listenOnIPv4()
	if err != nil {
		panic(err)
	}

	r.HandleFunc("/api/get/linkByShortId", apiGetLinkHandler).Methods("GET")
	r.HandleFunc("/api/delete/linkByShortId", apiDeleteLinkHandler).Methods("DELETE")
	r.HandleFunc("/api/create/ShortenedLink", apiAddLinkHandler).Methods("POST")
	r.HandleFunc("/api/create/AutoShortenedLink", apiAutoAddLinkHandler).Methods("POST")
	r.HandleFunc("/{shortId}", redirectLinkHandler).Methods("GET")

	err = srv.ListenAndServeTLS(certificateFilePath, privateKeyFilePath)

	if err != nil {
		panic(err)
	}
}
