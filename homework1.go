package gitnames

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
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

func fetchUsersData(username string) (User, error) {
	url := "https://api.github.com/users/" + username
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
	var user User
	if err := json.Unmarshal(body, &user); err != nil {
		log.Fatal(err)
	}
	return user, err
}

func fetchRepos(username string) ([]Repository, error) {
	url := "https://api.github.com/users/" + username + "/repos"
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
	if err := json.Unmarshal(body, &repos); err != nil {
		log.Fatal(err)
	}

	return repos, err
}

func fetchLanguages(username, repo string) (map[string]int, error) {
	url := "https://api.github.com/repos/" + username + "/" + repo + "/languages"
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

	var languages map[string]int

	if err = json.Unmarshal(body, &languages); err != nil {
		log.Fatal(err)
	}

	return languages, nil
}

func setup(filePath string) ([]string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var usernames []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		usernames = append(usernames, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return usernames, nil
}

func main() {
	filePath := "D:\\goprojects\\gitnames.txt" // Replace with the path to your file
	usernames, err := setup(filePath)
	if err != nil {
		fmt.Println("Error reading file:", err)
		return
	}

	for _, word := range usernames {
		fmt.Print(word)
	}
}
