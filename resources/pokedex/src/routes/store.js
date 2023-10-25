import {writable} from 'svelte/store';

export const cards = writable([])

export const url = writable( {} )

export const forceSearch = writable({})

export const showCardWithNb = writable(false)

export const getUrl = card => `${getBaseUrl()}/card/image?img=${card.Img.substring(card.Img.lastIndexOf("/") + 1)}`

export const getBaseUrl = () => {
    let loc = {};
    url.subscribe(d=>loc = d)
    if("localhost" === loc.hostname) {
        return "http://localhost:9002";
    }
    return loc.origin;
}

const energyImages = {
    normal:"https://static.wikia.nocookie.net/pokemongo/images/f/fb/Normal.png",
    combat:"https://static.wikia.nocookie.net/pokemongo/images/3/30/Fighting.png",
    feu:"https://static.wikia.nocookie.net/pokemongo/images/3/30/Fire.png",
    "électrique":"https://static.wikia.nocookie.net/pokemongo/images/2/2f/Electric.png",
    eau:"https://static.wikia.nocookie.net/pokemongo/images/9/9d/Water.png",
    plante:"https://static.wikia.nocookie.net/pokemongo/images/c/c5/Grass.png",
    psy:"https://static.wikia.nocookie.net/pokemongo/images/2/21/Psychic.png",
    dragon:"https://static.wikia.nocookie.net/pokemongo/images/c/c7/Dragon.png",
    fee:"https://static.wikia.nocookie.net/pokemongo/images/4/43/Fairy.png",
    "obscurité":"https://static.wikia.nocookie.net/pokemongo/images/0/0e/Dark.png",
    "métal":"https://static.wikia.nocookie.net/pokemongo/images/c/c9/Steel.png",
}

export const getImgEnergy = name => {
    return energyImages[name];
}

export const updateCards = () => fetch(`${getBaseUrl()}/card`).then(d=>d.json()).then(c => cards.update(()=>c))

export const removeCard = card => cards.update(cds => cds.filter(c=>c.Link !== card.Link));

export const updateCard = card => cards.update(cds => {
    const newCds = [...cds];
    const idx = newCds.findIndex(c=>c.Link === card.Link)
    if(idx === -1){
        newCds.push(card);
    }else {
        newCds[idx] = card;
    }
    return newCds;
});