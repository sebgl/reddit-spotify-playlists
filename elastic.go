package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

const index = "spotify"
const documentType = "playlist"

// TODO: variabilize
const elasticsearchHost = "http://localhost:9200"
const elasticsearchURL = elasticsearchHost + "/" + index + "/" + documentType

func ToElasticSearch(playlists []Playlist) error {
	for _, p := range playlists {
		pID := p.SpotifyPlaylist.ID.String()
		data, err := json.Marshal(p)
		if err != nil {
			return err
		}

		req, err := http.NewRequest("PUT", elasticsearchURL+"/"+pID, bytes.NewBuffer(data))
		if err != nil {
			return err
		}
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			return err
		}
		defer resp.Body.Close()

		if resp.StatusCode != 201 {
			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				return err
			}
			return fmt.Errorf(`Failed to write data to elasticsearch.
                Got status code %d. Response body: %s`, resp.StatusCode, body)
		}
	}
	return nil
}
