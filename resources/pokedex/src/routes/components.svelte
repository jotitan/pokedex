<script>
    import CardStack from "../components/plateau/card_stack.svelte"
    let counter = 0;
</script>

<h1>Liste composants utiles</h1>

<button on:click={()=>counter--}>-</button>
<button on:click={()=>counter++}>+</button>

<CardStack leftCards={counter}/>
<CardStack leftCards={counter+5} singleCards="false"/>