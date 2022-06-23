package main

import (
	"fmt"
	"os"

	ex "github.com/markus-wa/demoinfocs-golang/v2/examples"
	demoinfocs "github.com/markus-wa/demoinfocs-golang/v2/pkg/demoinfocs"
	common "github.com/markus-wa/demoinfocs-golang/v2/pkg/demoinfocs/common"
	events "github.com/markus-wa/demoinfocs-golang/v2/pkg/demoinfocs/events"
)

// Run code: go run main.go -demo /pad/naar/demo.dem
func main() {
	f, err := os.Open(ex.DemoPathFromArgs())
	checkError(err)

	defer f.Close()

	p := demoinfocs.NewParser(f)
	defer p.Close()

	// Parse header
	header, err := p.ParseHeader()
	checkError(err)
	fmt.Println("Map:", header.MapName)

	// Register handler voor kill events
	p.RegisterEventHandler(func(e events.Kill) {
		var hs string
		if e.IsHeadshot {
			hs = " (HS)"
		}
		var wallBang string
		if e.PenetratedObjects > 0 {
			wallBang = " (WB)"
		}
		fmt.Printf("%s <%v%s%s> %s\n", formatPlayer(e.Killer), e.Weapon, hs, wallBang, formatPlayer(e.Victim))
	})

	// Register handler voor het einde van de ronde
	p.RegisterEventHandler(func(e events.RoundEnd) {
		gs := p.GameState()
		switch e.Winner {
		case common.TeamTerrorists:
			// Winner's score + 1
			fmt.Printf("Round finished: winnerSide=T  ; score=%d:%d\n", gs.TeamTerrorists().Score()+1, gs.TeamCounterTerrorists().Score())
		case common.TeamCounterTerrorists:
			fmt.Printf("Round finished: winnerSide=CT ; score=%d:%d\n", gs.TeamCounterTerrorists().Score()+1, gs.TeamTerrorists().Score())
		default:
			// Gelijk spel
			fmt.Println("Round finished: No winner (tie)")
		}
	})

	// Register handler voor chat messages
	p.RegisterEventHandler(func(e events.ChatMessage) {
		fmt.Printf("Chat - %s says: %s\n", formatPlayer(e.Sender), e.Text)
	})
	// Register handler voor een rankup van spelers in de match
	p.RegisterEventHandler(func(e events.RankUpdate) {
		fmt.Printf("Rank Update: %d went from rank %d to rank %d, change: %f\n", e.SteamID32, e.RankOld, e.RankNew, e.RankChange)
	})

	err = p.ParseToEnd()
	checkError(err)
}

func formatPlayer(p *common.Player) string {
	if p == nil {
		return "?"
	}

	switch p.Team {
	case common.TeamTerrorists:
		return "[T]" + p.Name
	case common.TeamCounterTerrorists:
		return "[CT]" + p.Name
	}

	return p.Name
}

func checkError(err error) {
	if err != nil {
		panic(err)
	}
}
