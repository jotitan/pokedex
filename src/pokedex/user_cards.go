package pokedex

import (
	"encoding/json"
	"errors"
	"github.com/jotitan/pokedex/model"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"sync"
)

const (
	threshold = 20
)

type CopyCard struct {
	Card *model.Card
	Copy int
}

func NewCardsOwned(search *SearcherManager) *CardsOwned {
	return (&CardsOwned{
		search.dataFolder,
		make(map[string]*CopyCard),
		make(map[string]map[string]*CopyCard),
		make(map[string]map[string]*CopyCard),
		0,
		&sync.Mutex{},
	}).load(search)
}

type CardsOwned struct {
	dataFolder    string
	cards         map[string]*CopyCard
	originalDecks map[string]map[string]*CopyCard
	customDecks   map[string]map[string]*CopyCard
	counter       int
	locker        *sync.Mutex
}

func (co *CardsOwned) load(search *SearcherManager) *CardsOwned {
	data, err := ioutil.ReadFile(filepath.Join(co.dataFolder, "database", "player_cards.json"))
	if err == nil {
		dto := dtoToSave{}
		json.Unmarshal(data, &dto)
		for link, nb := range dto.Cards {
			foundCard := search.getCardByLink(link)
			if foundCard == nil {
				log.Println("Impossible to find card : ", link)
			} else {
				co.cards[link] = &CopyCard{foundCard, nb}
			}
		}
		transformDeck(dto.OriginalDecks, search, co.originalDecks)
		transformDeck(dto.CustomDecks, search, co.customDecks)
	}
	return co
}

func transformDeck(m map[string]map[string]int, search *SearcherManager, decks map[string]map[string]*CopyCard) {
	for deckName, cards := range m {
		list := make(map[string]*CopyCard, len(cards))
		for link, nb := range cards {
			list[link] = &CopyCard{search.getCardByLink(link), nb}
		}
		decks[deckName] = list
	}
}

func (co CardsOwned) GetNb(link string) int {
	if card, exist := co.cards[link]; exist {
		return card.Copy
	}
	return 0
}

func (co *CardsOwned) AddCard(cardToAdd *model.Card) int {
	defer co.saveIfNecessary()
	if card, exist := co.cards[cardToAdd.Link]; exist {
		card.Copy++
	} else {
		co.cards[cardToAdd.Link] = &CopyCard{Card: cardToAdd, Copy: 1}
	}
	return co.cards[cardToAdd.Link].Copy
}

func (co *CardsOwned) addCardToMap(m map[string]*CopyCard, cardToAdd *model.Card) int {
	defer co.saveIfNecessary()
	if card, exist := m[cardToAdd.Link]; exist {
		card.Copy++
	} else {
		m[cardToAdd.Link] = &CopyCard{Card: cardToAdd, Copy: 1}
	}
	return m[cardToAdd.Link].Copy
}

func (co *CardsOwned) saveIfNecessary() {
	co.locker.Lock()
	co.counter++
	// Add lock
	if co.counter > threshold {
		co.DoSave()
	}
	co.locker.Unlock()
}

func (co *CardsOwned) DoSave() {
	co.Save()
	co.counter = 0
}

func (co *CardsOwned) RemoveCard(cardToRemove *model.Card, removeAll bool) int {
	defer co.saveIfNecessary()
	if card, exist := co.cards[cardToRemove.Link]; exist {
		if removeAll || card.Copy == 1 {
			card.Copy = 0
			delete(co.cards, cardToRemove.Link)
			return 0
		} else {
			card.Copy--
			return card.Copy
		}
	}
	return 0
}

func (co *CardsOwned) AddDeck(name string, isOriginal bool) {
	if isOriginal {
		co.originalDecks[name] = make(map[string]*CopyCard, 0)
	} else {
		co.customDecks[name] = make(map[string]*CopyCard, 0)
	}
	co.saveIfNecessary()
}

func (co *CardsOwned) RemoveCardFromDeck(deckName string, isOriginal bool, cardToRemove *model.Card, removeAll bool) int {
	deck, err := co.getDeck(deckName, isOriginal)
	if err != nil {
		return 0
	}
	defer co.saveIfNecessary()
	if card, exist := deck[cardToRemove.Link]; exist {
		if removeAll || card.Copy == 1 {
			card.Copy = 0
			delete(deck, cardToRemove.Link)
			return 0
		} else {
			card.Copy--
			return card.Copy
		}
	}
	return 0
}

func (co *CardsOwned) getDeck(deckName string, isOriginal bool) (map[string]*CopyCard, error) {
	var deck map[string]*CopyCard
	exist := false
	if isOriginal {
		deck, exist = co.originalDecks[deckName]
	} else {
		deck, exist = co.customDecks[deckName]
	}
	if !exist {
		return nil, errors.New("no deck")
	}
	return deck, nil
}

func (co *CardsOwned) AddCardToDeck(deckName string, isOriginal bool, cardToAdd *model.Card) int {
	deck, err := co.getDeck(deckName, isOriginal)
	if err != nil {
		return 0
	}
	defer co.saveIfNecessary()
	if card, exist := deck[cardToAdd.Link]; exist {
		card.Copy++
	} else {
		deck[cardToAdd.Link] = &CopyCard{Card: cardToAdd, Copy: 1}
	}
	return deck[cardToAdd.Link].Copy
}

type dtoToSave struct {
	Cards         map[string]int
	OriginalDecks map[string]map[string]int
	CustomDecks   map[string]map[string]int
}

type Deck struct {
	Name string
	Kind string
}

func (co *CardsOwned) Save() error {
	owned := dtoToSave{cardsToList(co.cards), mapCardToMapLink(co.originalDecks), mapCardToMapLink(co.customDecks)}
	data, _ := json.Marshal(owned)
	co.counter = 0
	return ioutil.WriteFile(filepath.Join(co.dataFolder, "database", "player_cards.json"), data, os.ModePerm)
}

func (co *CardsOwned) GetAll() []*CopyCard {
	return mapToList(co.cards)
}

func (co *CardsOwned) GetDecks() []Deck {
	decks := make([]Deck, 0, len(co.originalDecks)+len(co.customDecks))
	for d := range co.originalDecks {
		decks = append(decks, Deck{d, "original"})
	}
	for d := range co.customDecks {
		decks = append(decks, Deck{d, "custom"})
	}
	return decks
}

func (co *CardsOwned) GetCardsOfDeck(deck string, isOriginal bool) []*CopyCard {
	if isOriginal {
		return getCardsOfDesk(deck, co.originalDecks)
	}
	return getCardsOfDesk(deck, co.customDecks)
}

func getCardsOfDesk(deck string, m map[string]map[string]*CopyCard) []*CopyCard {
	mcards, exist := m[deck]
	if exist {
		return mapToList(mcards)
	}
	return []*CopyCard{}
}

func mapToList[V interface{}](m map[string]V) []V {
	list := make([]V, 0, len(m))
	for _, value := range m {
		list = append(list, value)
	}
	return list
}

func cardsToList(m map[string]*CopyCard) map[string]int {
	cards := make(map[string]int, len(m))
	for _, card := range m {
		cards[card.Card.Link] = card.Copy
	}
	return cards
}

func mapCardToMapLink(m map[string]map[string]*CopyCard) map[string]map[string]int {
	m2 := make(map[string]map[string]int, len(m))
	for deck, cards := range m {
		m2[deck] = cardsToList(cards)
	}
	return m2
}
