package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os"
	"runtime"
	"strconv"
	"strings"
	"time"
)

var (
	AWSAccessKey = "AKIAIYDAIEFMMXKCOB3Q"
	AWSSecretKey = "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"

	cassandra   = flag.String("cassandra", "localhost", "Cassandra servers")
	concurrency = flag.Int("c", 10, "Concurrency")
	delimiter   = flag.String("delimiter", ",", "Field delimiter in input file")
	file        = flag.String("file", "", "File with records to delete (Format: emd5,publisher,network,offer)")
	iterations  = flag.Int("iterations", 100000000, "Random iterations to run")
	numTeams    = flag.Int("n", 10, "Produce the top n teams")
)

type Player struct {
	Name   string
	Team   string
	VsTeam string
	Pos    string
	Points float64
	Price  int
}

type Team struct {
	Pitcher Player
	First   Player
	Second  Player
	Short   Player
	Third   Player
	Catcher Player
	OF1     Player
	OF2     Player
	OF3     Player
	Points  float64
}

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	flag.Parse()

	if *file == "" {
		log.Fatal("Must provide a file")
	}
}

func main() {
	pitcher := []Player{}
	first := []Player{}
	second := []Player{}
	short := []Player{}
	third := []Player{}
	catcher := []Player{}
	of := []Player{}

	var scanner *bufio.Scanner

	fh, err := os.Open(*file)
	if err != nil {
		log.Fatal("Unable to open file: ", err)
	}
	defer fh.Close()

	scanner = bufio.NewScanner(fh)

	for scanner.Scan() {
		row := scanner.Text()
		cols := strings.Split(row, *delimiter)

		pts, err := strconv.ParseFloat(cols[4], 64)
		if err != nil {
			log.Fatal(err)
		}

		price, err := strconv.Atoi(cols[5])
		if err != nil {
			log.Fatal(err)
		}

		p := Player{
			Name:   cols[0],
			Team:   cols[1],
			Pos:    cols[2],
			VsTeam: strings.Replace(cols[3], "@", "", -1),
			Points: pts,
			Price:  price,
		}
		switch p.Pos {
		case "1B":
			first = append(first, p)
		case "2B":
			second = append(second, p)
		case "3B":
			third = append(third, p)
		case "SS":
			short = append(short, p)
		case "C":
			catcher = append(catcher, p)
		case "OF":
			of = append(of, p)
		case "P":
			pitcher = append(pitcher, p)
		default:
			log.Println("Unrecognized pos: ", p)
			continue
		}
	}

	resultTeams := []Team{}
	minScore := 0.0
	for i := 0; i < *iterations; i++ {
		ofIndexesUsed := make(map[int]bool)
		rand.Seed(time.Now().UTC().UnixNano())

		var ofIdx1, ofIdx2, ofIdx3 int

		ofIdx1 = rand.Intn(len(of))
		ofIndexesUsed[ofIdx1] = true

		for {
			ofIdx2 = rand.Intn(len(of))
			if _, ok := ofIndexesUsed[ofIdx2]; ok {
				continue
			}
			ofIndexesUsed[ofIdx2] = true
			break
		}

		for {
			ofIdx3 = rand.Intn(len(of))
			if _, ok := ofIndexesUsed[ofIdx3]; ok {
				continue
			}
			break
		}

		t := Team{
			Pitcher: pitcher[rand.Intn(len(pitcher))],
			First:   first[rand.Intn(len(first))],
			Second:  second[rand.Intn(len(second))],
			Short:   short[rand.Intn(len(short))],
			Third:   third[rand.Intn(len(third))],
			Catcher: catcher[rand.Intn(len(catcher))],
			OF1:     of[ofIdx1],
			OF2:     of[ofIdx2],
			OF3:     of[ofIdx3],
		}

		teamSalary := t.salary()
		t.Points = t.points()
		if teamSalary <= 35000 {
			//eligibleTeams = append(eligibleTeams, t)
			if t.Points > minScore && t.ValidTeam() {
				if len(resultTeams) >= *numTeams {
					idx, min := getIndexAndScoreWorstTeam(resultTeams)
					minScore = min
					resultTeams = removeTeam(resultTeams, idx)
				}
				resultTeams = append(resultTeams, t)
			}
		}

		if i%1000000 == 0 {
			log.Println(fmt.Sprintf("Completed 10000 iterations [%v percent complete]", (float64(i)/float64(*iterations))*100.0))
		}
	}
	for _, t := range resultTeams {
		log.Println(t)
		log.Println("\n")
	}
}

func removeTeam(t []Team, idx int) []Team {
	t = append(t[:idx], t[idx+1:]...)
	return t
}

func (t *Team) salary() int {
	var totalSalary int
	for _, p := range []Player{t.Pitcher, t.First, t.Second, t.Short, t.Third, t.Catcher, t.OF1, t.OF2, t.OF3} {
		totalSalary += p.Price
	}

	return totalSalary
}

func (t *Team) points() float64 {
	var totalPoints float64
	for _, p := range []Player{t.Pitcher, t.First, t.Second, t.Short, t.Third, t.Catcher, t.OF1, t.OF2, t.OF3} {
		totalPoints += p.Points
	}

	return totalPoints
}

func (t *Team) ValidTeam() bool {
	teamCounts := make(map[string]int)
	for _, p := range []Player{t.Pitcher, t.First, t.Second, t.Short, t.Third, t.Catcher, t.OF1, t.OF2, t.OF3} {
		// facing my starting pitcher
		if p.Team == t.Pitcher.VsTeam {
			return false
		}
		if _, ok := teamCounts[p.Team]; !ok {
			teamCounts[p.Team] = 1
		} else {
			teamCounts[p.Team] += 1
		}
	}

	for _, v := range teamCounts {
		if v > 4 {
			return false
		}
	}
	return true
}

func getIndexAndScoreWorstTeam(teams []Team) (int, float64) {
	var idx int
	score := -1.0
	for i, t := range teams {
		if score == -1.0 || t.Points < score {
			idx = i
			score = t.Points
		}
	}
	return idx, score
}
