<script>
    let name = "";
    let kind = "original"
    export let addDeck;
    const add = () =>addDeck(name,kind);
    const setKindDeck = e => kind = e.target.value;
</script>

<h2>Ajouter un deck</h2>
<input type="text" placeholder="nom" on:change={e=>name = e.target.value}/>
<input type="radio" id="originalDeck" value="original" name="deck_kind" checked on:click={setKindDeck}/>
<label for="originalDeck">Original deck</label>
<input type="radio" id="customDeck" value="custom" name="deck_kind" on:click={setKindDeck}/>
<label for="customDeck">Custom deck</label>
<button on:click={add}>Ajouter</button>