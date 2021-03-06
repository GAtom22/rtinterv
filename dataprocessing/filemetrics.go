package dataprocessing

import (
	"bufio"
	"fmt"
	"os"
	m "retargetly-exercise/models"
	"strconv"
	"strings"
	"time"
)

//GetFileMetrics returns the metrics required by "fileName" parameter
func GetFileMetrics(fileName string) (int64, []m.Segment, error) {
	segmentNum := 238		//Number of segments
	countArr := make(map[string][]string, segmentNum)
	finalDataMap := make(map[int]map[string]int, segmentNum)
	apiStruct := make([]m.Segment, 0, segmentNum)
	
	file, err := os.Open(fileName)
	if err != nil {
		return 0, apiStruct, fmt.Errorf("Failed to open file %s", fileName)
	}
	scanner := bufio.NewScanner(file)

	// Close when the function returns
	defer file.Close()



	// get a map with the segments as key and the countries as array of strings in a 100000-line batch
	var lineNum int
	for scanner.Scan() {
		if lineNum%100000 == 0 {
			// take the map that has segment:["Countries array"] and convert to a map for each segment with country as key and user count as value
			err := countPerCountryPerSeg(&countArr, &finalDataMap)
			if err != nil {
				return 0, apiStruct, fmt.Errorf("Error while processing data (cannot convert key to int)")
			}
			//Reset the count
			countArr = make(map[string][]string, segmentNum)
		}
		preProcessLinesBySeg(scanner.Text(), &countArr)
		lineNum++
	}

	// Make the count per country according to the segment. Then store this value in a map with the segment as the key
	err = countPerCountryPerSeg(&countArr, &finalDataMap)
	if err != nil {
		return 0, apiStruct, fmt.Errorf("Error while processing data (cannot convert key to int)")
	}

	//Parse the map with data to the api response structure
	parseToAPIResponse(&finalDataMap, &apiStruct)

	return time.Now().Unix(), apiStruct, nil
}

func parseToAPIResponse(dataMap *map[int]map[string]int, apiStruct *[]m.Segment) {
	for segmt, countries := range *dataMap {
		// create uniques by country
		uniquesArr := []m.Unique{}
		for countryID, count := range countries {
			newUnique := m.Unique{
				Country: countryID,
				Count:   count,
			}
			uniquesArr = append(uniquesArr, newUnique)
		}
		newSegment := m.Segment{
			SegmentID: segmt,
			Uniques:   uniquesArr,
		}
		(*apiStruct) = append((*apiStruct), newSegment)
	}
}

//preProcessLinesBySeg creates a map with the segment as keys and an array of country codes as values
func preProcessLinesBySeg(line string, countryArr *map[string][]string) {
	record := strings.Split(line, "\t")
	segments := strings.Split(record[1], ",")
	for _, s := range segments {
		(*countryArr)[s] = append((*countryArr)[s], record[2])
	}
}

//countPerSegment input is an array of countries and creates a map with the count per country
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

func countPerCountryPerSeg(data *map[string][]string, dataMap *map[int]map[string]int) error {
	for k, v := range *data {
		key, err := strconv.Atoi(k)
		if err != nil {
			return err
		}
		newCountriesCount := countPerSegment(v)
		if countries, exist := (*dataMap)[key]; exist {
			// Add the countries count to the existing object
			(*dataMap)[key] = updateDataMapObject(countries, newCountriesCount)
			continue
		}
		(*dataMap)[key] = newCountriesCount
	}
	return nil
}

func updateDataMapObject(currentCountMap map[string]int, newCountMap map[string]int) map[string]int {
	for country, count := range newCountMap {
		// check if the item/element exist in the duplicate_frequency map
		_, exist := currentCountMap[country]
		if exist {
			currentCountMap[country] += count // increase counter if already in the map
		} else {
			currentCountMap[country] = count // else start counting
		}
	}
	return currentCountMap
}
