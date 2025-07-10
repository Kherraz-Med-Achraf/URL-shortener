package main

import (
	"encoding/json"
	"os"
	"path/filepath"

)

type User struct {
	Username     string `json:"username"`
	PasswordHash string `json:"password_hash"`
	IsAdmin      bool   `json:"is_admin"`
}

const userDir = "data/users"

func saveUser(u User) error {
	os.MkdirAll(userDir, 0755)
	f, err := os.Create(filepath.Join(userDir, u.Username+".json"))
	if err != nil {
		return err
	}
	defer f.Close()
	return json.NewEncoder(f).Encode(u)
}

func loadUser(username string) (User, error) {
	var u User
	f, err := os.Open(filepath.Join(userDir, username+".json"))
	if err != nil {
		return u, err
	}
	defer f.Close()
	err = json.NewDecoder(f).Decode(&u)
	return u, err
}
