// Package ghsshauth implements functionality necessary to use the GitHub API
// as an authorization provider for the OpenSSH server's AuthorizedKeysCommand.
package ghsshauth

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"path"
	"strings"
	"time"
)

var httpClient = &http.Client{Timeout: 5 * time.Second}

// ReadAuthorizedGithubUsers attempts to read the authorized_github_users file
// for the given home directory and returns the usernames it finds as a slice
// of strings.
func ReadAuthorizedGithubUsers(homedir string) ([]string, error) {
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

// FetchUserKeys fetches the SSH keys for the specified username(s) from the
// GitHub API
func FetchUserKeys(users ...string) ([]string, error) {
	keys := []string{}
	for _, user := range users {
		url := fmt.Sprintf("https://api.github.com/users/%s/keys", user)
		res, err := httpClient.Get(url)
		if err != nil {
			return nil, err
		}
		if res.StatusCode != http.StatusOK {
			return nil, fmt.Errorf("http error: status %d %s", res.StatusCode, res.Status)
		}
		data, err := ioutil.ReadAll(res.Body)
		defer res.Body.Close()
		if err != nil {
			return nil, err
		}
		userKeys := []struct {
			ID  int    `json:"id"`
			Key string `json:"key"`
		}{}
		if err := json.Unmarshal(data, &userKeys); err != nil {
			return nil, err
		}
		for _, userKey := range userKeys {
			keys = append(keys, userKey.Key)
		}
	}
	return keys, nil
}
