<script>
    import Header from "../components/header.svelte";
    import ShowCards from "../components/show_cards.svelte";
    import CheckboxType from "../components/checkbox_type_pokemon.svelte";
    import {cards, getImgEnergy} from "./store.js";

    let selectedFilter = "";
    let allCards = [];
    let cardsOfPlayer = [];
    let filters = [];
    cards.subscribe(c => {
        allCards = c
        cardsOfPlayer = c
    })
    const filter = (selected, type) => {
        if (selected) {
            filters.push(type);
        } else {
            filters = filters.filter(v => v !== type);
        }
        doKindFilter()
    }
    const doKindFilter = () => {
        if (filters.length > 0) {
            cardsOfPlayer = allCards.filter(c => c.Pokemon != null && filters.includes(c.Pokemon.Details.TypePokemon))
        } else {
            cardsOfPlayer = allCards
        }
        if (selectedFilter !== "") {
            if (selectedFilter === "pokemon") {
                cardsOfPlayer = cardsOfPlayer.filter(c => c.Name.toLocaleLowerCase() !== "énergie" && c.Name.toLocaleLowerCase() !== "dresseur");
            } else {
                cardsOfPlayer = cardsOfPlayer.filter(c => c.Name.toLocaleLowerCase() === selectedFilter);
            }
        } else {
            cardsOfPlayer = cardsOfPlayer
        }
    }
    const filterByKind = e => {
        selectedFilter = e.target.value;
        doKindFilter();
    }

</script>

<Header/>

<div>
    <h3 style="display:inline-block; width:200px;top:-8px;position:relative">Mes cartes {cardsOfPlayer.length}</h3>
    <span style="top:-8px;position:relative">
    <input type="radio" value="" name="ft" id="id_all" on:click={filterByKind} checked={selectedFilter === ""}/>
    <label for="id_all">Tout</label>
    <input type="radio" value="dresseur" name="ft" id="id_dresseur" on:click={filterByKind}
           checked={selectedFilter === "dresseur"}/>
    <label for="id_dresseur">Dresseur</label>
    <input type="radio" value="énergie" name="ft" id="id_energie" on:click={filterByKind}
           checked={selectedFilter === "énergie"}/>
    <label for="id_energie">Energie</label>
    <input type="radio" value="pokemon" name="ft" id="id_pokemon" on:click={filterByKind}
           checked={selectedFilter === "pokemon"}/>
    <label for="id_pokemon">Pokemon</label>
        </span>
    {#if selectedFilter === "pokemon"}
        <CheckboxType path={getImgEnergy("normal")} onclick={on=>filter(on,'normal')} title="Normal"/>
        <CheckboxType path={getImgEnergy("combat")} onclick={on=>filter(on,'combat')} title="Combat"/>
        <CheckboxType path={getImgEnergy("feu")} onclick={on=>filter(on,'feu')} title="Feu"/>
        <CheckboxType path={getImgEnergy("électrique")} onclick={on=>filter(on,'électrique')} title="Électrique"/>
        <CheckboxType path={getImgEnergy("eau")} onclick={on=>filter(on,'eau')} title="Eau"/>
        <CheckboxType path={getImgEnergy("plante")} onclick={on=>filter(on,'plante')} title="Plante"/>
        <CheckboxType path={getImgEnergy("psy")} onclick={on=>filter(on,'psy')} title="Psy"/>
        <CheckboxType path={getImgEnergy("dragon")} onclick={on=>filter(on,'dragon')} title="Dragon"/>
        <CheckboxType path={getImgEnergy("fee")} onclick={on=>filter(on,'fée')} title="Fée"/>
        <CheckboxType path={getImgEnergy("obscurité")} onclick={on=>filter(on,'obscurité')} title="Obscurité"/>
        <CheckboxType path={getImgEnergy("métal")} onclick={on=>filter(on,'métal')} title="Métal"/>
    {/if}
</div>

<div>
    <ShowCards cards={cardsOfPlayer.sort((c1,c2)=>c1.Name.localeCompare(c2.Name))} isSearchByName={false}
               showSelected={false}/>
</div>




