<script>
    export let path = "";
    export let onclick = ()=>{};
    export let title = "";
    let selected = false;

    const click = ()=> {
        selected = !selected;
        onclick(selected);
    }
</script>

<style>
    img {
        width:32px;
        height:32px;
        filter:grayscale(100%);
    }
    img:hover {
        filter:grayscale(0%);
        cursor:pointer;
    }
    img.selected{
        filter:grayscale(0%);
    }
</style>

<img src={path} on:click={click} class={selected?'selected':''} title={title} alt={title}/>