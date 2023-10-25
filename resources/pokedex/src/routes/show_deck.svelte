<script>
    import Header from "../components/header.svelte";
    import {getBaseUrl, url, cards,showCardWithNb} from './store.js';
    import Modal, {bind} from 'svelte-simple-modal';
    import {writable} from 'svelte/store'
    import Popup from '../components/edit_deck_name.svelte';
    import ShowCards from '../components/show_cards.svelte';
    import {page} from "$app/stores";
    const show = writable(null);
    let decks = [];
    let promiseDeck = Promise.resolve([]);
    let cardsOfPlayer = []
    cards.subscribe(c=>cardsOfPlayer = c)
    url.update(()=>$page.url)
    let currentDeck = {};
    let totalOfDeck = 0;

    const getDecks = () => fetch(`${getBaseUrl()}/deck`).then(d=>d.json()).then(d=>decks = d);

    const computeSizeDeck = ()=>totalOfDeck = cardsOfPlayer.reduce((total,c)=>total+c.Nb,0);

    const addDeckAndHide = (deck,kind) =>
        fetch(`${getBaseUrl()}/deck?kind=${kind}&name=${deck}`,{method:'POST'}).then(d=> {
                decks = [...decks, {Name:deck,Kind:kind}]
                show.set(null)
            }
        );

    const getCards = deck=>{
        if(deck == null){
            promiseDeck = Promise.resolve([]);
            return;
        }
        currentDeck = deck;
        promiseDeck =
            fetch(`${getBaseUrl()}/deck/card?deck=${deck.Name}&kind=${deck.Kind}`)
                .then(d=>d.json())
                .then(foundCards=>{
                    cardsOfPlayer.forEach(c=>c.Nb = 0)
                    foundCards.forEach(c=>cardsOfPlayer.find(cp=>cp.Link === c.Card.Link).Nb = c.Copy)
                    computeSizeDeck()
                    return cardsOfPlayer
                })
    }

    const toggleAllCard = e => showCardWithNb.update(()=>e.target.checked);

    const selectDeck = e => getCards(decks[e.target.value]);

    const showAddDeck = ()=> show.set(bind(Popup,{addDeck:addDeckAndHide}))

    const addCardToDeck = card =>
        fetch(`${getBaseUrl()}/deck/card?kind=${currentDeck.Kind}&deck=${currentDeck.Name}&link=${card.Link}`,{method:'POST'})
            .then(d=>d.json())
            .then(d=>card.Nb = d.nb)
            .then(computeSizeDeck);

    const removeCardToDeck = (card,all) =>
        fetch(`${getBaseUrl()}/deck/card?kind=${currentDeck.Kind}&deck=${currentDeck.Name}&all=${all}&link=${card.Link}`,{method:'DELETE'})
            .then(d=>d.json())
            .then(d=> card.Nb = d.nb)
            .then(computeSizeDeck);
</script>

<Header/>
<h3>Liste des decks</h3>

{#await getDecks()}
    loading...
{:then []}
    <select on:change={selectDeck}>
        <option>Sélectionnez un deck</option>
            {#each decks as d,i}
                <option value={i}>{d.Name} ({d.Kind})</option>
            {/each}
    </select>
{/await}

<Modal show={$show}><button on:click={showAddDeck}>Ajouter deck</button></Modal>
<span style="font-weight:bold">{totalOfDeck} / 60</span>
<input type="checkbox" on:click={toggleAllCard}/> Cacher hors deck
<div>
    {#await promiseDeck}
        loading...
    {:then cards}
        {#if cards.length === 0}
        {:else}
            <ShowCards cards={cards} overrideAdd={addCardToDeck} overrideRemove={removeCardToDeck} isSearchByName={false}/>
        {/if}

    {/await}
</div>


