package gitnames

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

type User struct {
	Login       string    `json:"login"`
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	Bio         string    `json:"bio"`
	PublicRepos int       `json:"public_repos"`
	PublicGists int       `json:"public_gists"`
	Followers   int       `json:"followers"`
	Following   int       `json:"following"`
	CreatedAt   time.Time `json:"created_at"`
}

type Repository struct {
	Name     string `json:"name"`
	Language string `json:"language"`
	Forks    int    `json:"forks"`
	Created  string `json:"created_at"`
	Updated  string `json:"updated_at"`
}

func fetchUsersData(usernames []string) {
	url := "https://api.github.com/users/" + usernames[0]
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Add("Accept", "Accept: application/json")
	//req.Header.Add("Custom-Header", "Custom value")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Response status:", resp.Status)
	defer resp.Body.Close()
	//scanner := bufio.NewScanner(resp.Body)
	//for i := 0; scanner.Scan() && i < 15; i++ {
	//	fmt.Println(scanner.Text())
	//}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	var user User
	err = json.Unmarshal(body, &user)
	if err := json.Unmarshal(body, &user); err != nil {
		log.Fatal(err)
	}
}

func fetchRepos(usernames []string) {
	url := "https://api.github.com/users/" + usernames[0] + "/repos"
	req, err := http.NewRequest("GET", url, nil)

	if err != nil {
		log.Fatal(err)
	}

	req.Header.Add("Accept", "Accept: application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Response status:", resp.Status)
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	var repos []Repository
	err = json.Unmarshal(body, &repos)
	if err := json.Unmarshal(body, &repos); err != nil {
		log.Fatal(err)
	}
}
func main() {
	var usernames [3]string
	usernames[1] = "angeld55"
	usernames[0] = "YanaRGeorgieva"

}
