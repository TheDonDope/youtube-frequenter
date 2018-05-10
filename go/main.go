package main

import (
	"fmt"
	"time"

	"github.com/TheDonDope/youtube-tinfoil-expose/go/service"
)

func main() {
	start := time.Now()
	fmt.Println("Welcome to youtube-tinfoil-expose")
	youtubeService, serviceError := service.GetYouTubeService()

	service.HandleError(serviceError, "Error creating YouTube client")

	query := service.Query{}
	query.CustomURL = "wwwKenFMde"
	query.PlaylistName = "uploads"

	results := service.Exfoliator(youtubeService, query)
	fmt.Println(results)
	elapsed := time.Since(start)
	fmt.Println(elapsed)
}
