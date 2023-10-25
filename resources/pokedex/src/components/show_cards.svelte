<script>
    import ShowCard from './show_card.svelte'
    export let overrideRemove = null;
    export let overrideAdd = null;
    export let cards = [];
    export let isSearchByName = true;
    export let showSelected = true;
</script>

{#if cards != null && cards.length > 0}
    {#each cards as card}
        {#if overrideAdd != null && overrideRemove != null}
            <ShowCard card={card} isSearchByName={isSearchByName} add={overrideAdd} remove={overrideRemove} showSelected={showSelected}/>
        {:else}
            <ShowCard card={card} isSearchByName={isSearchByName} showSelected={showSelected}/>
        {/if}
    {/each}
{/if}