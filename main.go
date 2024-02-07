package main

import (
	"fmt"
	"math/rand"
	"time"
)

type Player struct {
	ID       int
	Dice     []int
	DiceTemp []int
	Points   int
	IsActive bool
}

func (p *Player) rollDice(rn *rand.Rand) {
	for i := range p.Dice {
		p.Dice[i] = rn.Intn(6) + 1
	}
}

func (p *Player) updatePlayer() {
	if len(p.DiceTemp) > 0 {
		p.Dice = append(p.Dice, p.DiceTemp...)
		p.DiceTemp = make([]int, 0)
	}

	if len(p.Dice) == 0 {
		p.IsActive = false
	}
}

func main() {
	seed := rand.NewSource(time.Now().UnixNano())
	rn := rand.New(seed)

	var numPlayers, numDice int
	var count int

	fmt.Print("Masukkan jumlah PEMAIN: ")
	fmt.Scan(&numPlayers)
	fmt.Print("Masukkan jumlah dadu per PEMAIN: ")
	fmt.Scan(&numDice)

	players := initializePlayers(numPlayers, numDice)
	activePlayersCount := len(players)

	for activePlayersCount > 1 {
		count++
		fmt.Println("\n=====================")
		for i := range players {
			players[i].rollDice(rn)
		}

		fmt.Printf("Lemparan dadu ke %v\n", count)
		print(players)

		// Evaluasi dadu
		for i := range players {
			evaluateDice(&players[i], &players)
		}

		// Update dadu
		for i := range players {
			players[i].updatePlayer()
		}

		fmt.Println("\nSetelah Evaluasi:")
		print(players)

		activePlayersCount = countActivePlayers(players)
	}

	var winner Player

	for _, player := range players {
		if player.Points > winner.Points {
			winner = player
		}
	}

	fmt.Printf("\nPEMAIN %d memenangkan permainan dengan %d poin!\n", winner.ID, winner.Points)
}

func initializePlayers(numPlayers, numDice int) []Player {
	players := make([]Player, numPlayers)
	for i := range players {
		players[i] = Player{
			ID:       i + 1,
			Dice:     make([]int, numDice),
			DiceTemp: make([]int, 0),
			Points:   0,
			IsActive: true,
		}
	}
	return players
}

func evaluateDice(player *Player, allPlayers *[]Player) {
	i := 0
	for i < len(player.Dice) {
		die := player.Dice[i]
		flagRemove := false
		switch die {
		case 1:
			flagRemove = true
			nextPlayer := getNextPlayer(*allPlayers, player.ID-1)
			nextPlayer.DiceTemp = append(nextPlayer.DiceTemp, 1)
			fmt.Printf("\n * PEMAIN %d memberikan dadu angka 1 kepada PEMAIN %d.\n", player.ID, nextPlayer.ID)
		case 6:
			flagRemove = true
			player.Points++
			fmt.Printf("\n * PEMAIN %d mendapatkan poin!\n", player.ID)
		}

		if flagRemove {
			player.Dice = removeDie(player.Dice, i)
		} else {
			i++
		}
	}
}

func getNextPlayer(players []Player, currentPlayerIndex int) *Player {
	nextPlayerIndex := (currentPlayerIndex + 1) % len(players)
	if !players[nextPlayerIndex].IsActive {
		return getNextPlayer(players, nextPlayerIndex)
	}
	return &players[nextPlayerIndex]
}

func removeDie(arr []int, i int) []int {
	temp := make([]int, 0)
	temp = append(temp, arr[:i]...)
	return append(temp, arr[i+1:]...)
}

func countActivePlayers(players []Player) int {
	count := 0
	for i := range players {
		if players[i].IsActive {
			count++
		}
	}
	return count
}

func print(players []Player) {
	for _, player := range players {
		fmt.Printf("PEMAIN %d: Dadu = %v, Poin = %d\n", player.ID, player.Dice, player.Points)
	}
}
