package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/paulmach/orb/encoding/mvt"
	"github.com/paulmach/orb/geojson"
	"github.com/paulmach/orb/maptile"
	"log"
	"net/http"
	"os"
	"strconv"
)

var data []byte

func main() {
	var err error
	data, err = os.ReadFile("data.json")
	if err != nil {
		log.Fatal(err)
	}

	r := mux.NewRouter()
	r.HandleFunc("/tiles/{z}/{x}/{y}", tileHandler)
	http.Handle("/", r)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func tileHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	z, _ := strconv.ParseUint(vars["z"], 10, 32)
	x, _ := strconv.ParseUint(vars["x"], 10, 32)
	y, _ := strconv.ParseUint(vars["y"], 10, 32)

	fc, err := geojson.UnmarshalFeatureCollection(data)
	if err != nil {
		log.Fatal(err)
	}

	collections := map[string]*geojson.FeatureCollection{
		"features": fc,
	}

	layers := mvt.NewLayers(collections)

	layers.ProjectToTile(maptile.New(uint32(x), uint32(y), maptile.Zoom(z)))
	layers.Clip(mvt.MapboxGLDefaultExtentBound)
	data, err := mvt.MarshalGzipped(layers)
	if err != nil {
		fmt.Println(err)
	}

	w.Header().Set("Content-Type", "application/vnd.mapbox-vector-tile")
	w.Header().Set("Content-Encoding", "gzip")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(data)
}
