package main

import (
	"encoding/json"
	"io/ioutil"
)

type user struct {
	Sid  string `json:"sid" binding:"required"`
	Name string `json:"name" binding:"required"`
}

func getUserMap() map[string]user {
	data, _ := ioutil.ReadFile("./data.json")
	m := make(map[string]user)
	json.Unmarshal([]byte(data), &m)
	return m
}
