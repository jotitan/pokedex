<script>
    import {getUrl} from "../../routes/store.js";

    export let card;

</script>
<div style="position:relative" class="img-place">
    {#if card != null}
        <img src={getUrl(card)}/>
    {/if}
</div>
<style>
    img {
        height: 200px;
        border:solid 1px black;
    }

    .img-place {
        border: solid 1px black;
        padding: 10px;
        border-radius: 10px;
        background-color:white;
    }

</style>