<script>
    export let pokemon;
</script>

{#if pokemon != null}
    <div>
        <h2>{pokemon.name}</h2>
        <img src={pokemon.ThumbnailImage}/>
    </div>



{/if}



