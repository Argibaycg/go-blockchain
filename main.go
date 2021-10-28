package main

import (
	"Argibaycg/go-blockchain/blockchain"
	"Argibaycg/go-blockchain/server"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/davecgh/go-spew/spew"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}

	go func() {
		t := time.Now()
		genesisBlock := blockchain.Block{Index: 0, Timestamp: t.String(), BPM: 0, Hash: "", PrevHash: ""}
		spew.Dump(genesisBlock)
		blockchain.Blockchain = append(blockchain.Blockchain, genesisBlock)
	}()
	log.Fatal(run())

}

//server

func run() error {
	mux := server.MakeMuxRouter()
	httpAddr := os.Getenv("PORT")
	log.Println("Listening on", httpAddr)
	s := &http.Server{
		Addr:           ":" + httpAddr,
		Handler:        mux,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	if err := s.ListenAndServe(); err != nil {
		return err
	}

	return nil
}
