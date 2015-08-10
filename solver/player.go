package main

import (
	"log"
)

type Player struct {
	Id     string
	Name   string
	Team   string
	VsTeam string
	Pos    string
	Points float64
	Price  int
}

// ByPrice implements sort.Interface for []Player based on
// the Price field.
type ByPrice []*Player

func (a ByPrice) Len() int           { return len(a) }
func (a ByPrice) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByPrice) Less(i, j int) bool { return a[i].Price < a[j].Price }

func GetAllPlayers() []*Player {
	var q string
	q = "SELECT id, name, team, vs_team, pos, points, price FROM players"
	rows, err := db.Query(q)
	if err != nil {
		log.Println("Error getting players from the DB", err)
		return []*Player{}
	}

	players := make([]*Player, 0)
	defer rows.Close()
	for rows.Next() {
		var p Player
		if err := rows.Scan(&p.Id, &p.Name, &p.Team, &p.VsTeam, &p.Pos, &p.Points, &p.Price); err != nil {
			log.Println("Error getting players from the DB")
			return []*Player{}
		}
		players = append(players, &p)
	}

	return players
}
