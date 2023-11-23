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
	"strings"
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

type Languages []struct {
	langs map[string]int
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
	url := "https://api.github.com/users/" + username + "/repos?per_page=80"
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

func fetchLanguages(username, repo string, l *map[string]int) {
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

	if err = json.Unmarshal(body, &l); err != nil {
		log.Fatal(err)
	}
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
	filePath := "D:\\goprojects\\gitnames.txt"
	usernames, err := setup(filePath)
	if err != nil {
		fmt.Print("error reading file:", err)
		return
	}
	var repoLangs = make([]Languages, len(usernames))
	userInd := 0
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

		repoInd := 0
		repoLangs[userInd] = make(Languages, len(repos))
		for _, repo := range repos {
			repoLangs[userInd][repoInd].langs = make(map[string]int)
			fetchLanguages(username, repo.Name, &repoLangs[userInd][repoInd].langs)
			if err != nil {
				log.Fatal(err)
			}
			repoInd++
		}
		printRepoTable(repos, repoLangs, userInd)
		userInd++
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

func printRepoTable(repos []Repository, repoLangs []Languages, userInd int) {
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"Name", "Forks", "Created", "Updated", "Languages"})
	repoInd := 0
	totalForks := 0
	for _, repo := range repos {
		totalForks += repo.Forks
		proportion := calculateProportion(repoLangs[userInd][repoInd].langs)
		ls := formatLanguages(proportion)
		t.AppendRow([]interface{}{repo.Name, repo.Forks, repo.Created[:10], repo.Updated[:10],
			ls})
		repoInd++
	}
	t.AppendRow([]interface{}{"Total forks: ", totalForks, " ------ ", " ------ ", " ------ "})
	t.Render()
}
func formatLanguages(proportion map[string]float64) string {
	var result string
	for lang, percent := range proportion {
		result += fmt.Sprintf("%s->%.2f%%|", lang, percent)
	}

	return strings.TrimSuffix(result, " || ")
}

func calculateProportion(repoLanguages map[string]int) map[string]float64 {
	totalBytes := 0
	var languages []string
	for lang, count := range repoLanguages {
		totalBytes += count
		languages = append(languages, lang)
	}

	sort.Slice(languages, func(i, j int) bool {
		return repoLanguages[languages[i]] > repoLanguages[languages[j]]
	})

	proportion := make(map[string]float64)
	numLanguages := min(2, len(languages))
	for i := 0; i < numLanguages; i++ {
		proportion[fmt.Sprintf(languages[i])] = float64(repoLanguages[languages[i]]) / float64(totalBytes) * 100
	}

	proportion["Others"] = float64(totalBytes)
	for i := 0; i < numLanguages; i++ {
		proportion["Others"] -= float64(repoLanguages[languages[i]])
	}
	proportion["Others"] /= float64(totalBytes) * 100

	return proportion
}
