package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

var collectionUUIDs = map[string]string{}

func getColsSize() int {
	return len(collectionUUIDs)
}

func init() {
	yaml, err := os.Open("collections.yml")
	if err != nil {
		panic(err)
	}

	scanner := bufio.NewScanner(yaml)
	for scanner.Scan() {
		line := strings.Split(scanner.Text(), ": ")
		collectionUUIDs[line[0]] = line[1]
	}

	if scanner.Err() != nil {
		panic(scanner.Err())
	}
}

func getColUUID(col string) string {
	for k, v := range collectionUUIDs {
		if k == col {
			return v
		}
	}

	return ""
}

func printCols() {
	for k, v := range collectionUUIDs {
		fmt.Printf("key: %s value: %s\n", k, v)
	}
}
