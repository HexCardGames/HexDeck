<script lang="ts">
    import { ChartNoAxesCombined } from "lucide-svelte";
    import { onDestroy, onMount } from "svelte";
    import { _ } from "svelte-i18n";

    let stats = {
        online_player_count: null,
        current_game_rooms: null,
        games_played: null,
    }

    async function getStats() {
        try {
            const res = await fetch("/api/stats");
            const resJson = await res.json();
            stats.online_player_count = resJson?.OnlinePlayerCount
            stats.current_game_rooms = resJson?.RunningGames
            stats.games_played = resJson?.TotalGamesPlayed
        } catch {}
    }

    let getStateInterval: any = undefined;

    onMount(() => {
        getStats();
        // Request stats update every 10s
        getStateInterval = setInterval(getStats, 10 * 60 * 1000);
    })

    onDestroy(() => {
        clearInterval(getStateInterval);
    })
</script>

<div class="p-4 bg-primary-50 dark:bg-primary-950 rounded-xl grid content-start justify-items-center w-3xs text-center space-y-2 border-1 border-primary-200 dark:border-primary-800">
    <ChartNoAxesCombined size="48px" />
    <h4 class="text-xl font-semibold">{$_("landing_page.stats_container.title")}</h4>
    <!-- content div -->
    <div class="grid justify-items-start text-start">
        <span>{$_("landing_page.stats_container.online_player_count", { values: { count: stats.online_player_count ?? "..." } })}</span>
        <span>{$_("landing_page.stats_container.current_game_rooms", { values: { count: stats.current_game_rooms ?? "..." } })}</span>
        <span>{$_("landing_page.stats_container.games_played", { values: { count: stats.games_played ?? "..." } })}</span>
    </div>
</div>