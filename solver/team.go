package main

import (
	"errors"
	"fmt"
	"log"
)

type Team struct {
	Pitcher *Player
	Catcher *Player
	First   *Player
	Second  *Player
	Third   *Player
	Short   *Player
	OF1     *Player
	OF2     *Player
	OF3     *Player
	Points  float64
}

func (t *Team) save() error {
	if t.Pitcher == nil {
		return errors.New("Missing pitcher")
	}
	if t.Catcher == nil {
		return errors.New("Missing catcher")
	}
	if t.First == nil {
		return errors.New("Missing first")
	}
	if t.Second == nil {
		return errors.New("Missing second")
	}
	if t.Third == nil {
		return errors.New("Missing third")
	}
	if t.Short == nil {
		return errors.New("Missing short")
	}
	if t.OF1 == nil {
		return errors.New("Missing of1")
	}
	if t.OF2 == nil {
		return errors.New("Missing of2")
	}
	if t.OF3 == nil {
		return errors.New("Missing of3")
	}

	_, err := db.Exec(fmt.Sprintf(`INSERT INTO teams
		(pitcher, catcher, first, second, third, short, of1, of2, of3)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`),
		t.Pitcher.Id, t.Catcher.Id, t.First.Id, t.Second.Id, t.Third.Id, t.Short.Id, t.OF1.Id, t.OF2.Id, t.OF3.Id)

	if err != nil {
		log.Println("Error saving team: ", err)
		return err
	}

	return nil
}

func (t *Team) salary() int {
	var totalSalary int
	for _, p := range []*Player{t.Pitcher, t.First, t.Second, t.Short, t.Third, t.Catcher, t.OF1, t.OF2, t.OF3} {
		totalSalary += p.Price
	}

	return totalSalary
}

func (t *Team) points() float64 {
	var totalPoints float64
	for _, p := range []*Player{t.Pitcher, t.First, t.Second, t.Short, t.Third, t.Catcher, t.OF1, t.OF2, t.OF3} {
		totalPoints += p.Points
	}

	return totalPoints
}

func (t *Team) Valid() bool {
	teamCounts := make(map[string]int)
	for _, p := range []*Player{t.Pitcher, t.First, t.Second, t.Short, t.Third, t.Catcher, t.OF1, t.OF2, t.OF3} {
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
		if v == *stackSize {
			return true
		}
	}
	return false
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

func removeTeam(t []Team, idx int) []Team {
	t = append(t[:idx], t[idx+1:]...)
	return t
}
