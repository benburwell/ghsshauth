package main

import (
	"fmt"
	"github.com/benburwell/ghsshauth"
	"os"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "no home directory provided")
		os.Exit(1)
	}
	homedir := os.Args[1]
	users, err := ghsshauth.ReadAuthorizedGithubUsers(homedir)
	if err != nil {
		if os.IsNotExist(err) {
			os.Exit(0)
		} else {
			fmt.Fprintf(os.Stderr, "could not read authorized_github_users: %v\n", err)
			os.Exit(1)
		}
	}
	keys, err := ghsshauth.FetchUserKeys(users...)
	if err != nil {
		fmt.Fprintf(os.Stderr, "could not get public keys: %v\n", err)
		os.Exit(1)
	}
	for _, key := range keys {
		fmt.Println(key)
	}
	os.Exit(0)
}
