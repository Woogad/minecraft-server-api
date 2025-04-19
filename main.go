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
	Online       bool     `json:"online"`
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

	serverPort := os.Getenv("SERVERPORT")
	if serverPort == "" {
		serverPort = ":8080"
	}
	mux := http.NewServeMux()
	mux.HandleFunc("/favicon.ico", func(w http.ResponseWriter, r *http.Request) {
		http.NotFound(w, r) // or just do nothing
	})
	mux.HandleFunc("/api/status", serverStatusHandler)
	mux.Handle("/", http.FileServer(http.Dir("./public")))

	fmt.Println("Server is running on " + serverPort)
	log.Fatal(http.ListenAndServe(serverPort, mux))
}

func serverStatusHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	minecraftAddress := os.Getenv("MINECRAFT_ADDRESS")
	response := ServerStatus{
		Host: minecraftAddress,
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	resp, err := status.Modern(ctx, minecraftAddress, 25565)
	if err != nil {
		response.Online = false
		fmt.Println("Failed to get server status: " + err.Error())
	} else {

		if len(resp.Players.Sample) > 0 {
			for _, player := range resp.Players.Sample {
				response.PlayerList = append(response.PlayerList, player.Name.Raw)
			}
		}

		response.Online = true
		response.PlayerOnline = int(*resp.Players.Online)
		response.MaxPlayers = int(*resp.Players.Max)
		response.Version = resp.Version.Name.Raw
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
