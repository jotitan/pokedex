<script>
    import {page} from '$app/stores';
    import {goto} from '$app/navigation';

    import {cards, updateCards, url,forceSearch,showCardWithNb} from '../routes/store.js';
    showCardWithNb.set(false)
    let len = 0
    let nbCards = 0
    cards.subscribe(c => {
        len = c.length
        nbCards = c.reduce((total,c)=>total+c.Nb,0)
    });
    url.update(()=>$page.url)
    updateCards()


    let searchType = "name";
    const launchSearch = value => {
        if(value !== ""){
            goto(`/?name=${value}&kind=${searchType}`)
            forceSearch.update(()=>{return {kind:searchType,name:value}})
        }
    }
    const checkSearch = e => e.key === 'Enter'?launchSearch(e.target.value):{};

</script>


<div class="header">
    <span class="title">Pokemon</span>
    <input on:keyup={checkSearch} placeholder="Pokemon ou extension"/>
    <input type="radio" id="name" value="name" name="type" checked on:click={e=>searchType = e.target.value}/>
    <label for="name">Par nom</label>
    <input type="radio" id="extension" value="extension" name="type" on:click={e=>searchType = e.target.value}/>
    <label for="extension">Par extension</label>

    <a href="user_cards">
    <span class="details" title="Cartes pokemon">
        {len} - {nbCards}
    </span>
    </a>
    <a href="show_deck">
    <span class="decks">
        <img src="/deck.png" style="width:26px;height:26px;"/>
    </span>
    </a>
</div>

<style>
    .decks {
        position:absolute;
        right:120px;
        line-height:1.5;
        height:25px;
    }

    .details {
        position:absolute;
        right:30px;
        padding-left:25px;
        background-image: url("/pokeball.svg");
        background-repeat: no-repeat;
        background-size:22px;
        line-height:1.5;
        height:25px;
        font-weight:bold;
        cursor:pointer;
    }
    .details > a {
        text-decoration: none;
        color:black
    }

    .title {
        font-weight:bold;
        font-variant: small-caps;
        font-size:20px;
        margin-left:20px;
    }

    .header {
        height:30px;
        border-bottom:solid 1px gray;
        margin-bottom:5px;
    }
</style>