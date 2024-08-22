package main

import (
	"log"
	"net/http"
	"os"
	"path"
	"sync"

	"github.com/rajatsharma/raftkv/config"
	"github.com/rajatsharma/raftkv/fsm"
	"github.com/rajatsharma/raftkv/httpserver"
	"github.com/rajatsharma/raftkv/node"
)

func main() {
	cfg := config.GetConfig()

	db := &sync.Map{}
	kf := &fsm.Fsm{Db: db}

	dataDir := "data"
	err := os.MkdirAll(dataDir, os.ModePerm)
	if err != nil {
		log.Fatalf("Could not create data directory: %s", err)
	}

	r, err := node.SetupRaft(path.Join(dataDir, "node"+cfg.Id), cfg.Id, "localhost:"+cfg.RaftPort, kf)
	if err != nil {
		log.Fatal(err)
	}

	hs := httpserver.HttpServer{R: r, Db: db}

	http.HandleFunc("/set", hs.SetHandler)
	http.HandleFunc("/get", hs.GetHandler)
	http.HandleFunc("/join", hs.JoinHandler)

	if err := http.ListenAndServe(":"+cfg.HttpPort, nil); err != nil {
		panic(err)
	}
}
