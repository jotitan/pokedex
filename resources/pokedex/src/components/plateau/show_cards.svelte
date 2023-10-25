<script>
    import {getUrl} from "../../routes/store.js";
    import ShowCard from "./show_card.svelte";
    export let size = 5;
    export let cards = [];

</script>
<div style="position:relative" class="space">
    {#each cards as card}
        <ShowCard card={card}/>
    {/each}

</div>
<style>
    img {
        height: 200px;
        border: solid 1px black;
    }

    :global(.space > div) {
        display:inline-block;
        margin-left:10px;
    }

    .space {
        border: solid 1px black;
        padding: 10px;
        border-radius: 10px;
        background-color: salmon;
    }

</style>