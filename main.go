package main

import (
	"bufio"
	"flag"
	"fmt"
	"github.com/nyudlts/erecs-directory-generator/collections"
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

var partners = map[string][]string{
	"2": []string{"tamwag", "TW"},
	"3": []string{"fales", "FA"},
	"6": []string {"nyarchives", "UA"},
}


var workOrderPtr = flag.String("workorder", "", "the location of the work order")
var resourceIdPtr = flag.Int("resourceId", 0, "the aspace id of the resource")
var rwsPtr = flag.String("rwsLocation", "", "the location of ACM-BornDigital")
var batchPtr = flag.Int("batchNum", 0, "the batch number")
var rwsLocation string

func main() {
	flag.Parse()
	setup()

	//open the workorder as a slice of WorkOrderEntries
	workOrder, err := openWorkOrder(*workOrderPtr)
	if err != nil {
		panic(err)
	}
	fmt.Println(workOrder)
	//get the name of the partner and collection directory from the first line of the workorder
	partnerId := strings.Split(workOrder[0].URI, "/")[2]
	collectionPrefix := strings.Split(workOrder[0].ResourceID, ".")[0]
	collectionNum := strings.Split(workOrder[0].ResourceID, ".")[1]
	directoryName := filepath.Join(rwsLocation, partners[partnerId][1] + "_" + collectionPrefix + "_" + collectionNum)

	if *batchPtr > 0 {
		directoryName = fmt.Sprintf("%s-Batch-%d", directoryName, *batchPtr)
	}


	err = createDirectories(directoryName); if err != nil {
		panic(err)
	}

	metadataDirectory := filepath.Join(directoryName, "metadata")

	err = CopyWorkOrder(*workOrderPtr, metadataDirectory); if err != nil {
		panic(err)
	}

	//create the transfer-info.txt file
	err = CreateTransferInfoFile(metadataDirectory, partnerId, strings.ToLower(collectionPrefix), collectionNum)
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

	fmt.Println("complete")

}

func setup() {
	//check that Archivematica-Stagin is available
	rwsLocation = filepath.Join(*rwsPtr, "Archivematica-Staging")
	if _, err := os.Stat(rwsLocation); os.IsNotExist(err) {
		panic(fmt.Errorf("Archivematica-Staging does not exist at location %s", *rwsPtr))
	}

	//check that the work order exists
	if _, err := os.Stat(*workOrderPtr); os.IsNotExist(err) {
		panic(fmt.Errorf("Work Order does not exist at %s", *workOrderPtr))
	}

	//check that there is a resource ID and it is not 0
	if *resourceIdPtr == 0 {
		panic(fmt.Errorf("Resource ID must be defined and not equal to zero"))
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

func createDirectories(directoryName string) error {
	//create the root directory
	err := os.Mkdir(directoryName, 0755)
	if err != nil {
		return(err)
	}

	//create the metadata directory
	metadataDir := filepath.Join(directoryName, "metadata")
	err = os.Mkdir(metadataDir, 0755)
	if err != nil {
		return(err)
	}

	return nil
}

func CopyWorkOrder(workorder string, mdLocation string) error {
	wo, err := ioutil.ReadFile(workorder)
	if err != nil {
		return err
	}

	wo2, err := os.Create(filepath.Join(mdLocation, workorder))
	if err != nil {
		return err
	}

	writer := bufio.NewWriter(wo2)
	writer.Write(wo)
	writer.Flush()
	wo2.Close()

	return nil
}

func CreateTransferInfoFile(metadataDir string, partnerId string, collectionPrefix string, collectionNum string) error {
	partner := partners[partnerId]
	code := strings.TrimSpace(collectionPrefix + collectionNum)
	uuid := collections.GetUUID(partner[0], code)
	transferInfoFileLoc := filepath.Join(metadataDir, "transfer-info.txt")
	transferInfoFile, err := os.Create(transferInfoFileLoc)
	if err != nil {
		return err
	}
	defer transferInfoFile.Close()
	writer := bufio.NewWriter(transferInfoFile)
	writer.WriteString("Internal-sender-identifier: " + partner[0] + "/" + code + "\n")
	writer.WriteString(transferInfo)
	writer.WriteString("nyu-dl-project-name: " + partner[0] + "/" + code + "\n")
	writer.WriteString("nyu-dl-rstar-collection-id: " + uuid + "\n")
	writer.WriteString(fmt.Sprintf("nyu-dl-archivesspace-resource-url: https://archivesspace.library.nyu.edu:8489/repositories/%s/resources/%d", partnerId, *resourceIdPtr))
	writer.Flush()
	return nil
}

var transferInfo = `Source-organization: ACM
Organization-address: 70 Washington Square South, New York, NY 10012
Contact-name:
Contact-phone:
Contact-email:
nyu-dl-content-classification: processed_collection
nyu-dl-content-type: electronic_records	
nyu-dl-package-type: AIP
nyu-dl-hostname:
nyu-dl-pathname:
`

