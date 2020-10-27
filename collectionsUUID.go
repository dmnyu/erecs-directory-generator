package main

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

var Repos = []Repo{}
var repoCodes = map[string]int{}

type YamlEntry struct {
	Code string `yaml:"code"`
	Id   string `yaml:"id"`
}

type Repo struct {
	Name    string      `yaml:"name"`
	Entries []YamlEntry `yaml:"entries"`
}

func init() {
	source, err := ioutil.ReadFile("collections.yml")
	if err != nil {
		panic(err)
	}
	err = yaml.Unmarshal(source, &Repos)
	if err != nil {
		panic(err)
	}

	for i, r := range Repos {
		repoCodes[r.Name] = i
	}
}

func GetEntryMapForRepo(repo string) map[string]string {
	uuids := map[string]string{}
	entries := Repos[repoCodes[repo]].Entries
	for _, entry := range entries {
		uuids[entry.Code] = entry.Id
	}
	return uuids
}

func GetUUID(repo string, code string) string {
	uuids := GetEntryMapForRepo(repo)
	return uuids[code]
}
