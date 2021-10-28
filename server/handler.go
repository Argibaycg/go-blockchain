package server

import (
	"Argibaycg/go-blockchain/blockchain"
	"encoding/json"
	"io"
	"net/http"

	"github.com/davecgh/go-spew/spew"
	"github.com/gorilla/mux"
)

type Message struct {
	BPM int
}

func MakeMuxRouter() http.Handler {
	muxRouter := mux.NewRouter()
	muxRouter.HandleFunc("/", HandleGetBlockchain).Methods("GET")
	muxRouter.HandleFunc("/", HandleWriteBlock).Methods("POST")
	return muxRouter
}

func HandleGetBlockchain(w http.ResponseWriter, r *http.Request) {
	bytes, err := json.MarshalIndent(blockchain.Blockchain, "", "  ")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	io.WriteString(w, string(bytes))
}

func HandleWriteBlock(w http.ResponseWriter, r *http.Request) {
	var m Message

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&m); err != nil {
		RespondWithJSON(w, r, http.StatusBadRequest, r.Body)
		return
	}
	defer r.Body.Close()

	newBlock, err := blockchain.GenerateBlock(blockchain.Blockchain[len(blockchain.Blockchain)-1], m.BPM)
	if err != nil {
		RespondWithJSON(w, r, http.StatusInternalServerError, m)
		return
	}
	if blockchain.IsBlockValid(newBlock, blockchain.Blockchain[len(blockchain.Blockchain)-1]) {
		newBlockchain := append(blockchain.Blockchain, newBlock)
		blockchain.ReplaceChain(newBlockchain)
		spew.Dump(blockchain.Blockchain)
	}

	RespondWithJSON(w, r, http.StatusCreated, newBlock)

}

func RespondWithJSON(w http.ResponseWriter, r *http.Request, code int, payload interface{}) {
	response, err := json.MarshalIndent(payload, "", "  ")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("HTTP 500: Internal Server Error"))
		return
	}
	w.WriteHeader(code)
	w.Write(response)
}
