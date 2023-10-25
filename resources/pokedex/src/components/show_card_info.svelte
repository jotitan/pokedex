<script>
    import {getBaseUrl, getImgEnergy} from "../routes/store.js";

    export let card = {};
    console.log("=>", card)
    const getLevel = c => {
        switch (c.Pokemon.Details.Level) {
            case 0 :
                return "Base";
            default:
                return `Niveau ${c.Pokemon.Details.Level} (évolution de ${c.Pokemon.Details.EvolutionOf})`;
        }
    }
    const getImg = c => `${getBaseUrl()}/card/image?img=${c.Img.substring(c.Img.lastIndexOf("/") + 1)}`;
</script>

{#if card != null}
    <div class="detail_card">
        <h2>{card.SubName !== "" ? card.SubName : card.Name} ({card.Extension} - {card.Number})</h2>
        <div style="display: inline-block; width: 245px">
            <img src={getImg(card)} alt="icon"/>
        </div>
        <div style="display: inline-block; vertical-align:top; margin-left:10px; width:590px;">
            <div>
                <img alt={card.Pokemon.Details.TypePokemon} src={getImgEnergy(card.Pokemon.Details.TypePokemon)}
                     class="energy_attack"/>
                {getLevel(card)} - {card.Pokemon.Details.PV} PV
            </div>
            {#if card.Pokemon.Weekness.Energy !== ""}
                <div>
                    <span class="label">Faiblesse : </span>
                    <img alt={card.Pokemon.Weekness.Energy} src={getImgEnergy(card.Pokemon.Weekness.Energy)}
                         class="energy_attack"/>
                    x {card.Pokemon.Weekness.Factor}
                </div>
            {/if}
            {#if card.Pokemon.Resistance.Energy !== ""}
                <div>
                    <span class="label">Résistance : </span>
                    <img alt={card.Pokemon.Resistance.Energy} src={getImgEnergy(card.Pokemon.Resistance.Energy)}
                         class="energy_attack"/>
                    - {card.Pokemon.Resistance.Value}
                </div>
            {/if}
            <div><span class="label">Cout de retraite : </span>{card.Pokemon.Retirement.Nb}</div>
            {#if card.Pokemon.Attacks != null}
                <h3 style="text-align:center">{card.Pokemon.Attacks.length}
                    attaque{card.Pokemon.Attacks.length > 1 ? "s" : ""} </h3>

                {#each card.Pokemon.Attacks as attack}
                    <div class="attack">
                        {#each attack.Energies as energy}
                            <img alt={energy} src={getImgEnergy(energy)} class="energy_attack"/>
                        {/each}
                        <span class="label"> {attack.Name} (+{attack.Cost})</span> :
                        {attack.Description}
                    </div>
                {/each}
            {/if}
        </div>
    </div>
{/if}

<style>
    .attack {
        margin-bottom: 15px;
    }

    .label {
        font-weight: 800;
    }

    .energy_attack {
        width: 22px;
        height: 22px;
        vertical-align: text-bottom;
    }
</style>