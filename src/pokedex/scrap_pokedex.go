package pokedex

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/jotitan/pokedex/model"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"
)

const BaseUrl = "https://www.pokemon.com"

func ScrapDresseur(dataFolder string) {
	begin := time.Now()
	cards := extractDresseurs(dataFolder)
	log.Println("LOAD ALL DRESSEURS", len(cards), "in", time.Now().Sub(begin))
}

func ScrapEneryies(dataFolder string) {
	begin := time.Now()
	cards := extractEnergies(dataFolder)
	log.Println("LOAD ALL ENERGIES", len(cards), "in", time.Now().Sub(begin))
}

func loadExistingImages(folder string) map[string]bool {
	m := make(map[string]bool)
	dir, err := os.Open(folder)
	if err == nil {
		names, _ := dir.Readdirnames(-1)
		for _, name := range names {
			m[name] = true
		}
	}
	return m
}

func ScrapImage(dataFolder string) {
	s := NewSearcher(dataFolder)
	limiter := make(chan struct{}, 5)
	waiter := sync.WaitGroup{}
	waiter.Add(1)
	existings := loadExistingImages(filepath.Join(dataFolder, "images"))
	for i, c := range s.cards.GetImagesUrl() {
		waiter.Add(1)
		limiter <- struct{}{}
		if i%100 == 0 {
			log.Println("Copy", i, "image")
		}
		go func(img string) {
			name := filepath.Base(img)
			if _, exists := existings[name]; !exists {
				output := filepath.Join(dataFolder, filepath.Base(name))
				// Check if already exists
				r, _ := http.Get(img)
				f, _ := os.OpenFile(output, os.O_RDWR|os.O_CREATE|os.O_TRUNC, os.ModePerm)
				io.Copy(f, r.Body)
				f.Close()
			}
			waiter.Done()
			<-limiter
		}(c)
	}
	waiter.Done()
	waiter.Wait()
	log.Println("End of scrap image")
}

func ScrapByName(name string) {
	cards := extractCards(model.Pokemon{Name: name})
	for _, t := range cards {
		improveCard(t)
	}
	fmt.Println(cards)
}

func Scrap(dataFolder string) {
	s := NewSearcher(dataFolder)
	begin := time.Now()
	cards := make([]*model.Card, 0)
	threshold := make(chan struct{}, 10)
	limit := make(chan struct{}, 1)
	waiter := sync.WaitGroup{}
	waiter.Add(1)
	for _, poke := range getPokemons() {
		waiter.Add(1)
		threshold <- struct{}{}
		go func(p model.Pokemon) {
			tempCards := extractCards(p)
			tempPokemonCards := make([]*model.Card, 0, len(tempCards))
			for _, t := range tempCards {
				if _, exists := s.cardsByLink[t.Link]; !exists {
					if p, err := improveCard(t); err == nil {
						tempPokemonCards = append(tempPokemonCards, p)
					} else {
						log.Println("error", t.Link)
					}
				}

			}
			// release limit 25
			<-threshold
			// mutex to limit multiple write
			limit <- struct{}{}
			cards = append(cards, tempPokemonCards...)
			<-limit
			waiter.Done()
		}(poke)

	}
	waiter.Done()
	waiter.Wait()
	s.cards.Pokemons = append(s.cards.Pokemons, cards...)
	log.Println("LOAD ALL", len(cards), "in", time.Now().Sub(begin))
	save(s.cards, filepath.Join(dataFolder, "database", "cards_pokemon.json"))
}

func save(cards model.Cards, path string) {
	data, _ := json.Marshal(cards)
	ioutil.WriteFile(path, data, os.ModePerm)
	log.Println("Save cards")
}

func extractDresseurs(dataFolder string) []*model.Card {
	s := NewSearcher(dataFolder)
	cards := make([]*model.Card, 0)
	dresseursUrl := fmt.Sprintf("%s/fr/jcc-pokemon/cartes-pokemon/%d?cardName=&cardText=&evolvesFrom=&simpleSubmit=&trainer=on&trainer-pokemon-tool=on&trainer-stadium=on&trainer-supporter=on&trainer-technical-machine=on&trainer-rockets-secret-machine=on&format=unlimited&hitPointsMin=0&hitPointsMax=340&retreatCostMin=0&retreatCostMax=5&totalAttackCostMin=0&totalAttackCostMax=5&particularArtist=", BaseUrl)
	p := model.Pokemon{Name: "Dresseur"}
	for i := 1; i <= 125; i++ {
		data, _ := loadDataFromUrl(fmt.Sprintf(dresseursUrl, i))
		if strings.EqualFold("", data) {
			break
		}
		tempCards := extractCard(data, p)
		// Improve cards

		cards = append(cards, improveDresseurCards(tempCards)...)
	}
	s.cards.Dresseurs = cards
	save(s.cards, filepath.Join(dataFolder, "database", "cards_pokemon.json"))
	return cards
}

func extractEnergies(dataFolder string) []*model.Card {
	s := NewSearcher(dataFolder)
	cards := make([]*model.Card, 0)
	energiesUrl := fmt.Sprintf("%s/fr/jcc-pokemon/cartes-pokemon/%d?basic-energy=on&special-energy=on&format=unlimited&hitPointsMin=0&hitPointsMax=340&retreatCostMin=0&retreatCostMax=5&totalAttackCostMin=0&totalAttackCostMax=5&particularArtist=&advancedSubmit=&sort=number&sort=number", BaseUrl)
	p := model.Pokemon{Name: "Énergie"}
	for i := 1; i <= 15; i++ {
		data, _ := loadDataFromUrl(fmt.Sprintf(energiesUrl, i))
		if strings.EqualFold("", data) {
			break
		}
		tempCards := extractCard(data, p)
		// Improve cards

		cards = append(cards, improveDresseurCards(tempCards)...)
	}
	s.cards.Energies = cards
	save(s.cards, filepath.Join(dataFolder, "database", "cards_pokemon.json"))
	return cards
}

func improveDresseurCards(commonCards []*model.CommonCard) []*model.Card {
	waiter := sync.WaitGroup{}
	waiter.Add(len(commonCards))
	limit := make(chan struct{}, 1)
	cards := make([]*model.Card, len(commonCards))
	for i, c := range commonCards {
		go func(card *model.CommonCard, pos int) {
			if c, err := improveCard(card); err == nil {
				limit <- struct{}{}
				cards[pos] = c
				<-limit
			}
			waiter.Done()
		}(c, i)
	}
	waiter.Wait()
	return cards
}

func extractCards(p model.Pokemon) []*model.CommonCard {
	begin := time.Now()
	data, nb := loadData(p.Name, 1)
	if nb == 0 {
		return []*model.CommonCard{}
	}
	cards := extractCard(data, p)
	for i := 2; i <= nb; i++ {
		data, _ = loadData(p.Name, i)
		cards = append(cards, extractCard(data, p)...)
	}
	log.Println("GOT", len(cards), p.Short, "in", time.Now().Sub(begin))
	return cards
}

func loadDataFromUrl(urlToCall string) (string, string) {
	begin := time.Now()
	resp, err := http.Get(urlToCall)
	if err != nil || resp.StatusCode != 200 {
		log.Println("Impossible to get", urlToCall, err)
		return "", ""
	}
	d, _ := ioutil.ReadAll(resp.Body)

	data := strings.ReplaceAll(string(d), "\n", "")
	c, _ := regexp.Compile("<ul class=\"cards-grid clear\" id=\"cardResults\">(.*?)<\\/ul>")
	results := c.FindAllStringSubmatch(data, -1)
	if len(results) == 0 {
		return "", ""
	}
	log.Println("Load", urlToCall, time.Now().Sub(begin))
	return results[0][1], data
}

func loadData(name string, page int) (string, int) {
	r, data := loadDataFromUrl(fmt.Sprintf("%s/fr/jcc-pokemon/cartes-pokemon/%d?cardName=%s", BaseUrl, page, url.QueryEscape(name)))
	if strings.EqualFold("", data) {
		return "", 0
	}
	nbPage := 0
	if page == 1 {
		nbPage, _ = extractNbPage(data)
	}
	return r, nbPage
}

func extractNbPage(data string) (int, error) {
	r, _ := regexp.Compile("<span >1 sur ([0-9]+)<\\/span>")
	result := r.FindAllStringSubmatch(data, 1)
	if len(result) == 1 {
		return strconv.Atoi(result[0][1])
	}
	return 0, errors.New("impossible")
}

func extractCard(text string, p model.Pokemon) []*model.CommonCard {
	cards := make([]*model.CommonCard, 0)
	if strings.EqualFold("", text) {
		return cards
	}

	hrefs := strings.Split(text, "<li>")
	regexHref, _ := regexp.Compile("href=\"(.*?)\"")
	regexSrc, _ := regexp.Compile("src=\"(.*?)\"")
	for _, a := range hrefs {
		if !strings.EqualFold("", strings.Trim(a, " ")) {
			href := regexHref.FindAllStringSubmatch(a, -1)[0][1]
			src := regexSrc.FindAllStringSubmatch(a, -1)[0][1]
			cards = append(cards, &model.CommonCard{Img: src, Link: href, Name: p.Name})
		}
	}
	return cards
}

func getPokemons() map[string]model.Pokemon {
	begin := time.Now()
	pokemons := make([]model.Pokemon, 0)
	resp, _ := http.Get(fmt.Sprintf("%s/fr/api/pokedex/kalos", BaseUrl))
	data, _ := ioutil.ReadAll(resp.Body)
	json.Unmarshal(data, &pokemons)
	pokemonsAsMap := make(map[string]model.Pokemon)
	for _, pokemon := range pokemons {
		pokemonsAsMap[pokemon.Number] = pokemon
	}
	log.Println("LOAD POKEMONS", len(pokemonsAsMap), "in", time.Now().Sub(begin))
	return pokemonsAsMap
}

func improveCard(c *model.CommonCard) (*model.Card, error) {
	// Already treat
	if !strings.EqualFold("", c.Extension) {
		return nil, nil
	}
	log.Println("Improve model.Card", c.Name, c.Link)
	resp, err := http.Get(fmt.Sprintf("%s%s", BaseUrl, c.Link))
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != 200 {
		return nil, errors.New(fmt.Sprintf("Status %d", resp.StatusCode))
	}
	d, _ := ioutil.ReadAll(resp.Body)
	data := strings.ReplaceAll(string(d), "\n", "")
	r, _ := regexp.Compile("<div class=\"stats-footer\">(.*?)<\\/div>")
	results := r.FindAllStringSubmatch(data, -1)
	r, _ = regexp.Compile("<h1>(.*?)<")
	subNames := r.FindAllStringSubmatch(data, -1)
	if len(subNames) > 0 {
		c.SubName = subNames[0][1]
	}
	r, _ = regexp.Compile("\">(.*?)</a><\\/h3>")
	if len(results) > 0 {
		extension := r.FindAllStringSubmatch(results[0][1], -1)
		r, _ = regexp.Compile("<span>([0-9]+)/(.*?)</span>")

		nbResults := r.FindAllStringSubmatch(results[0][1], -1)
		if len(nbResults) > 0 {
			number, _ := strconv.Atoi(nbResults[0][1])
			c.Number = number
		} else {
			r, _ = regexp.Compile("<span>([^<]+)</span>")
			nbResults = r.FindAllStringSubmatch(results[0][1], -1)
			if len(nbResults) > 0 {
				c.Special = nbResults[0][1]
			}
		}
		c.Extension = extension[0][1]
	} else {
		log.Println("error card", c.Link)
	}
	p := scrapPokemonCard(d, c)
	return &p, nil
}

func scrapPokemonCard(d []byte, card *model.CommonCard) model.Card {
	p := model.NewCardPokemon(*card)

	data := strings.ReplaceAll(string(d), "\n", "")
	cardDetails := extractFromRegex("<div class=\"card-basic-info\">(.+?)<div class=\"pokemon-abilities\">", data)
	if len(cardDetails) > 0 {
		detail := cardDetails[0][1]
		pokemonCardDetail := model.PokemonCardDetail{}
		level, e := extractOneFromRegex("<h2>Pokémon (.+?)</h2>", detail)
		if e == nil {
			pokemonCardDetail.Level = model.GetLevel(level)
		}
		pv, e := extractOneFromRegex("<span class=\"card-hp\"><span>PV</span>(.+?)</span>", detail)
		if e == nil {
			if val, e2 := strconv.Atoi(pv); e2 == nil {
				pokemonCardDetail.PV = val
			}
		}
		evolution, e := extractOneFromRegex("cardName=(.+?)\">", detail)
		if e == nil {
			pokemonCardDetail.EvolutionOf = evolution
		}
		typePokemon, e := extractOneFromRegex("class=\"energy icon-(.+?)\"", detail)
		if e == nil {
			pokemonCardDetail.TypePokemon = model.GetEnergyFromName(typePokemon)
		}
		p.Pokemon = &model.CardPokemon{Details: pokemonCardDetail, Attacks: []model.Attack{}}
	}

	abilities := extractFromRegex("<div class=\"ability\">(.+?)</div>", data)
	for _, ab := range abilities {
		nameAttack, e := extractOneFromRegex("<h4 class=\"left label\">(.+?)</h4>", ab[1])
		if e == nil {
			attack := model.Attack{Name: nameAttack, Energies: make([]model.Energy, 0)}

			desc, e := extractOneFromRegex("<pre>(.+?)</pre>", ab[1])
			if e == nil {
				attack.Description = desc
			}
			cost, e := extractOneFromRegex("<span class=\"right plus\">(.+?)</span>", ab[1])
			if e == nil {
				attack.Cost = cost
			}
			energies := extractFromRegex("data-energy-type=\"(.+?)\"", ab[1])
			if len(energies) > 0 {
				for _, energy := range energies {
					attack.Energies = append(attack.Energies, model.GetEnergyFromName(energy[1]))
				}
			}
			p.Pokemon.Attacks = append(p.Pokemon.Attacks, attack)
		}
	}
	stats := extractFromRegex("<div class=\"stat(.+?)</div>", data)
	for _, stat := range stats {
		if name, e := extractOneFromRegex("<h4>(.+?)</h4>", stat[1]); e == nil {
			icons := extractFromRegex("<i class=\"energy icon-(.+?)\">", stat[1])
			if len(icons) > 0 {
				energies := make([]model.Energy, len(icons))
				for i, icon := range icons {
					energies[i] = model.GetEnergyFromName(icon[1])
				}
				switch {
				case strings.Contains(strings.ToLower(name), "retraite"):
					p.Pokemon.Retirement = model.Retirement{Nb: len(energies)}
				case strings.Contains(strings.ToLower(name), "faiblesse"):
					ratio, err := extractOneFromRegex("×([0-9]+)", stat[1])
					if value, err2 := strconv.Atoi(ratio); err == nil && err2 == nil {
						p.Pokemon.Weekness = model.Weekness{Energy: energies[0], Factor: value}
					}
				case strings.Contains(strings.ToLower(name), "résistance"):
					ratio, err := extractOneFromRegex("-([0-9]+)", stat[1])
					if value, err2 := strconv.Atoi(ratio); err == nil && err2 == nil {
						p.Pokemon.Resistance = model.Resistance{Energy: energies[0], Value: value}
					}
				}
			}
		}
	}
	return p
}

func extractFromRegex(regex, data string) [][]string {
	r, _ := regexp.Compile(regex)
	return r.FindAllStringSubmatch(data, -1)
}

func extractOneFromRegex(regex, data string) (string, error) {
	r, _ := regexp.Compile(regex)
	results := r.FindAllStringSubmatch(data, 1)
	if len(results) != 1 {
		return "", errors.New("")
	}
	return results[0][1], nil
}
