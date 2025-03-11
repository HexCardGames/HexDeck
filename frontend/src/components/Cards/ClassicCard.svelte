<script lang="ts">
    import type { Card } from "../../stores/sessionStore";

    interface ClassicCard {
        Symbol: string;
        Color: string;
    }
    export let data: ClassicCard = { Color: "", Symbol: "" };
    export let canUpdateCard: boolean = false;
    export let updateCard: (newCard: Card) => void = () => {};
    export let width: number;
    export let height: number;

    const COLOR_MAP: { [key: string]: string } = {
        red: "0#990000",
        green: "#009900",
        blue: "#000099",
        yellow: "#999900",
        black: "#000000",
        "": "#808080",
    };

    function selectColor(color: string) {
        let newCard: ClassicCard = {
            Color: color,
            Symbol: data.Symbol,
        };
        updateCard(newCard);
    }
</script>

<div class="card" style:background={COLOR_MAP[data.Color]} style:--width={`${width}px`} style:--height={`${height}px`}>
    {#if data.Symbol.length <= 2}
        <span class="symbol">{data.Symbol}</span>
    {:else}
        <span class="symbol large">{data.Symbol.split(":")[1].replace("_", " ")}</span>
    {/if}
    {#if data.Color == "black" && canUpdateCard}
        <div class="select mt-5">
            <div>
                <button
                    class="bg-red-400"
                    aria-label="Red"
                    on:click={() => {
                        selectColor("red");
                    }}
                ></button>
                <button
                    class="bg-green-400"
                    aria-label="Green"
                    on:click={() => {
                        selectColor("green");
                    }}
                ></button>
            </div>
            <div>
                <button
                    class="bg-blue-400"
                    aria-label="Blue"
                    on:click={() => {
                        selectColor("blue");
                    }}
                ></button>
                <button
                    class="bg-yellow-400"
                    aria-label="Yellow"
                    on:click={() => {
                        selectColor("yellow");
                    }}
                ></button>
            </div>
        </div>
    {/if}
</div>

<style>
    .card {
        display: flex;
        flex-direction: column;
        width: var(--width);
        height: var(--height);
        padding: 5px;
        border-radius: 8px;
    }

    .select button {
        width: 20px;
        height: 20px;
        border-radius: 30px;
    }

    .card span {
        text-align: left;
        word-wrap: anywhere;
    }

    .card .symbol {
        font-size: 25px;
        text-transform: capitalize;
    }

    .card .symbol.large {
        font-size: 17px;
    }

    button {
        cursor: pointer;
    }
</style>
