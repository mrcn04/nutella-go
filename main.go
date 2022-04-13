package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
)

func main() {
	fmt.Println("Starting the server...")

	http.HandleFunc("/cache", Handler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func Handler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("**** in the handler ****")

	q := r.URL.Query().Get("q")
	data, err := getData(q)
	if err != nil {
		fmt.Printf("error calling data source: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	resp := APIResponse{
		Cache: false,
		Data:  data,
	}

	err = json.NewEncoder(w).Encode(resp)
	if err != nil {
		fmt.Printf("error encoding response: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func getData(q string) ([]NominatimResponse, error) {
	escapedQ := url.PathEscape(q)
	address := fmt.Sprintf("https://nominatim.openstreetmap.org/search?q=%s&format=json", escapedQ)

	resp, err := http.Get(address)
	if err != nil {
		return nil, err
	}

	data := make([]NominatimResponse, 0)
	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		return nil, err
	}

	return data, nil
}
