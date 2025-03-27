<script lang="ts">
    import EndScreen from "./../views/Game/EndScreen.svelte";
    import Main from "./../views/Game/Main.svelte";
    import Lobby from "../views/Game/Lobby.svelte";
    import { _ } from "svelte-i18n";
    import { onMount } from "svelte";
    import { GameState, sessionStore } from "../stores/sessionStore";
    import { Spinner } from "flowbite-svelte";
    import { SvelteDate } from "svelte/reactivity";
    import { requestJoinRoom } from "../stores/roomStore";
    import gameStore from "../stores/gameStore";

    let maxRotationDeg = 20;
    let centerDistancePx = 200;
    let cardWidth = 100;
    let cardHeight = 150;

    onMount(async () => {
        // TODO: check if already connected to room, currently its overwriting the session
        const params = new URLSearchParams(window.location.search);
        const joinParam = params.get("join");

        if (joinParam) {
            await requestJoinRoom(joinParam);
            // Maybe show message instead redirecting to / if the join was unsuccessful
        }

        if (!sessionStore.hasSessionData()) {
            console.warn("No sessionData found! Go back home.");
            window.history.replaceState({}, "", "/");
        }
        sessionStore.connect();
    });
</script>

{#if !$sessionStore.connected}
    <div class="flex flex-row w-full mt-32 h-full justify-center items-center">
        <div class="flex flex-col items-center gap-6 p-7 md:flex-row md:gap-8 rounded-2xl">
            <div>
                <Spinner size="12" class="text-primary-100" />
            </div>
            <div class="grid items-center text-center md:items-start">
                <span class="text-2xl font-medium">
                    {$_("game_screen.loading")}
                </span>
                <span class="font-medium text-sky-500">
                    {$sessionStore.players?.find((player) => player.PlayerId == $sessionStore.userId)?.Username}
                </span>
                <span class="flex gap-2 font-medium text-gray-600 dark:text-gray-400">
                    <span>{new SvelteDate().toLocaleString()}</span>
                </span>
            </div>
        </div>
    </div>
{:else}
    {#if $sessionStore.gameState == GameState.Lobby}
        <div>
            <!-- Lobby and player list -->
            <Lobby {cardWidth} {cardHeight} {centerDistancePx} {maxRotationDeg} />
        </div>
    {/if}

    {#if $sessionStore.gameState == GameState.Running}
        <div class="size-full">
            {#if $gameStore.isLobbyOverlayShown}
                <div class="absolute inset-0 z-10 bg-white/30 dark:bg-black/30 backdrop-blur-sm mt-24">
                    <Lobby {cardWidth} {cardHeight} {centerDistancePx} {maxRotationDeg} />
                </div>
            {/if}
            <!-- Running game -->
            <Main {cardWidth} {cardHeight} {centerDistancePx} {maxRotationDeg} />
        </div>
    {:else if $sessionStore.gameState == GameState.Ended}
        <EndScreen />
    {/if}
{/if}
