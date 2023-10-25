package model

import (
	"fmt"
	"strings"
)

type Card struct {
	CommonCard
	Pokemon  *CardPokemon
	Dresseur *CardDresseur `json:"Dresseur,omitempty"`
	Energy   *CardEnergy   `json:"Energy,omitempty"`
}

type CommonCard struct {
	Img       string
	Link      string
	Extension string
	Number    int
	Special   string
	Name      string
	SubName   string
}

type CardPokemon struct {
	Details    PokemonCardDetail
	Attacks    []Attack `json:"Attacks,omitempty"`
	Resistance Resistance
	Retirement Retirement
	Weekness   Weekness
	//Pokemon    *Pokemon
}

func NewCardPokemon(c CommonCard) Card {
	return Card{CommonCard: c}
}

type Resistance struct {
	Energy Energy
	Value  int
}

type Retirement struct {
	Nb int
}

type Weekness struct {
	Energy Energy
	Factor int
}

type CardDresseur struct {
}
type CardEnergy struct {
}

type PokemonCardDetail struct {
	PV          int
	Level       int
	EvolutionOf string
	TypePokemon Energy
}

func (pcd PokemonCardDetail) String() string {
	return fmt.Sprintf("Detail : Type %s, %d PV, Level %d, evolution of %s", pcd.TypePokemon, pcd.PV, pcd.Level, pcd.EvolutionOf)
}

type Attack struct {
	Name        string
	Description string
	Cost        string
	Energies    []Energy
}

func (a Attack) String() string {
	return fmt.Sprintf("Attack %s : %s (%s). %d", a.Name, a.Description, a.Cost, len(a.Energies))
}

type Energy string

const (
	grassEnergy    = Energy("plante")
	fireEnergy     = Energy("feu")
	waterEnergy    = Energy("eau")
	electricEnergy = Energy("électrique")
	psyEnergy      = Energy("psy")
	fightingEnergy = Energy("combat")
	darkEnergy     = Energy("obscurité")
	metalEnergy    = Energy("métal")
	normalEnergy   = Energy("normal")
	fairyEnergy    = Energy("fée")
	dragonEnergy   = Energy("dragon")
)

var mapEnergies = map[string]Energy{
	"fire":      fireEnergy,
	"water":     waterEnergy,
	"fighting":  fightingEnergy,
	"lightning": electricEnergy,
	"psychic":   psyEnergy,
	"grass":     grassEnergy,
	"darkness":  darkEnergy,
	"metal":     metalEnergy,
	"dragon":    dragonEnergy,
	"dairy":     fairyEnergy,
}

func GetEnergyFromName(name string) Energy {
	if energy, exist := mapEnergies[strings.ToLower(name)]; exist {
		return energy
	}
	return normalEnergy
}

func GetLevel(name string) int {
	switch name {
	case "Niveau 2":
		return 2
	case "Niveau 1":
		return 1
	default:
		return 0
	}

}

func (c Card) WithNb(nb int) CardWithNb {
	return CardWithNb{c, nb}
}

type CardWithNb struct {
	Card
	Nb int
}

func (c Card) String() string {
	return fmt.Sprintf("%s (%s) : %d (%s)", c.Name, c.Extension, c.Number, c.Special)
}

type Pokemon struct {
	Name   string `json:"name"`
	Short  string `json:"slug"`
	Image  string `json:"ThumbnailImage"`
	Number string `json:"number"`
}

type Cards struct {
	Dresseurs []*Card
	Energies  []*Card
	Pokemons  []*Card
}

func (cards Cards) GetImagesUrl() []string {
	images := make([]string, 0, len(cards.Dresseurs)+len(cards.Pokemons)+len(cards.Energies))
	for _, c := range cards.Pokemons {
		images = append(images, c.Img)
	}
	for _, c := range cards.Dresseurs {
		images = append(images, c.Img)
	}
	for _, c := range cards.Energies {
		images = append(images, c.Img)
	}
	return images
}
