package main

import (
	"bufio"
	"flag"
	"fmt"
	"io/ioutil"
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

var partners = map[string]string{
	"2": "tamwag",
	"3": "fales",
	"6": "nyarchives",
}

var workOrderPtr = flag.String("workorder", "digitization_work_order_report.tsv", "the location of the work order")

func main() {
	flag.Parse()

	//open the workorder as a slice of WorkOrderEntries
	workOrder, err := openWorkOrder(*workOrderPtr)
	if err != nil {
		panic(err)
	}

	//get the name of the partner and collection directory from the first line of the workorder
	partnerId := strings.Split(workOrder[1].URI, "/")[2]
	collectionPrefix := strings.Split(workOrder[1].ResourceID, ".")[0]
	collectionNum := strings.Split(workOrder[1].ResourceID, ".")[1]
	directoryName := collectionPrefix + "-" + collectionNum

	//create the root directory
	err = os.Mkdir(directoryName, 0755)
	if err != nil {
		panic(err)
	}

	//create the metadata directory
	metadataDir := filepath.Join(directoryName, "metadata")
	err = os.Mkdir(metadataDir, 0755)
	if err != nil {
		panic(err)
	}

	//copy the work order to the metadata directory
	err = CopyWorkOrder(*workOrderPtr, metadataDir)
	if err != nil {
		panic(err)
	}

	//create the transfer-info.txt file
	err = CreateTransferInfoFile(metadataDir, partnerId, strings.ToLower(collectionPrefix),
		strings.Replace(collectionNum, "0", "", 1))
	if err != nil {
		panic(err)
	}

	//create cuid directories
	for _, entry := range workOrder {
		subdir := filepath.Join(directoryName, entry.ComponentID)
		err := os.Mkdir(subdir, 0755)
		if err != nil {
			panic(err)
		}
	}
}

func CreateTransferInfoFile(metadataDir string, partnerId string, collectionPrefix string, collectionNum string) error {
	transferInfoFileLoc := filepath.Join(metadataDir, "transfer-info.txt")
	transferInfoFile, err := os.Create(transferInfoFileLoc)
	if err != nil {
		return err
	}
	defer transferInfoFile.Close()
	writer := bufio.NewWriter(transferInfoFile)
	writer.WriteString("Internal-sender-identifier: " + partners[partnerId] + "/" + collectionPrefix + collectionNum + "\n")
	writer.WriteString(transferInfo)
	writer.WriteString("nyu-dl-project-name: " + partners[partnerId] + "/" + collectionPrefix + collectionNum + "\n")
	writer.Flush()
	colId := strings.TrimSpace(collectionPrefix + collectionNum)
	uuid := getColUUID(colId)
	fmt.Println(uuid)
	if uuid == "" {
		fmt.Println("WARNING: no UUID in configuration file for: ", colId)
	} else {
		writer.WriteString("nyu-dl-rstar-collection-id: " + uuid)
		writer.Flush()
	}
	return nil
}

func CopyWorkOrder(workorder string, mdLocation string) error {
	wo, err := ioutil.ReadFile(workorder)
	if err != nil {
		return err
	}

	wo2, err := os.Create(filepath.Join(mdLocation, "digitization_work_order_report.tsv"))
	if err != nil {
		return err
	}

	writer := bufio.NewWriter(wo2)
	writer.Write(wo)
	writer.Flush()
	wo2.Close()

	return nil
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

var transferInfo = `Source-organization: ACM
Organization-address: 70 Washington Square South, New York, NY 10012
Contact-name: 
Contact-phone: 
Contact-email: 
nyu-dl-content-classification: processed_collection
nyu-dl-package-type: AIP
`

//Internal-sender-identifier: fales/mss413
//nyu-dl-project-name: fales/mss413
//nyu-dl-archivesspace-resource-url: https://archivesspace.library.nyu.edu:8089/repositories/3/resources/1883
