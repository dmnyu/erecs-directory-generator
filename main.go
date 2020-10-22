package main

import (
	"bufio"
	"os"
	"path/filepath"
	"strings"
)

type WorkOrderEntry struct {
	ResourceID          string
	RefId               string
	URI                 string
	ContainerIndicator1 string
	ContainerIndicator2 string
	ContainerIndicator3 string
	Title               string
	ComponentID         string
}

func main() {
	workOrderLoc := "digitization_work_order_report.tsv"
	workOrder, err := openWorkOrder(workOrderLoc)
	if err != nil {
		panic(err)
	}

	directoryName := strings.Replace(workOrder[1].ResourceID, ".", "-", 1)

	err = os.Mkdir(directoryName, 0755)
	if err != nil {
		panic(err)
	}

	metadataDir := filepath.Join(directoryName, "metadata")
	err = os.Mkdir(metadataDir, 0755)
	if err != nil {
		panic(err)
	}

	for _, entry := range workOrder {
		subdir := filepath.Join(directoryName, entry.ComponentID)
		err := os.Mkdir(subdir, 0755)
		if err != nil {
			panic(err)
		}
	}
}

func openWorkOrder(fileLoc string) ([]WorkOrderEntry, error) {
	var workOrder = []WorkOrderEntry{}
	workOrderFile, err := os.Open(fileLoc)
	if err != nil {
		return workOrder, err
	}

	scanner := bufio.NewScanner(workOrderFile)
	scanner.Scan() // skip the header

	for scanner.Scan() {
		line := strings.Split(scanner.Text(), "\t")
		workOrder = append(workOrder, WorkOrderEntry{
			line[0], line[1], line[2], line[3], line[4], line[5], line[6], line[7],
		})
	}

	if scanner.Err() != nil {
		return workOrder, err
	}

	return workOrder, nil
}
