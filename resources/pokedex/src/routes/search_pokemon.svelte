<script>
    import {forceSearch, getBaseUrl, showCardWithNb} from './store.js';
    import ShowCards from '../components/show_cards.svelte';
    import ShowPokemon from '../components/show_pokemon.svelte'

    let searchType = "name";
    let isDresseur = false;
    let searchPromise = Promise.resolve([]);
    let showExtension = false;
    let saveResults = [];
    let selectedFilter = "";
    const search = criterias => {
        if (criterias.name === "") {
            return Promise.resolve([]);
        }
        showExtension = criterias.name !== "dresseur";
        searchPromise = fetch(`${getBaseUrl()}/search?kind=${criterias.kind}&name=${criterias.name}`)
            .then(d => d.json())
            .then(d => d.sort((c1, c2) => c1.Number - c2.Number))
            .then(d => {
                saveResults = d
                return saveResults
            })
    }

    forceSearch.subscribe(search)
    const toggleAllCard = e => showCardWithNb.update(() => e.target.checked);
    const filterCard = e => {
        selectedFilter = e.target.value;
        if(selectedFilter !== ""){
            if(selectedFilter === "pokemon"){
                searchPromise = Promise.resolve(saveResults.filter(c => c.Name.toLocaleLowerCase() !== "énergie" && c.Name.toLocaleLowerCase() !== "dresseur"));
            }else {
                searchPromise = Promise.resolve(saveResults.filter(c => c.Name.toLocaleLowerCase() === selectedFilter));
            }
        } else {
            searchPromise = Promise.resolve(saveResults);
        }
    }

</script>

{#await searchPromise}...searching
{:then cards}
    <span class="nb-results">{cards != null ? cards.length : 0} résultat(s)</span>
    <input type="checkbox" on:click={toggleAllCard}/> Mes cartes

    <input type="radio" value="" name="filter-type" id="id_all" on:click={filterCard} checked={selectedFilter === ""}/>
    <label for="id_all">Tout</label>
    <input type="radio" value="dresseur" name="filter-type" id="id_dresseur" on:click={filterCard} checked={selectedFilter === "dresseur"}/>
    <label for="id_dresseur">Dresseur</label>
    <input type="radio" value="énergie" name="filter-type" id="id_energie" on:click={filterCard} checked={selectedFilter === "énergie"}/>
    <label for="id_energie">Energie</label>
    <input type="radio" value="pokemon" name="filter-type" id="id_pokemon" on:click={filterCard} checked={selectedFilter === "pokemon"}/>
    <label for="id_pokemon">Pokemon</label>

    {#if cards.length > 0}

        {#if searchType === "name"}
            <ShowPokemon pokemon={cards[0].Poke}/>
        {:else}
            <br/>
        {/if}
        <div>
            <ShowCards cards={cards} isSearchByName={searchType === "name" || showExtension }/>
        </div>
    {/if}
{/await}

<style>
    .nb-results {
        font-weight: bold;
        padding-left: 20px;
    }

</style>