package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/go-redis/redis/v8"
)

func main() {
	fmt.Println("Starting the server...")

	api := NewAPI()

	http.HandleFunc("/cache", api.Handler)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", os.Getenv("PORT")), nil))
}

func (a *API) Handler(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query().Get("q")
	data, cached, err := a.getData(r.Context(), q)
	if err != nil {
		fmt.Printf("error calling data source: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	fmt.Printf("send a request for %s:", q)

	resp := APIResponse{
		Cache: cached,
		Data:  data,
	}

	err = json.NewEncoder(w).Encode(resp)
	if err != nil {
		fmt.Printf("error encoding response: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (a *API) getData(ctx context.Context, q string) ([]NominatimResponse, bool, error) {
	// is query cached?
	val, err := a.cache.Get(ctx, q).Result()
	if err == redis.Nil {
		escapedQ := url.PathEscape(q)
		address := fmt.Sprintf("https://nominatim.openstreetmap.org/search?q=%s&format=json", escapedQ)

		resp, err := http.Get(address)
		if err != nil {
			return nil, false, err
		}

		data := make([]NominatimResponse, 0)

		err = json.NewDecoder(resp.Body).Decode(&data)
		if err != nil {
			return nil, false, err
		}

		b, err := json.Marshal(data)
		if err != nil {
			return nil, false, err
		}

		err = a.cache.Set(ctx, q, bytes.NewBuffer(b).Bytes(), time.Second*15).Err()
		if err != nil {
			return nil, false, err
		}

		return data, false, nil

	} else if err != nil {
		fmt.Printf("error on redis: %v\n", err)
		return nil, false, err
	} else {
		data := make([]NominatimResponse, 0)

		err := json.Unmarshal(bytes.NewBufferString(val).Bytes(), &data)
		if err != nil {
			return nil, false, err
		}

		return data, true, nil
	}
}

type API struct {
	cache *redis.Client
}

func NewAPI() *API {
	redisAddress := fmt.Sprintf("%s:6379", os.Getenv("REDIS_URL"))

	rdb := redis.NewClient(&redis.Options{
		Addr:     redisAddress,
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	return &API{
		cache: rdb,
	}
}
