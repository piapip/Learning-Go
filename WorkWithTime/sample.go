package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"
)

//History describes the day information
type History struct {
	date   time.Time
	action string
}

//TimeLine describes a chain of History
type TimeLine []History

//Split do this
func Split(r rune) bool {
	return r == '[' || r == ']' || r == ' ' || r == '-' || r == ':' || r == '#'
}

//Len is for sorting
func (p TimeLine) Len() int {
	return len(p)
}

//Less is for sorting
func (p TimeLine) Less(i, j int) bool {
	return p[i].date.Before(p[j].date)
}

//Swap is for sorting
func (p TimeLine) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}

func main() {
	file, err := os.Open("input.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	var data TimeLine
	for scanner.Scan() {
		var log History
		instruction := scanner.Text()
		split := strings.FieldsFunc(instruction, Split)
		year, _ := strconv.Atoi(split[0])
		month, _ := strconv.Atoi(split[1])
		date, _ := strconv.Atoi(split[2])
		hour, _ := strconv.Atoi(split[3])
		minute, _ := strconv.Atoi(split[4])
		log.date = time.Date(year, time.Month(month), date, hour, minute, 0, 0, time.UTC)
		action := split[len(split)-2] + " " + split[len(split)-1]
		if len(split) == 9 {
			action = split[len(split)-3] + " " + action
		}
		log.action = action
		data = append(data, log)
	}

	sort.Sort(data)
	var id int
	var sleepStart time.Time
	var sleepEnd time.Time
	sleep := make(map[int]float64)
	for _, log := range data {
		if strings.Contains(log.action, "begins shift") {
			id, _ = strconv.Atoi(strings.Split(log.action, " ")[0])
		} else if strings.Contains(log.action, "falls asleep") {
			sleepStart = log.date
		} else {
			sleepEnd = log.date
			sleep[id] = sleep[id] + sleepEnd.Sub(sleepStart).Minutes()
		}
	}

	for id, sleepTime := range sleep {
		sleepDuration := fmt.Sprintf("%f", sleepTime)
		fmt.Println(strconv.Itoa(id) + " " + sleepDuration)
	}
}

//1777
