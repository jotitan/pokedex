<script>
    import {getBaseUrl} from "../../routes/store.js";

    export let leftCards = 0;
    export let singleCards = true;
</script>
<div style="position:relative">
    {#if singleCards === true}

        <img src={`${getBaseUrl()}/card/image?img=cards/pokemon_card_240.png`}/>
        <div class="counter">{leftCards}</div>

    {:else}
        {#if leftCards > 0}
            {#each [...Array(leftCards).keys()] as i}
                <img src={`${getBaseUrl()}/card/image?img=cards/pokemon_card_240.png`}
                     style={`position:absolute;top:${i*30}px;left:${i*5}px`}/>
            {/each}
        {/if}
        <div></div>
    {/if}
</div>
<style>
    img {
        height:200px;
    }
    .counter {
        position: absolute;
        top: 125px;
        font-size: 18px;
        left: 80px;
        font-weight: bold;
        background-color: black;
        color: white;
        padding: 0px 10px;
    }
</style>