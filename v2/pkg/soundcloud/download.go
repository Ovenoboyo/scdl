package soundcloud

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/Ovenoboyo/scdl/v2/pkg/mp3"
	"github.com/antchfx/htmlquery"
)

// Download queries the SoundCloud api and receives a m3u8 file, then binds the segments received into a .mp3 file
func (s *Soundcloud) Download(url string) {

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Println(err)
	}

	// set Non Hacker User Agent
	req.Header.Set("Accept", s.UserAgent)

	resp, err := s.Client.Do(req)
	if err != nil {
		log.Println(err)
	}

	// parse html
	doc, err := htmlquery.Parse(resp.Body)
	if err != nil {
		log.Println(err)
	}

	streamURL, err := s.ConstructStreamURL(doc, StreamTypeHLS)
	if err != nil {
		log.Println(err)
	}

	songName, err := s.GetTitle(doc)
	if err != nil {
		log.Println(err)
	}

	artwork, err := s.GetArtwork(doc)
	if err != nil {
		log.Println(err)
	}

	println("Got stream url", streamURL)
	// Get the response from the URL
	streamResp, err := http.Get(streamURL)
	if err != nil {
		log.Println(err)
		return
	}

	defer streamResp.Body.Close()

	// Read the body of the response
	body, err := ioutil.ReadAll(streamResp.Body)
	if err != nil {
		log.Println("Error reading response body:", err)
		return
	}

	// Unmarshal the JSON into the struct
	var audioResp AudioLink
	err = json.Unmarshal(body, &audioResp)
	if err != nil {
		fmt.Println("Error unmarshalling JSON:", err)
		return
	}

	// merge segments
	mp3.Merge(audioResp.URL, songName)

	artworkResp, err := http.Get(artwork)
	image, err := ioutil.ReadAll(artworkResp.Body)
	if err != nil {
		log.Println(err)
	}

	// set cover image for mp3 file
	mp3.SetCoverImage(songName+".mp3", image)
}
