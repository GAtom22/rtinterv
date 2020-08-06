package dataprocessing

import (
	// "runtime"
	"bufio"
	"fmt"
	"log"
	"os"
	"retargetly-exercise/models"
	"strconv"
	"strings"
	"sync"
)

var wg = sync.WaitGroup{}
var m = sync.RWMutex{}

func StartProcesss(fileName string) {
	// runtime.GOMAXPROCS(1)
	file, err := os.Open(fileName)
	if err != nil {
		log.Fatal(err)
	}

	// Close when the functin returns
	defer file.Close()

	scanner := bufio.NewScanner(file)
	countArr := map[string][]string{}
	finalDataMap := map[int]map[string]int{}
	apiStruct := []models.Segment{}
	workQueue := make(chan map[string][]string, 5)
	secondStep := make(chan map[int]map[string]int)
	defer close(workQueue)
	defer close(secondStep)
	// get a map with the segments as key and the countries as array of strings
	var lineNum int

	for scanner.Scan() {
		if lineNum%1000000 == 0 {
			countPerCountryPerSeg(&countArr, &finalDataMap)
			countArr = map[string][]string{}
		}
		preProcessLinesBySeg(scanner.Text(), &countArr)
		lineNum++
	}

	// Make the count per country according to the segment. Then store this value in a map with the segment as the key
	countPerCountryPerSeg(&countArr, &finalDataMap)

	//Generate last data structure
	for segmt, countries := range finalDataMap {
		// create uniques by country
		uniquesArr := []models.Unique{}
		for countryID, count := range countries {
			newUnique := models.Unique{
				Country: countryID,
				Count:   count,
			}
			uniquesArr = append(uniquesArr, newUnique)
		}

		newSegment := models.Segment{
			SegmentID: segmt,
			Uniques:   uniquesArr,
		}

		apiStruct = append(apiStruct, newSegment)
	}

	fmt.Printf("ready")
	fmt.Printf("%v", apiStruct)
}

func preProcessLinesBySeg(line string, countryArr *map[string][]string) {
	record := strings.Split(line, "\t")
	segments := strings.Split(record[1], ",")
	for _, s := range segments {
		(*countryArr)[s] = append((*countryArr)[s], record[2])
	}
}

func countPerSegment(list []string) map[string]int {
	countryFrequency := make(map[string]int)
	for _, country := range list {
		// check if the item/element exist in the duplicate_frequency map
		_, exist := countryFrequency[country]
		if exist {
			countryFrequency[country]++ // increase counter by 1 if already in the map
		} else {
			countryFrequency[country] = 1 // else start counting from 1
		}
	}
	return countryFrequency
}

func countPerCountryPerSeg(data *map[string][]string, dataMap *map[int]map[string]int) {
	for k, v := range *data {
		key, err := strconv.Atoi(k)
		if err != nil {
			fmt.Println("Error")
			return
		}
		(*dataMap)[key] = countPerSegment(v)
	}
}
