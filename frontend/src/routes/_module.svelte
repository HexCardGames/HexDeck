<script>
    import { theme, toggleTheme } from "../stores/theme";
    import { Gamepad2, Moon, Sun, SunMoon, UsersRound } from "lucide-svelte";
    import { Tooltip, Button } from "flowbite-svelte";
    import options from "../stores/pageoptions";
    import { _ } from "svelte-i18n";
    import gameStore, { toggleLobbyOverlay } from "../stores/gameStore";
    import { GameState, sessionStore } from "../stores/sessionStore";
</script>

<header class="Header">
    <div class="Header-bg"></div>
    <div class="Header-content">
        <div class="left-header-group header-group">
            <div class="page-header-icon">
                <svelte:component this={options.page_icon} size="2.4rem" />
            </div>
            <h1 class="text-3xl">{$_("page_name")}</h1>
        </div>
        <div class="middle-header-group header-group"></div>
        <div class="right-header-group header-group gap-2">
            <!-- Theme btn -->
            <Button on:click={toggleTheme} class="!p-2 mt-2 rounded-full focus:bg-primary-700 hover:bg-primary-600 focus:ring-0" color="none">
                {#if $theme === "dark"}
                    <Moon size="2rem" />
                {:else if $theme === "light"}
                    <Sun size="2rem" />
                {:else if $theme === "system"}
                    <SunMoon size="2rem" />
                {/if}
            </Button>
            <Tooltip type="auto">
                {$_("header.theme_btn.tooltip", {
                    values: { current_theme: $_(`header.theme_btn.${$theme}`) },
                })}
            </Tooltip>

            <!-- Player list btn (ingame) -->
            {#if $sessionStore.gameState == GameState.Running}
                <Button
                    on:click={() => {
                        toggleLobbyOverlay();
                    }}
                    class="!p-2 mt-2 rounded-full focus:bg-primary-700 hover:bg-primary-600 focus:ring-0"
                    color="none"
                >
                    {#if $gameStore.isLobbyOverlayShown}
                        <Gamepad2 size="2rem" />
                    {:else}
                        <UsersRound size="2rem" />
                    {/if}
                </Button>
                {#if $gameStore.isLobbyOverlayShown}
                    <Tooltip type="auto">
                        {$_("lobby.return_to_game")}
                    </Tooltip>
                {:else}
                    <Tooltip type="auto">
                        {$_("lobby.player")}
                    </Tooltip>
                {/if}
            {/if}
        </div>
    </div>
</header>

<div class="main-container">
    <div class="page-slot">
        <slot />
    </div>
</div>

<style>
    .main-container {
        display: flex;
        flex-direction: column;
        height: 100vh;
        padding-top: 100px;
    }

    .page-slot {
        flex-grow: 1;
        overflow: hidden;
    }

    .Header {
        background: linear-gradient(180deg, var(--default-background-color) 30%, transparent 100%);
        opacity: 1;
        position: fixed;
        height: 100px;
        top: 0;
        left: 0;
        width: 100%;
        z-index: 30;
    }

    .Header-bg {
        background: linear-gradient(180deg, var(--primary) 50%, transparent 100%);
        z-index: -1;
        margin: 0px;
        position: absolute;
        width: 100%;
        height: 100%;
    }

    .Header-content {
        display: flex;
        justify-content: space-between;
        align-items: center;
        padding: 0 1rem;
    }

    .left-header-group,
    .middle-header-group,
    .right-header-group {
        display: flex;
        align-items: center;
    }

    .middle-header-group {
        flex-grow: 1;
        justify-content: center;
    }

    .right-header-group {
        margin-left: auto;
    }

    .page-header-icon {
        display: flex;
        align-items: center;
        justify-content: center;
        margin-right: 10px;
    }
</style>
