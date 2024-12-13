package main

import (
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/voltageeee/anilibria-go-wrapper/anilibria"
)

func main() {
	title, err := anilibria.Search("токийский гуль", []string{})
	if err != nil {
		panic(err)
	}

	fmt.Println("Title name:", title[0].Names.Ru)
	fmt.Println("Title ID:", title[0].Id)
	fmt.Println("Title franchise:", title[0].Franchises[0].Franchise.Name)

	titleId := title[0].Id

	titleFranchse, err := anilibria.GetFranchises([]string{}, titleId)
	if err != nil {
		panic(err)
	}

	fmt.Println("Franchise releases:", titleFranchse[0].Releases[0].Names.Ru)

	titlePoster := title[0].Posters.Original.URL

	poster, err := http.Get(fmt.Sprintf("https://anilibria.tv/%s", titlePoster))
	if err != nil {
		panic(err)
	}
	defer poster.Body.Close()

	posterFile, err := os.Create("poster.jpg")
	if err != nil {
		panic(err)
	}
	defer posterFile.Close()

	io.Copy(posterFile, poster.Body)

	fmt.Println("Poster saved to poster.jpg!")

	titleTorrent := title[0].Torrents.List[0].URL

	torrent, err := http.Get(fmt.Sprintf("https://anilibria.tv/%s", titleTorrent))
	if err != nil {
		panic(err)
	}
	defer torrent.Body.Close()

	torrentFile, err := os.Create("torrent.torrent")
	if err != nil {
		panic(err)
	}
	defer torrentFile.Close()

	io.Copy(torrentFile, torrent.Body)

	fmt.Println("Torrent saved to torrent.torrent!")
}
