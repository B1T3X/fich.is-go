package cli

import (
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	toml "github.com/pelletier/go-toml"
)
# Todo fix this
type Config struct {
	Url string
	APIKey string
}
var config_raw []byte
var err error
config_raw, err := os.ReadFile("~/.config/fich.is/config")

if err != nil {
	fmt.Print("Configuration at ~/.config/fich.is/config is unreadable")
	os.Exit(1)
}
config := Config{}
toml.Unmarshal(config_raw, &config)

	domain := config.Url
	apiKey = config.APIKey

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
	if apiKey != nil {
		requestPath := domain + fmt.Sprintf("/api/create/ShortenedLink?shortId=%v&url=%v", *id, *url)
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
