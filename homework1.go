package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/jedib0t/go-pretty/table"
	"io"
	"log"
	"net/http"
	"os"
	"sort"
	"time"
)

type User struct {
	Login       string    `json:"login"`
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	PublicRepos int       `json:"public_repos"`
	PublicGists int       `json:"public_gists"`
	Followers   int       `json:"followers"`
	Following   int       `json:"following"`
	CreatedAt   time.Time `json:"created_at"`
}

type Repository struct {
	Name    string `json:"name"`
	Forks   int    `json:"forks"`
	Created string `json:"created_at"`
	Updated string `json:"updated_at"`
}

func fetchUsersData(username string) (User, error) {
	url := "https://api.github.com/users/" + username
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Add("Authorization", "Bearer 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")
	req.Header.Add("User-Agent", "LearningToFetch")
	req.Header.Add("Accept", "Accept: application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	//fmt.Println("Response status:", resp.Status)
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
	req.Header.Add("Authorization", "Bearer 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")
	req.Header.Add("Accept", "Accept: application/json")
	req.Header.Add("User-Agent", "ForGoCourse")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	//fmt.Println("Response status:", resp.Status)
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
	req.Header.Add("Authorization", "Bearer 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")
	req.Header.Add("Accept", "Accept: application/json")
	req.Header.Add("User-Agent", "ForGoCourse")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	//fmt.Println("Response statusForLanguage:", resp.Status)
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
		log.Fatal(err)
	}
	defer file.Close()

	var usernames []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		usernames = append(usernames, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	return usernames, nil
}

func main() {
	filePath := "D:\\goprojects\\gitnames.txt" // Replace with the path to your file
	usernames, err := setup(filePath)
	if err != nil {
		fmt.Print("error reading file:", err)
		return
	}

	for _, username := range usernames {
		user, err := fetchUsersData(username)
		if err != nil {
			log.Fatal(err)
		}
		printUserTable(user)
		repos, err := fetchRepos(username)
		if err != nil {
			log.Fatal(err)
		}

		printRepoTable(repos, username)
	}

}

func printUserTable(user User) {
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"Login", "ID", "Name", "Public Repos", "Public Gists", "Followers", "Following", "Created At"})
	t.AppendRow([]interface{}{
		user.Login,
		user.ID,
		user.Name,
		user.PublicRepos,
		user.PublicGists,
		user.Followers,
		user.Following,
		user.CreatedAt.Format("2006-01-02"),
	})
	t.Render()
}

func printRepoTable(repos []Repository, username string) {
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"Name", "Forks", "Created", "Updated", "Language 1", "Language 2", "Language 3", "Language 4", "Others"})
	var totalForks int
	for _, repo := range repos {
		totalForks += repo.Forks
		lang, err := fetchLanguages(username, repo.Name)
		if err != nil {
			log.Fatal(err)
		}
		proportion := calculateProportion(lang)
		t.AppendRow([]interface{}{repo.Name, repo.Forks, repo.Created, repo.Updated,
			fmt.Sprintf("%.2f%%", proportion["Language 1"]),
			fmt.Sprintf("%.2f%%", proportion["Language 2"]),
			fmt.Sprintf("%.2f%%", proportion["Language 3"]),
			fmt.Sprintf("%.2f%%", proportion["Language 4"]),
			fmt.Sprintf("%.2f%%", proportion["Others"])})
	}
	//t.AppendRow([]interface{}{fmt.Printf()})
	t.Render()
}

func calculateProportion(repoLanguages map[string]int) map[string]float64 {
	totalBytes := 0

	for _, count := range repoLanguages {
		totalBytes += count
	}

	var languages []string
	for lang := range repoLanguages {
		languages = append(languages, lang)
	}
	sort.Slice(languages, func(i, j int) bool {
		return repoLanguages[languages[i]] > repoLanguages[languages[j]]
	})

	proportion := make(map[string]float64)
	numLanguages := min(4, len(languages)) // Take the minimum of 4 and the actual number of languages
	for i := 0; i < numLanguages; i++ {
		proportion[languages[i]] = float64(repoLanguages[languages[i]]) / float64(totalBytes) * 100
	}

	proportion["Others"] = float64(totalBytes)
	for i := 0; i < numLanguages; i++ {
		proportion["Others"] -= float64(repoLanguages[languages[i]])
	}
	proportion["Others"] /= float64(totalBytes) * 100

	return proportion
}
