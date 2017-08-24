package gmaps

import (
	"context"
	"encoding/json"
	"os"
	"time"

	"github.com/goware/lg"
	"github.com/pkg/errors"

	"googlemaps.github.io/maps"
)

type Client struct {
	*maps.Client
	cache map[string]*Distance
}

type Distance struct {
	Duration time.Duration `json:"duration"`
	Readable string        `json:"readable"`
	Distance maps.Distance `json:"distance"`
}

func New() (*Client, error) {
	// cache file
	cacheFile, err := os.OpenFile("./tmp/maps_cache.json", os.O_CREATE|os.O_RDWR, os.ModeAppend)
	if err != nil {
		return nil, errors.Wrap(err, "file open")
	}
	defer cacheFile.Close()

	client := &Client{
		cache: make(map[string]*Distance),
	}

	// load the cache file into lookup table
	if err := json.NewDecoder(cacheFile).Decode(&client.cache); err != nil {
		return nil, errors.Wrap(err, "json")
	}

	client.Client, err = maps.NewClient(maps.WithAPIKey(os.Getenv("MAPS_KEY")))
	if err != nil {
		return nil, err
	}
	return client, nil
}

func (c *Client) GetDistance(origin, destination string) (*Distance, error) {
	// look up in cache
	dist, found := c.cache[destination]
	if found {
		return dist, nil
	}

	rq := &maps.DistanceMatrixRequest{
		Origins:      []string{origin},
		Destinations: []string{destination},
		Mode:         maps.TravelModeWalking,
	}
	res, err := c.DistanceMatrix(context.Background(), rq)
	if err != nil {
		return nil, err
	}

	lg.Debugf("distane to %s", destination)
	var d *Distance
	for _, r := range res.Rows {
		for _, e := range r.Elements {
			d = &Distance{
				Duration: e.Duration,
				Distance: e.Distance,
				Readable: e.Duration.String(),
			}
			break
		}
		break
	}

	// cache the result
	c.cache[destination] = d
	// write to the cache file
	cacheFile, err := os.OpenFile("./tmp/maps_cache.json", os.O_RDWR, os.ModeAppend)
	if err != nil {
		return nil, errors.Wrap(err, "file open")
	}
	if err := json.NewEncoder(cacheFile).Encode(c.cache); err != nil {
		lg.Warn(errors.Wrap(err, "cache file writing"))
	}

	return d, nil
}
