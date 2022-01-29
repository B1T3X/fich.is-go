package cli

import (
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

var domain string = "https://fich.is"

func shortenLink(shortenCmd *flag.FlagSet, id *string, url *string) {
	shortenCmd.Parse(os.Args[2:])

	if *id == "" || *url == "" {
		fmt.Printf("id and url required\n")
		for _, i := range os.Args[2:] {
			fmt.Println(i)
		}
		shortenCmd.PrintDefaults()
		os.Exit(1)
	}
	requestPath := domain + fmt.Sprintf("/api/create/ShortenedLink?shortId=%v&url=%v", *id, *url)
	resp, err := http.Post(requestPath, "text/*", nil)

	fmt.Printf("Sending request %v\n", requestPath)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		panic(err)
	}

	shortenedLink := string(body)

	fmt.Println(shortenedLink)
}

func deleteLink(deleteCmd *flag.FlagSet, id *string) {
	deleteCmd.Parse(os.Args[2:])

	if *id == "" {
		fmt.Println("Expected id")
		os.Exit(1)
	}

	client := &http.Client{}
	req, err := http.NewRequest("DELETE", domain+fmt.Sprintf("/api/delete/linkByShortId?id=%v", id), nil)
	if err != nil {
		panic(err)
	}

	resp, err := client.Do(req)

	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	statuscode := fmt.Sprint(resp.StatusCode)
	fmt.Printf("Status code is: %v\n", statuscode)
	if statuscode == "200" {
		fmt.Printf("Link ID %v deleted successfully\n", *id)
	} else {
		fmt.Printf("Link ID %v was not deleted for some reason\n", *id)
	}

	return
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Expected 'shorten' or 'delete'")
		os.Exit(1)
	}

	shortenCmd := flag.NewFlagSet("shorten", flag.ExitOnError)

	shortenId := shortenCmd.String("id", "", "ShortId")
	shortenUrl := shortenCmd.String("url", "", "URL")

	deleteCmd := flag.NewFlagSet("delete", flag.ExitOnError)
	deleteId := deleteCmd.String("id", "", "ShortId")

	switch os.Args[1] {
	case "shorten":
		shortenLink(shortenCmd, shortenId, shortenUrl)
	case "delete":
		deleteLink(deleteCmd, deleteId)
	default:
		fmt.Println("Mama mia!")
	}
}
