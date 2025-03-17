<script>
    import { CrownIcon, TimerIcon } from "lucide-svelte";
    import { _ } from "svelte-i18n";
    import { sessionStore } from "../../stores/sessionStore";
    import { Button } from "flowbite-svelte";
</script>

<div class="flex size-full justify-center items-center">
    <div class="flex flex-col items-center gap-2">
        {#if $sessionStore.winner == undefined}
            <TimerIcon size={64} />
            <span>{$_("end_screen.game_has_ended")}</span>
        {:else}
            <CrownIcon size={64} />
            {#if $sessionStore.winner == sessionStore.getUserId()}
                <span>{$_("end_screen.you_won_the_game")}</span>
            {:else}
                <span
                    >{$_("end_screen.player_won_the_game", {
                        values: { player_name: sessionStore.getUser($sessionStore.winner)?.Username },
                    })}</span
                >
            {/if}
        {/if}
        <div class="flex flex-row gap-3 mt-5">
            <Button
                color="alternative"
                on:click={() => {
                    sessionStore.leaveRoom();
                    window.history.replaceState({}, "/Game", "/");
                }}>{$_("end_screen.go_back")}</Button
            >
            <!-- TODO: implement rematch with the same players -->
            <Button>{$_("end_screen.play_again")}</Button>
        </div>
    </div>
</div>
