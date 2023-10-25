package pokedex

import (
	"fmt"
	"sort"
	"strings"
)

func Analyze(dataFolder string) {
	s := NewSearcher(dataFolder)
	AnalyzeAttacks(s)
	AnalyzeDeepAttacks(s)
}

func AnalyzeAttacks(s *SearcherManager) {
	max := 0
	nbAttacks := 0
	attacksRepartitions := make([]int, 5)
	for _, p := range s.cards.Pokemons {
		if p.Pokemon != nil {
			att := p.Pokemon.Attacks
			if len(att) > max {
				max = len(att)
			}
			nbAttacks += len(att)
			attacksRepartitions[len(att)]++
		}
	}
	fmt.Println("--------- ANALYSE ATTACKS ---------")
	fmt.Println(fmt.Sprintf("POKEMON NUMBER : %d", len(s.cards.Pokemons)))
	fmt.Println(fmt.Sprintf("ATTACK MAX %d", max))
	fmt.Println(fmt.Sprintf("AVERAGE ATTACH %.2f", float64(nbAttacks)/float64(len(s.cards.Pokemons))))
	fmt.Println("REPARTITION", attacksRepartitions)
}

// Cartographie les mots
func AnalyzeDeepAttacks(s *SearcherManager) {
	repartitions := make(map[string]int)
	equality := make(map[string]int)
	names := make(map[string]int)
	nbAttacks := 0
	for _, p := range s.cards.Pokemons {
		if p.Pokemon != nil {
			att := p.Pokemon.Attacks
			nbAttacks += len(att)
			for _, a := range att {
				nbLine := equality[a.Description]
				equality[a.Description] = nbLine + 1
				nbName := names[a.Name]
				names[a.Name] = nbName + 1
				for _, word := range strings.Split(strings.ToLower(a.Description), " ") {
					nb := repartitions[word]
					repartitions[word] = nb + 1
				}
			}
		}
	}
	instances := transform(repartitions)
	instancesNames := transform(names)
	instancesDescriptions := transform(equality)
	sort.SliceStable(instances, func(i, j int) bool { return instances[i].nb < instances[j].nb })
	fmt.Println("--------- ANALYSE DEEP ATTACKS ---------")
	fmt.Println(instances)
	fmt.Println("Nb attacks", nbAttacks)
	fmt.Println("Different names", len(instancesNames))
	fmt.Println("Different descriptions", len(instancesDescriptions))
	fmt.Println("Empty descriptions", equality[""])
}

type sub struct {
	value string
	nb    int
}

func transform(m map[string]int) []sub {
	instances := make([]sub, 0)
	for word, nb := range m {
		instances = append(instances, sub{word, nb})
	}
	sort.SliceStable(instances, func(i, j int) bool { return instances[i].nb < instances[j].nb })
	return instances
}
