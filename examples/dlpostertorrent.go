package main

import (
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/voltageeee/anilibria-go-wrapper/anilibria"
)

func main() {
	// Get an array of titles that have a name of "токийский гуль"
	title, err := anilibria.Search("токийский гуль", []string{})
	if err != nil {
		panic(err)
	}

	fmt.Println("Title name:", title[0].Names.Ru)
	fmt.Println("Title ID:", title[0].Id)
	fmt.Println("Title franchise:", title[0].Franchises[0].Franchise.Name)

	titleId := title[0].Id

	// Get the title's franchises
	titleFranchse, err := anilibria.GetFranchises([]string{}, titleId)
	if err != nil {
		panic(err)
	}

	fmt.Println("Franchise releases:", titleFranchse[0].Releases[0].Names.Ru)

	titlePoster := title[0].Posters.Original.URL

	// Get the title's posters
	poster, err := http.Get(fmt.Sprintf("https://anilibria.tv/%s", titlePoster))
	if err != nil {
		panic(err)
	}
	defer poster.Body.Close()

	// Save the poster
	posterFile, err := os.Create("poster.jpg")
	if err != nil {
		panic(err)
	}
	defer posterFile.Close()

	io.Copy(posterFile, poster.Body)

	fmt.Println("Poster saved to poster.jpg!")

	titleTorrent := title[0].Torrents.List[0].URL

	// Get the torrent
	torrent, err := http.Get(fmt.Sprintf("https://anilibria.tv/%s", titleTorrent))
	if err != nil {
		panic(err)
	}
	defer torrent.Body.Close()

	// Save the torrent
	torrentFile, err := os.Create("torrent.torrent")
	if err != nil {
		panic(err)
	}
	defer torrentFile.Close()

	io.Copy(torrentFile, torrent.Body)

	fmt.Println("Torrent saved to torrent.torrent!")
}
