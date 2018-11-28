package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"strings"
	"time"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "no home directory provided")
		os.Exit(1)
	}
	homedir := os.Args[1]
	users, err := getUsers(homedir)
	if err != nil {
		if os.IsNotExist(err) {
			os.Exit(0)
		} else {
			fmt.Fprintf(os.Stderr, "could not read authorized_github_users: %v\n", err)
			os.Exit(1)
		}
	}
	keys, err := getUserKeys(users)
	if err != nil {
		fmt.Fprintf(os.Stderr, "could not get public keys: %v\n", err)
		os.Exit(1)
	}
	for _, key := range keys {
		fmt.Println(key)
	}
	os.Exit(0)
}

func getUsers(homedir string) ([]string, error) {
	filepath := path.Join(homedir, ".ssh", "authorized_github_users")
	data, err := ioutil.ReadFile(filepath)
	if err != nil {
		return nil, err
	}
	lines := strings.Split(string(data), "\n")
	users := map[string]struct{}{}
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if len(trimmed) > 0 && !strings.HasPrefix(trimmed, "#") {
			users[trimmed] = struct{}{}
		}
	}
	uniqueUsers := []string{}
	for user := range users {
		uniqueUsers = append(uniqueUsers, user)
	}
	return uniqueUsers, nil
}

// GHKey is a key structure from the GitHub REST API
type GHKey struct {
	ID  int    `json:"id"`
	Key string `json:"key"`
}

func getUserKeys(users []string) ([]string, error) {
	client := &http.Client{Timeout: 5 * time.Second}
	keys := []string{}
	for _, user := range users {
		url := fmt.Sprintf("https://api.github.com/users/%s/keys", user)
		res, err := client.Get(url)
		if err != nil {
			return nil, err
		}
		if res.StatusCode != http.StatusOK {
			return nil, fmt.Errorf("error fetching keys: got http %d %s", res.StatusCode, res.Status)
		}
		data, err := ioutil.ReadAll(res.Body)
		defer res.Body.Close()
		if err != nil {
			return nil, err
		}
		userKeys := []GHKey{}
		if err := json.Unmarshal(data, &userKeys); err != nil {
			return nil, err
		}
		for _, userKey := range userKeys {
			keys = append(keys, userKey.Key)
		}
	}
	return keys, nil
}
