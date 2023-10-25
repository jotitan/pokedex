package server

import (
	"encoding/json"
	"fmt"
	"github.com/jotitan/pokedex/model"
	"github.com/jotitan/pokedex/pokedex"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

type PokemonServer struct {
	searcher     *pokedex.SearcherManager
	cardsOwned   *pokedex.CardsOwned
	folderImages string
}

func (ps PokemonServer) forceSave(w http.ResponseWriter, r *http.Request) {
	ps.cardsOwned.DoSave()
}

func (ps PokemonServer) manageDeck(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "DELETE,POST,GET")
	switch r.Method {
	case http.MethodPost:
		originalKind := strings.EqualFold("original", r.FormValue("kind"))
		ps.cardsOwned.AddDeck(r.FormValue("name"), originalKind)
	case http.MethodGet:
		decks := ps.cardsOwned.GetDecks()
		data, _ := json.Marshal(decks)
		w.Write(data)
	}

}

func (ps PokemonServer) getImageOfCard(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	if f, err := os.Open(filepath.Join(ps.folderImages, r.FormValue("img"))); err == nil {
		defer f.Close()
		io.Copy(w, f)
	} else {
		http.Error(w, "image not found", 404)
	}

}

func (ps PokemonServer) manageCard(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "DELETE,POST,GET")

	switch r.Method {
	case http.MethodPost:
		card, error := ps.searcher.FindByLink(r.FormValue("link"))
		if error != nil {
			http.Error(w, error.Error(), http.StatusNotFound)
		}
		nb := ps.cardsOwned.AddCard(card)
		w.Write([]byte(fmt.Sprintf("{\"nb\":%d}", nb)))
	case http.MethodDelete:
		card, error := ps.searcher.FindByLink(r.FormValue("link"))
		if error != nil {
			http.Error(w, error.Error(), http.StatusNotFound)
		}
		nb := ps.cardsOwned.RemoveCard(card, strings.EqualFold("true", r.FormValue("all")))
		w.Write([]byte(fmt.Sprintf("{\"nb\":%d}", nb)))
	case http.MethodGet:
		tempCards := ps.cardsOwned.GetAll()
		cards := make([]model.CardWithNb, len(tempCards))
		for i, c := range tempCards {
			if c.Card != nil {
				cards[i] = c.Card.WithNb(c.Copy)
			}
		}
		data, _ := json.Marshal(cards)
		w.Write(data)
	}
}

func (ps PokemonServer) manageCardOfDeck(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "DELETE,POST,GET")

	originalKind := strings.EqualFold("original", r.FormValue("kind"))
	nameDeck := r.FormValue("deck")
	switch r.Method {
	case http.MethodPost:
		card, error := ps.searcher.FindByLink(r.FormValue("link"))
		if error != nil {
			http.Error(w, error.Error(), http.StatusNotFound)
		}
		nb := ps.cardsOwned.AddCardToDeck(nameDeck, originalKind, card)
		w.Write([]byte(fmt.Sprintf("{\"nb\":%d}", nb)))
	case http.MethodDelete:
		card, error := ps.searcher.FindByLink(r.FormValue("link"))
		if error != nil {
			http.Error(w, error.Error(), http.StatusNotFound)
		}
		nb := ps.cardsOwned.RemoveCardFromDeck(nameDeck, originalKind, card, strings.EqualFold("true", r.FormValue("all")))
		w.Write([]byte(fmt.Sprintf("{\"nb\":%d}", nb)))
	case http.MethodGet:
		cards := ps.cardsOwned.GetCardsOfDeck(nameDeck, originalKind)
		data, _ := json.Marshal(cards)
		w.Write(data)
	}
}

func (ps PokemonServer) search(w http.ResponseWriter, r *http.Request) {
	var cards []*model.Card
	switch r.FormValue("kind") {
	case "name":
		cards = ps.searcher.FindByName(r.FormValue("name"))
	case "extension":
		cards = ps.searcher.FindByExtension(r.FormValue("name"))
	}
	// Improve with owned cards
	cardsWithNb := make([]model.CardWithNb, len(cards))
	for i, card := range cards {
		cardsWithNb[i] = card.WithNb(ps.cardsOwned.GetNb(card.Link))
	}
	data, _ := json.Marshal(cardsWithNb)
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Write(data)
}

func RunServer(dataFolder string) {
	searcher := pokedex.NewSearcher(dataFolder)
	ps := PokemonServer{searcher, pokedex.NewCardsOwned(searcher), filepath.Join(dataFolder, "images")}
	server := http.NewServeMux()
	server.HandleFunc("/search", ps.search)
	server.HandleFunc("/card", ps.manageCard)
	server.HandleFunc("/deck", ps.manageDeck)
	server.HandleFunc("/deck/save", ps.forceSave)
	server.HandleFunc("/deck/card", ps.manageCardOfDeck)
	server.HandleFunc("/card/image", ps.getImageOfCard)

	log.Println("Run app on 9002")
	http.ListenAndServe(":9002", server)

}
