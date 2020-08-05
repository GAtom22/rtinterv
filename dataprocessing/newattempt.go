package dataprocessing

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

type Countries struct{
	Name string
	Count int
}

func StartProcess(fileName string) {
	file, err := os.Open(fileName)
	if err != nil {
		log.Fatal(err)
	}

	// Close when the functin returns
	defer file.Close()

	scanner := bufio.NewScanner(file)
	userArr := map[string][]Countries{}

	for scanner.Scan() {
		// fmt.Printf("%v", scanner.Text())
		userArr = preProcessLines(scanner.Text(), userArr)
	}
	fmt.Printf("%v",userArr)
}

// func preProcessLines(line string, userArr []string) []string{

// 	record := strings.Split(line, "\t")
// 	segments := strings.Split(record[1],",")
// 	for _,v := range segments{
// 		userArr = append(userArr, v+record[2])
// 	}	
// 	return userArr
// }

// func preProcessLines(line string, userArr map[string]int) map[string]int{
// 	record := strings.Split(line, "\t")
// 	segments := strings.Split(record[1],",")
// 	for _,v := range segments{
// 		key := v+record[2]
// 		if _, ok:=userArr[key]; ok{
// 			userArr[key]++
// 			continue
// 		}
// 		userArr[key] = 1
// 	}	
// 	return userArr
// }

func preProcessLines(line string, userArr map[string][]Countries) map[string][]Countries{
	record := strings.Split(line, "\t")
	segments := strings.Split(record[1],",")
	for _,s := range segments{
		if _, ok:=userArr[s]; ok{
			for _, v := range userArr[s]{
				if v.Name == record[2]{
					v.Count++
					continue
				}
			}
		}
		newCountry := Countries{
			Name: record[2],
			Count: 1,
		}
		userArr[s] = append(userArr[s], newCountry)
	}	
	return userArr
}