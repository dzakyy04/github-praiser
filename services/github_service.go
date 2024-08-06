package services

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

type GithubData struct {
	Data   interface{}
	Readme string
}

var token = os.Getenv("GITHUB_TOKEN")

func GetGithubData(username string) (GithubData, error) {
	client := &http.Client{}

	profileUrl := fmt.Sprintf("https://api.github.com/users/%s", username)
	repoUrl := fmt.Sprintf("https://api.github.com/users/%s/repos?sort=updated", username)
	readmeUrl := fmt.Sprintf("https://raw.githubusercontent.com/%s/%s/main/README.md", username, username)

	profileData, err := fetchGithubApi(client, profileUrl, token)
	if err != nil {
		return GithubData{}, err
	}

	repoData, err := fetchGithubApi(client, repoUrl, token)
	if err != nil {
		return GithubData{}, err
	}

	readme, _ := fetchReadme(client, readmeUrl, token)

	// Convert profileData to map
	profileMap, ok := profileData.(map[string]interface{})
	if !ok {
		return GithubData{}, fmt.Errorf("unexpected profile data format")
	}

	data := map[string]interface{}{
		"name":         profileMap["name"],
		"bio":          profileMap["bio"],
		"company":      profileMap["company"],
		"location":     profileMap["location"],
		"followers":    profileMap["followers"],
		"following":    profileMap["following"],
		"public_repos": profileMap["public_repos"],
		"created_at":   profileMap["created_at"],
		"updated_at":   profileMap["updated_at"],
		"repositories": repoData,
	}

	return GithubData{Data: data, Readme: readme}, nil
}

func fetchGithubApi(client *http.Client, url, token string) (interface{}, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	if token != "" {
		req.Header.Set("Authorization", "token "+token)
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("GitHub API returned status: %d", resp.StatusCode)
	}

	var data interface{}
	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func fetchReadme(client *http.Client, url, token string) (string, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", err
	}

	if token != "" {
		req.Header.Set("Authorization", "token "+token)
	}

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", nil
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}
