package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os"
	"runtime"
	"sort"
	"sync"
	"time"

	_ "github.com/lib/pq"
)

var (
	dbstr = flag.String("dbstr", "dbname=dfs sslmode=disable", "DB connection string")
	db    *sql.DB

	stackSize   = flag.Int("stacksize", 4, "Team must include a stack of this many players")
	concurrency = flag.Int("c", 10, "Number of parallel workers")
	maxGT       = flag.Int("maxgt", 24, "Max game time")
	minGT       = flag.Int("mingt", -12, "Min game time")
	sigChan     = make(chan os.Signal, 1)
	resultTeams = []Team{}
)

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	flag.Parse()
	initDb()
}

func initDb() {
	var (
		err error
		row string
	)

	if db, err = sql.Open("postgres", *dbstr); err != nil {
		log.Println("Could not open database connection")
		log.Fatal(err)
	}

	if err = db.QueryRow("SELECT 'OK'").Scan(&row); err != nil {
		log.Println("Could not connect to database")
		log.Println(*dbstr)
		log.Fatal(err)
	}

	db.SetMaxIdleConns(*concurrency)
	db.SetMaxOpenConns(*concurrency)
}

func main() {
	pitcher := []*Player{}
	first := []*Player{}
	second := []*Player{}
	short := []*Player{}
	third := []*Player{}
	catcher := []*Player{}
	of := []*Player{}

	for _, p := range GetAllPlayers() {
		if !(p.Time <= *maxGT && p.Time >= *minGT) {
			continue
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

	iterations := 0
	wg := &sync.WaitGroup{}
	for i := 0; i < *concurrency; i++ {
		wg.Add(1)
		go func() {
			for {
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

				outfielders := []*Player{of[ofIdx1], of[ofIdx2], of[ofIdx3]}
				sort.Sort(ByPrice(outfielders))
				t := Team{
					Pitcher: pitcher[rand.Intn(len(pitcher))],
					First:   first[rand.Intn(len(first))],
					Second:  second[rand.Intn(len(second))],
					Short:   short[rand.Intn(len(short))],
					Third:   third[rand.Intn(len(third))],
					Catcher: catcher[rand.Intn(len(catcher))],
					OF1:     outfielders[0],
					OF2:     outfielders[1],
					OF3:     outfielders[2],
				}

				teamSalary := t.salary()
				if teamSalary <= 35000 && teamSalary >= 34000 && t.Valid() {
					err := t.save()
					if err != nil {
						log.Println(err)
					}
				}
				iterations += 1
				if iterations%10000000 == 0 {
					log.Println(fmt.Sprintf("Completed 1MM iterations [%vB total]", (float64(iterations) / float64(1000000000))))
				}
			}
		}()
	}
	wg.Wait()
}
