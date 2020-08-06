package dataprocessing

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"retargetly-exercise/models"
	"strconv"
	"strings"
)

func StartProcess(fileName string) {
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

	// get a map with the segments as key and the countries as array of strings
	for scanner.Scan() {
		preProcessLinesBySeg(scanner.Text(), &countArr)
	}

	// Make the count per country according to the segment. Then store this value in a map with the segment as the key 
	for k, v := range countArr {
		key, err := strconv.Atoi(k)
		if err != nil {
			fmt.Println("Error")
			return
		}
		finalDataMap[key] = countPerSegment(v)
	}
	
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
//	fmt.Printf("%v", apiStruct)
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
