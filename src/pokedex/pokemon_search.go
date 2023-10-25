package pokedex

import (
	"encoding/json"
	"errors"
	"github.com/jotitan/pokedex/model"
	"io/ioutil"
	"path/filepath"
	"strings"
)

type SearcherManager struct {
	dataFolder string
	cards      model.Cards
	// Used to search a specific card, by link card (unique)
	cardsByLink map[string]*model.Card
	byName      map[string][]*model.Card
	byExtension map[string][]*model.Card
}

func NewSearcher(dataFolder string) *SearcherManager {
	s := &SearcherManager{dataFolder: dataFolder}
	s.load()
	return s
}

func (s *SearcherManager) group() {
	s.byName = make(map[string][]*model.Card)
	s.byExtension = make(map[string][]*model.Card)
	s.cardsByLink = make(map[string]*model.Card)
	s.groupByType(s.cards.Pokemons, false)
	s.groupByType(s.cards.Dresseurs, true)
	s.groupByType(s.cards.Energies, true)
}

func (s *SearcherManager) groupByType(cards []*model.Card, useSub bool) {
	for _, c := range cards {
		name := extractName(c, useSub)
		s.indexName(name, c)
		if s.byExtension[clean(c.Extension)] == nil {
			s.byExtension[clean(c.Extension)] = make([]*model.Card, 0)
		}
		s.byExtension[clean(c.Extension)] = append(s.byExtension[clean(c.Extension)], c)
		s.cardsByLink[c.Link] = c
	}
}

func extractName(c *model.Card, useSub bool) string {
	if useSub {
		return clean(c.SubName)
	}
	return clean(c.Name)
}

func (s *SearcherManager) indexName(name string, c *model.Card) {
	for _, sub := range strings.Split(name, " ") {
		if s.byName[sub] == nil {
			s.byName[sub] = make([]*model.Card, 0)
		}
		s.byName[sub] = append(s.byName[sub], c)
	}
}

func (s SearcherManager) getCardByLink(link string) *model.Card {
	return s.cardsByLink[link]
}

func clean(value string) string {
	return strings.ToLower(strings.Trim(value, " "))
}

func (s *SearcherManager) FindByLink(link string) (*model.Card, error) {
	if card, exist := s.cardsByLink[link]; exist {
		return card, nil
	}
	return nil, errors.New("card not found")
}

func (s *SearcherManager) FindByName(name string) []*model.Card {
	if cards, exist := s.byName[name]; exist {
		return cards
	}
	return make([]*model.Card, 0)
}

func (s *SearcherManager) FindByExtension(extension string) []*model.Card {
	if cards, exist := s.byExtension[extension]; exist {
		return cards
	}
	return make([]*model.Card, 0)
}

func (s *SearcherManager) load() {
	data, _ := ioutil.ReadFile(filepath.Join(s.dataFolder, "database", "cards_pokemon.json"))
	s.cards = model.Cards{}
	json.Unmarshal(data, &s.cards)
	s.group()
}

func (s SearcherManager) extractPokemon() map[string]*model.Pokemon {
	pokemons := make(map[string]*model.Pokemon)
	for _, c := range s.cards.Pokemons {
		/*p*/ _, exist := pokemons[c.Name]
		if exist {
			c.Pokemon = &model.CardPokemon{} //p
		} else {
			pokemons[c.Name] = &model.Pokemon{} //c.Pokemon
		}
	}
	return pokemons
}
