package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/mcstatus-io/mcutil/v4/status"
)

type ServerStatus struct {
	Online       string   `json:"online"`
	Host         string   `json:"host"`
	Version      string   `json:"version"`
	MaxPlayers   int      `json:"max_players"`
	PlayerOnline int      `json:"player_count"`
	PlayerList   []string `json:"player_list"`
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	http.HandleFunc("/favicon.ico", func(w http.ResponseWriter, r *http.Request) {
		http.NotFound(w, r) // or just do nothing
	})
	http.HandleFunc("/", serverStatusHandler)

	fmt.Println("Server is running on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func serverStatusHandler(w http.ResponseWriter, r *http.Request) {
	address := os.Getenv("ADDRESS")
	response := ServerStatus{
		Host: address,
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	resp, err := status.Modern(ctx, address, 25565)
	if err != nil {
		response.Online = "No"
		fmt.Println("Failed to get server status: " + err.Error())
	} else {

		if len(resp.Players.Sample) > 0 {
			for _, player := range resp.Players.Sample {
				response.PlayerList = append(response.PlayerList, player.Name.Raw)
			}
		}

		response.Online = "Yes"
		response.PlayerOnline = int(*resp.Players.Online)
		response.MaxPlayers = int(*resp.Players.Max)
		response.Version = resp.Version.Name.Raw
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
