package dataprocessing

import (
	"sync"
	"bufio"
	"fmt"
	"log"
	"os"
	"retargetly-exercise/models"
	"strconv"
	"strings"
)

func ReadTsvFile(fileName string) {
	tsv, err := os.Open(fileName)
	if err != nil {
		fmt.Println(err)
	}

	scn := bufio.NewScanner(tsv)

	var lines []string

	for scn.Scan() {
		line := scn.Text()
		lines = append(lines, line)
	}

	if err := scn.Err(); err != nil {
		fmt.Println(err)
	}

	//lines = lines[1:] // First line is header
	out := make(map[string][]string, len(lines))

	for _, line := range lines {
		record := strings.Split(line, "\t")
		out[record[0]] = record
	}
	fmt.Printf("%v", out["0f0747d151d4bf1bb692b30de526501"][1])
}

var concurrency = 50

func StartMakingCalcs(fileName string) {
	// This channel has no buffer, so it only accepts input when something is ready
	// to take it out. This keeps the reading from getting ahead of the writers.
	workQueue := make(chan string, 10000)

	// We need to know when everyone is done so we can exit.
	complete := make(chan bool, 10000)

	defer close(complete)

	var wg sync.WaitGroup

	segments := map[string][]string{}
	// Read the lines into the work queue.
	wg.Add(1)
	go func() {
		file, err := os.Open(fileName)
		if err != nil {
			log.Fatal(err)
		}

		// Close when the functin returns
		defer file.Close()

		scanner := bufio.NewScanner(file)

		for scanner.Scan() {
			workQueue <- scanner.Text()
		}

		// Close the channel so everyone reading from it knows we're done.
		wg.Done()
		close(workQueue)
	}()

	// Now read them all off, concurrently.
	for i := 0; i < concurrency; i++ {
		wg.Add(1)
		go func() {
			startWorking2(workQueue, complete, &segments)
			wg.Done()
		}()
	}

	// Wait for everyone to finish.
	for i := 0; i < concurrency; i++ {
		<-complete
	}
	wg.Wait()  
	fmt.Printf("%v", segments)
}

func startWorking(queue chan string, complete chan bool, segments *map[int][]models.Unique) {
	var newUnique models.Unique
	for line := range queue {
		// Do the work with the line.
		record := strings.Split(line, "\t")
		segmentArr := strings.Split(record[1], ",")
		for _, v := range segmentArr {
			segNum, err := strconv.Atoi(v)
			if err != nil {
				fmt.Println("Error while converting segment to int")
			}
			if uniques, ok := (*segments)[segNum]; ok {
				//check if there is a Unique with the same country that already exist, then add 1 if exists
				for _, v := range uniques {
					if v.Country == record[2] {
						v.Count++
						continue
					}
					newUnique.Count++
					newUnique.Country = record[2]
					(*segments)[segNum] = append((*segments)[segNum], newUnique)
				}

				break
			}
			newUnique.Count++
			newUnique.Country = record[2]
			(*segments)[segNum] = append((*segments)[segNum], newUnique)
		}
	}

	// Let the main process know we're done.
	complete <- true
}

func startWorking2(queue chan string, complete chan bool, segments *map[string][]string) {
	//var newUnique models.Unique
	for line := range queue {
		// Do the work with the line.
		record := strings.Split(line, "\t")
		(*segments)[record[2]] = append((*segments)[record[2]], record[1])
	}

	// Let the main process know we're done.
	complete <- true
}


var numGoWriters = 10

func process(r string) map[string]string {
	record := strings.Split(r, "\t")
	res := map[string]string{}
	res[record[2]] = record[1]
	return res
}

func processRow(r string, ch chan<- map[string]string) {
    res := process(r)
    ch <- res
}

func writeRow(f map[string]string, ch <-chan map[string]string) {
    for k := range ch {
			for i,v := range k{
				fmt.Println(i)
				f[i] = v
			}
		}
	}

func ProcessFile(fileName string) {
    // outFile, err := os.Create("/path/to/file.out")
    // if err != nil {
    //     // handle it
    // }
		// defer outFile.Close()
		file, err := os.Open(fileName)
		if err != nil {
			log.Fatal(err)
		}
		// Close when the functin returns
		defer file.Close()

    var wg sync.WaitGroup
    ch := make(chan map[string]string, 10)  // play with this number for performance
    defer close(ch) // once we're done processing rows, we close the channel
                    // so our worker threads exit
		fScanner := bufio.NewScanner(file)
		
		result := map[string]string{}

    for fScanner.Scan() {
        wg.Add(1)
        go func() {
            processRow(fScanner.Text(), ch)
            wg.Done()
        }()
    }
    for i := 0; i < numGoWriters; i++ {
        go writeRow(result ,ch)
    }
		wg.Wait()  
		fmt.Printf("%v", result)
}