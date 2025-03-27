<script lang="ts">
    import type { PlayerObj } from "../../stores/sessionStore";
    import { derived, get } from "svelte/store";
    import ClassicCard from "../../components/Cards/ClassicCard.svelte";
    import CardDisplay from "../../components/Game/CardDisplay.svelte";
    import { sessionStore } from "../../stores/sessionStore";
    import OpponentDisplay from "../../components/Game/OpponentDisplay.svelte";

    let maxRotationDeg = 20;
    let centerDistancePx = 200;
    let cardWidth = 100;
    let cardHeight = 150;
    let cardComponent = ClassicCard;

    let opponents = derived(sessionStore.store, ($store) => $store.players.filter((e) => e.PlayerId != sessionStore.getUserId()));
    let playerActive = derived(sessionStore.store, ($store) => ($store.playerStates[$store.userId ?? ""] ?? "").Active ?? false);
    let perSide = 0;
    $: perSide = Math.ceil($opponents.length / 3);
</script>

{#snippet OpponentCards(players: PlayerObj[], rotationDeg: number)}
    {#each players as player}
        <OpponentDisplay {cardComponent} {cardHeight} {cardWidth} {rotationDeg} {player} state={sessionStore.getPlayerState(player.PlayerId)} {centerDistancePx} {maxRotationDeg} />
    {/each}
{/snippet}

<div class="game relative grid h-full grid-rows-[1fr_2fr_1fr]">
    <div class="top flex justify-center items-center flex-row gap-20">{@render OpponentCards($opponents.slice(0, perSide), 0)}</div>
    <div class="middle grid grid-cols-[1fr_300px_1fr]">
        <div class="left flex justify-center items-center flex-col gap-20">
            {@render OpponentCards($opponents.slice(perSide, perSide * 2), -90)}
        </div>
        <div class="center grid grid-cols-2">
            <div class="drawCard flex justify-center items-center">
                <button
                    class="draw"
                    on:click={() => {
                        if (get(playerActive)) sessionStore.drawCard();
                    }}
                >
                    <svelte:component this={cardComponent} width={cardWidth} height={cardHeight} data={{ Color: "black", Symbol: "special:draw" }} />
                </button>
            </div>
            <div class="cardStack flex justify-center items-center">
                <CardDisplay
                    {cardComponent}
                    {cardHeight}
                    {cardWidth}
                    canPlayCards={false}
                    canUpdateCards={true}
                    updateCard={(_, data) => {
                        if (get(playerActive)) sessionStore.updatePlayedCard(data);
                    }}
                    cards={[{ Card: undefined, PlayedBy: "", CardIndex: -1 }, ...$sessionStore.playedCards]}
                    centerDistancePx={0}
                    maxRotationDeg={10}
                />
            </div>
        </div>
        <div class="right flex justify-center items-center flex-col gap-20">
            {@render OpponentCards($opponents.slice(perSide * 2, perSide * 3), 90)}
        </div>
    </div>
    <div class="cardContainer w-full overflow-y-clip overflow-x-auto">
        <div class="ownCards w-full min-w-min box-content flex items-end justify-center" class:opacity-60={!$playerActive}>
            <CardDisplay
                {cardHeight}
                {cardWidth}
                {cardComponent}
                fullwidth
                click={(i) => {
                    if (get(playerActive)) sessionStore.playCard(i);
                }}
                centerDistancePx={centerDistancePx * 3}
                {maxRotationDeg}
                cards={$sessionStore.ownCards}
            />
        </div>
    </div>
</div>
