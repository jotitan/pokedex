<script>
    import {page} from '$app/stores';
    import CardStack from "./card_stack.svelte"
    import BlockCard from "./show_card.svelte"
    import ShowCards from "./show_cards.svelte"
    import {updateCards, cards, url} from "../../routes/store.js";

    export let isTop = false;

    let myCards;
    url.update(()=>$page.url)
    cards.subscribe(c => myCards=c)
    updateCards()
</script>

<div style="position:absolute;left:30px">
    <CardStack leftCards={6} singleCards="false"/>
</div>

<div style={`position:absolute;left:calc(50% - 120px); top:${isTop?'calc(49vh - 240px)':'20px'}`}>
    <BlockCard card={myCards[0]}/>
</div>

<div style={`position:absolute;right:30px;top:${isTop?'calc(49vh - 210px)':'20px'}`}>
    <CardStack leftCards={50}/>
</div>

<div style={`position:absolute;left:330px;bottom:10px;`}>
    <ShowCards cards={[myCards[1],myCards[15],myCards[10]]}/>
</div>

