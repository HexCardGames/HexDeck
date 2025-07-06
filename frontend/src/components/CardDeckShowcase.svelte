<script lang="ts">
    import CardDisplay from "./Game/CardDisplay.svelte";

    export let cardComponent;
    export let cardDeckId: number;
    export let cardWidth: number;
    export let cardHeight: number;
    export let centerDistancePx: number;
    export let maxRotationDeg: number;

    const SHOWCASED_CARDS: any = {
        0: [
            { CanPlay: true, Card: { Symbol: "1", Color: "red" } },
            { CanPlay: true, Card: { Symbol: "9", Color: "green" } },
            { CanPlay: true, Card: { Symbol: "action:draw_2", Color: "blue" } },
            { CanPlay: true, Card: { Symbol: "action:skip", Color: "yellow" } },
            { CanPlay: true, Card: { Symbol: "action:draw_4", Color: "black" } },
        ],
        1: [
            { CanPlay: true, Card: { Symbol: "1", Color: "blue" } },
            { CanPlay: true, Card: { Symbol: "F", Color: "green" } },
            { CanPlay: true, Card: { Symbol: "action:shuffle", Color: "yellow" } },
            { CanPlay: true, Card: { Symbol: "action:skip", Color: "purple" } },
            { CanPlay: true, Card: { Symbol: "action:draw", Color: "rainbow" } },
        ],
    };
</script>

<CardDisplay canPlayCards={true} canUpdateCards={false} {cardComponent} cards={SHOWCASED_CARDS[cardDeckId]} {cardWidth} {cardHeight} {centerDistancePx} {maxRotationDeg} />
