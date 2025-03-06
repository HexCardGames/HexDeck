<script lang="ts">
    import { onDestroy } from "svelte";
    import { sessionStore } from "../stores/sessionStore";
    import { ButtonGroup, InputAddon, Spinner } from "flowbite-svelte";
    import { _ } from "svelte-i18n";

    export let playerId: string = sessionStore.getUserId() || "";

    let playerName: string = sessionStore.getUser(playerId)?.Username || "ERROR USERNAME NOT FOUND";
    let isLoading: boolean = false;
    let debounceTimer: NodeJS.Timeout | null = null;
    let loadingTimer: NodeJS.Timeout | null = null;

    let inputRef: HTMLInputElement | null = null;

    function onInput(event: Event) {
        const newName = (event.target as HTMLInputElement).value;
        playerName = newName;

        if (debounceTimer) clearTimeout(debounceTimer);
        if (loadingTimer) {
            clearTimeout(loadingTimer);
            isLoading = false;
        }

        // Start loading spinner after 0.8s
        loadingTimer = setTimeout(() => {
            isLoading = true;
        }, 800);

        // Start a debounce timer (3s)
        debounceTimer = setTimeout(() => {
            sessionStore.renamePlayer(playerId, newName);
            isLoading = false;
            unfocusInput()
        }, 1800);
    }

    function focusInput() {
        console.log(focusInput);
        inputRef?.focus();
    }

    function unfocusInput() {
        console.log(focusInput);
        inputRef?.blur();
    }

    onDestroy(() => {
        if (debounceTimer) clearTimeout(debounceTimer);
        if (loadingTimer) clearTimeout(loadingTimer);
    });
</script>

<div
    class="group w-xs mx-auto text-dark bg-gray-100 dark:bg-gray-900 focus-within:bg-gray-50 dark:focus-within:bg-gray-800 backdrop-blur-lg border border-black/20 dark:border-white/20 shadow-lg p-4 rounded-2xl flex items-center justify-between transition-all"
    role="button"
    tabindex="0"
    on:click={focusInput}
    on:keydown={(event) => { if (event.key === 'Enter' || event.key === ' ') focusInput(); }}>
    <div class="grid justify-items-start w-full">
        {#if playerId == sessionStore.getUserId()}
        <span class="text-sm">{$_("lobby.rename_yourself")}</span>
        {:else}
        <span class="text-sm">{$_("lobby.rename_player")}</span>
        {/if}

        <!-- Rename Player -->
        <div class="w-full">
            <input
                class="text-black w-full dark:text-white mr-md w-full border-0 focus:outline-none focus:ring-0 h-8 bg-transparent"
                bind:value={playerName}
                on:input={onInput}
                bind:this={inputRef}
            />
        </div>
    </div>
    {#if isLoading}
    <div class="">
        <Spinner
            class="text-primary-350 w-8 h-8 dark:text-primary-200"
        />
    </div>
{/if}
</div>
