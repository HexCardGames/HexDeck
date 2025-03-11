<script lang="ts">
    import type { PlayerObj, PlayerStateObj } from "../../stores/sessionStore";
    import CardDisplay from "./CardDisplay.svelte";

    export let player: PlayerObj;
    export let state: PlayerStateObj | undefined;
    export let cardComponent;
    export let cardWidth: number;
    export let cardHeight: number;

    export let centerDistancePx;
    export let maxRotationDeg;
    export let rotationDeg;
</script>

{#if state}
    <div class="opponentDisplay flex flex-col justify-center items-center" style:--rotation={`${rotationDeg}deg`}>
        <p class:font-bold={state?.Active}>{player.Username}</p>
        <CardDisplay
            {cardComponent}
            {cardWidth}
            {cardHeight}
            cards={Array(state?.NumCards).fill({ Card: undefined })}
            canPlayCards={false}
            canUpdateCards={false}
            hoverOffset={0}
            {centerDistancePx}
            {maxRotationDeg}
        />
    </div>
{/if}

<style>
    .opponentDisplay {
        transform: rotate(var(--rotation));
    }
</style>
