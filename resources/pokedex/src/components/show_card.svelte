<script>
    import {removeCard, updateCard, getBaseUrl, url,showCardWithNb} from '../routes/store.js';
    import {page} from "$app/stores";
    import Modal, {bind} from 'svelte-simple-modal';
    import {writable} from "svelte/store";
    import InfoPanel from "../components/show_card_info.svelte";
    const info = writable(null);
    export let card = {};
    export let showSelected = true;
    export let isSearchByName = true;
    const useCache = true;
    showCardWithNb.subscribe(v=>{
        card.Hide=v && !card.Nb > 0;
        card = card;
    })
    const toggle = () => card.Nb === 0 ? doAdd(card):()=>{};//doRemove(card);
    url.update(()=>$page.url)

    const reduce = (e)=> {
        e.stopPropagation()
        doRemove(card,false)
    }

    const increase = e => {
        e.stopPropagation()
        doAdd(card)
    }

    const doRemove = (c,all = true)=> {
        remove(c,all).then(()=>card=card)
    }

    const doAdd = c => {
        add(c).then(()=>card=card)
    }

    // Can be override
    export let remove = (c,all = true) => fetch(`${getBaseUrl()}/card?all=${all}&link=${c.Link}`,{method:'DELETE'}).then(d=>d.json()).then(d=>{
        c.Nb = d.nb;
        d.nb === 0 ? removeCard(c):updateCard(c)
        card = c;
    });

    export let add = c => fetch(`${getBaseUrl()}/card?link=${c.Link}`,{method:'POST'}).then(d=>d.json()).then(d=>{
        c.Nb = d.nb;
        updateCard(c)
        card = c;
    });
    const getImg = c => useCache ?`${getBaseUrl()}/card/image?img=${c.Img.substring(c.Img.lastIndexOf("/")+1)}`:c.Img;

    const showInfo = (e,c)=> {
        e.stopPropagation()
        info.set(bind(InfoPanel,{card:c}))
    }
</script>

{#if card != null && !card.Hide}
    <div class="card" >
        {#if card.Nb > 0}
            <span class="counter">
                <span class="counter">
                    <button on:click={reduce}>-</button>
                    <span class="nb-card">{card.Nb}</span>
                    <button on:click={increase}>+</button>
                </span>
            </span>
        {/if}
        <img src={getImg(card)} class={card.Nb > 0 && showSelected? "selected":""} on:click={toggle} alt="pokemon card"/>
        <Modal show={$info} >
            <span class="extension" title={isSearchByName?'':card.Extension} on:click={e=>showInfo(e,card)} style="cursor:help">
                {isSearchByName ?card.Extension :`${card.Name}`}
            </span> / {card.Number !== 0 ? card.Number : card.Special}
        </Modal>
    </div>
{/if}

<style>
    :global(.window) {
        width:900px !important;
    }
    .counter {
        position:absolute;
        text-align:center;
        top:320px;
        width:245px;
        opacity: 0.95;
    }
    .counter > span {
        background-color:white;
        position:relative;
        top:0px;
        padding: 0px 5px;
    }
    .extension {
        font-weight:bold;
        margin-left:10px;
    }
    .card {
        position:relative;
        width:245px;
        display:inline-block;
        margin:10px;
        vertical-align:top;
    }

    .selected {
        outline:5px solid darkgreen;
        border-radius: 10px;

    }
    .nb-card {
        font-weight:bold;
    }
</style>
