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

type ElasticsearchWriter struct {
	Host string
	URL  string
}

func NewElasticsearchWriter(host string) *ElasticsearchWriter {
	return &ElasticsearchWriter{
		Host: host,
		URL:  host + "/" + index + "/" + documentType,
	}
}

func (ew *ElasticsearchWriter) Write(playlists []Playlist) error {
	for _, p := range playlists {
		pID := p.SpotifyData.ID

		data, err := json.Marshal(p)
		if err != nil {
			return err
		}
		req, err := http.NewRequest("PUT", ew.URL+"/"+pID, bytes.NewBuffer(data))
		if err != nil {
			return err
		}
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			return err
		}
		defer resp.Body.Close()

		if resp.StatusCode != 201 && resp.StatusCode != 200 {
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
