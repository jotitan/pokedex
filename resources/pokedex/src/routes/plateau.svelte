<script>
    import {page} from '$app/stores';
    import PlayerSpace from "../components/plateau/player_space.svelte"
    import {updateCards, cards, url} from "./store.js";

    let myCards;
    url.update(()=>$page.url)
    cards.subscribe(c => myCards=c)
    updateCards()
</script>

<div style="background-color:lightgreen;height:47vh;padding-top:20px;">
    <PlayerSpace isTop={true}/>
</div>
<div style="background-color:lightsalmon;height:47vh; position:relative;padding-top:20px;">
    <PlayerSpace/>
</div>
