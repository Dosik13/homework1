package gitnames

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/rodaine/table"
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

	for _, word := range usernames {
		fmt.Print(word)
	}

	for _, username := range usernames {
		user, err := fetchUsersData(username)
		if err != nil {
			log.Fatal(err)
		}

		repos, err := fetchRepos(username)
		if err != nil {
			log.Fatal(err)
		}

		var totalForks int
		languages := make(map[string]int)
		var creationYears, updateYears []int

		for _, repo := range repos {
			totalForks += repo.Forks
			lang, err := fetchLanguages(username, repo.Name)
			if err != nil {
				log.Fatal(err)
			}

			for key, value := range lang {
				languages[key] += value
			}

			creationYear, _ := getYear(repo.Created)
			creationYears = append(creationYears, creationYear)

			updateYear, _ := getYear(repo.Updated)
			updateYears = append(updateYears, updateYear)
		}

		creationDistribution := calculateYearDistribution(creationYears)
		updateDistribution := calculateYearDistribution(updateYears)

		printStatisticsReport(username, user, len(repos), languages, totalForks, creationDistribution, updateDistribution)
	}

}

func getYear(dateString string) (int, error) {
	date, err := time.Parse(time.RFC3339, dateString)
	if err != nil {
		return 0, err
	}
	return date.Year(), nil
}

func calculateYearDistribution(years []int) map[int]int {
	distribution := make(map[int]int)
	for _, year := range years {
		distribution[year]++
	}
	return distribution
}
func printStatisticsReport(username string, user User, repoCount int, languages map[string]int, totalForks int, creationDistribution, updateDistribution map[int]int) {
	fmt.Printf("Statistics Report for %s\n", username)

	t := table.New("Username", "Repos", "Languages", "Followers", "Forks", "Creation Year", "Update Year", "Activity Distribution")

	languagesStr := formatLanguages(languages)
	creationDistributionStr := formatYearDistribution(creationDistribution)
	updateDistributionStr := formatYearDistribution(updateDistribution)
	activityDistributionStr := formatYearDistribution(calculateActivityDistribution(creationDistribution, updateDistribution))

	t.AddRow(username, fmt.Sprintf("%d", repoCount), languagesStr, fmt.Sprintf("%d", user.Followers), fmt.Sprintf("%d", totalForks), creationDistributionStr, updateDistributionStr, activityDistributionStr)

	fmt.Println(t)
}
func formatLanguages(languages map[string]int) string {
	var langList []string
	for lang, count := range languages {
		langList = append(langList, fmt.Sprintf("%s: %d", lang, count))
	}

	sort.Strings(langList)
	return strings.Join(langList, ", ")
}

func formatYearDistribution(distribution map[int]int) string {
	var result []string
	for year, count := range distribution {
		result = append(result, fmt.Sprintf("%d: %d", year, count))
	}
	sort.Strings(result)
	return strings.Join(result, ", ")
}
func calculateActivityDistribution(creationDistribution, updateDistribution map[int]int) map[int]int {
	activityDistribution := make(map[int]int)

	for year, count := range creationDistribution {
		activityDistribution[year] += count
	}

	for year, count := range updateDistribution {
		activityDistribution[year] += count
	}

	return activityDistribution
}
