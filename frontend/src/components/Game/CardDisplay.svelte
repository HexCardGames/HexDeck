<script lang="ts">
    import type { Card, CardInfoObj, CardPlayedObj } from "../../stores/sessionStore";

    export let cardWidth: number;
    export let cardHeight: number;
    export let cards: CardInfoObj[] | CardPlayedObj[];
    export let cardComponent;
    export let click: (index: number) => void = () => {};
    export let updateCard: (index: number, newCard: Card) => void = () => {};
    export let canPlayCards: boolean = true;
    export let canUpdateCards: boolean = false;
    export let centerDistancePx: number;
    export let maxRotationDeg: number;
    export let fullwidth: boolean = false;
    export let hoverOffset: number = canPlayCards ? 40 : 0;

    let rotationDeg = 0;
    $: rotationDeg = Math.min(maxRotationDeg, 1.5 * cards.length);
    let maxOffset = [0, 0];
    let offset = [0, 0];
    $: if (centerDistancePx && rotationDeg) maxOffset = getCardOffset(getCardRotation(0, cards.length, maxRotationDeg));
    $: if (centerDistancePx && rotationDeg) offset = getCardOffset(getCardRotation(0, cards.length, rotationDeg));

    function getCardRotation(i: number, totalCards: number, maxRotationDeg: number): number {
        if (totalCards == 1) return 0;
        return -maxRotationDeg + (i / (totalCards - 1)) * 2 * maxRotationDeg;
    }

    function getCardOffset(angle: number): [number, number] {
        let slope = Math.tan(angle * (Math.PI / 180));
        let y = Math.sqrt(slope ** 2 + 1) * (centerDistancePx / (slope ** 2 + 1));
        let x = slope * y;
        return [x, centerDistancePx - y];
    }
</script>

<div
    class="cards relative flex justify-center box-content"
    class:fullwidth
    style:--height={`${maxOffset[1] + cardHeight}px`}
    style:--width={`${-offset[0] * 2 + cardWidth}px`}
    style:--hover-offset={`${hoverOffset}px`}
>
    {#each cards as cardInfo, i}
        {@const rotation = getCardRotation(i, cards.length, rotationDeg)}
        {@const position = getCardOffset(rotation)}
        <button
            class="absolute card drop-shadow-lg"
            disabled={!(cardInfo as CardInfoObj).CanPlay}
            on:click={() => click(i)}
            class:canPlayCards
            style:--rotation={`${rotation}deg`}
            style:--left={`${position[0]}px`}
            style:--top={`${position[1]}px`}
        >
            <svelte:component
                this={cardComponent}
                width={cardWidth}
                height={cardHeight}
                data={cardInfo.Card}
                canUpdateCard={canUpdateCards}
                updateCard={(newCard: Card) => {
                    updateCard(i, newCard);
                }}
            />
        </button>
    {/each}
</div>

<style>
    .cards {
        height: var(--height);
        width: var(--width);
        padding-top: var(--hover-offset);
    }
    .card {
        transform: translate(var(--left), var(--top)) rotate(var(--rotation));
        transition: 0.2s transform;
        cursor: unset;
        user-select: none;
    }
    .card.canPlayCards:enabled {
        cursor: pointer;
    }
    .card.canPlayCards:disabled {
        filter: brightness(0.4);
    }
    .card.canPlayCards:hover {
        transform: translate(var(--left), var(--top)) rotate(var(--rotation)) translate(0px, calc(0px - var(--hover-offset)));
    }
</style>
