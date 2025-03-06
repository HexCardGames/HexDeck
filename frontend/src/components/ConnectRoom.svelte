<script lang="ts">
    import { onMount } from "svelte";
    import {
        Button,
        Spinner,
        InputAddon,
        ButtonGroup,
        Helper,
    } from "flowbite-svelte";
    import { _ } from "svelte-i18n";
    import { loading, join_error, create_error, rejoinRoomCode, rejoinRoomSessionData, requestJoinRoom, requestCreateRoom, joinSession, checkSessionData } from "../stores/roomStore";

    let joinRoomId = "";
    let inputRef: HTMLInputElement | null = null;

    function formatInput(event: any) {
        let rawValue = event.target.value.replace(/\D/g, "");
        join_error.set(false);

        if (rawValue.length > 6) {
            rawValue = rawValue.slice(0, 6);
        }

        let formattedValue = rawValue
            .replace(/(\d{3})(\d{0,3})/, "$1-$2")
            .trim();

        joinRoomId = formattedValue;

        if (joinRoomId.length > 6) {
            requestJoinRoom(joinRoomId);
        }
    }

    function handleKeyDown(event: any) {
        if (event.key === "Backspace") {
            join_error.set(false);

            let cursorPosition = event.target.selectionStart;

            if (cursorPosition === 4) {
                joinRoomId = joinRoomId.slice(0, 2);
                event.preventDefault();
            }
        }
    }

    function focusInput() {
        inputRef?.focus();
    }

    onMount(() => {
        focusInput();
        checkSessionData();
    });
</script>

<div class="w-full max-w-md">
    {#if $rejoinRoomSessionData?.sessionToken}
        <div class="mb-6">
            <Button
                class="bg-primary-400 focus:bg-primary-100 dark:bg-primary-800 focus:dark:bg-primary-700 focus:ring-0 text-dark w-full overflow-hidden rounded-xl border-1 border-primary-800"
                size="lg"
                disabled={$loading && $loading != "create"}
                on:click={() =>
                    joinSession(
                        $rejoinRoomSessionData?.sessionToken,
                        $rejoinRoomSessionData?.userId,
                    )}
            >
                {#if $loading == "create"}
                    <Spinner
                        class="text-primary-350 dark:text-primary-200 me-3"
                        size="4"
                    />
                {/if}
                {$_("landing_page.connect_room.rejoin_last_room")}
            </Button>
            {#if $create_error}
                <Helper class="mt-2">
                    <span class="text-red-900 dark:text-red-300 text-sm">
                        {$_(`error_messages.${$create_error}`, {
                            default: $_("error_messages.error_message", {
                                values: { error_message: $create_error },
                            }),
                        })}
                    </span>
                </Helper>
            {/if}
        </div>
        <div class="w-full text-center mb-4">
            {$_("landing_page.connect_room.or")}
        </div>
    {/if}
    {#if $rejoinRoomCode}
        <div class="mb-6">
            <Button
                class="bg-primary-400 focus:bg-primary-100 dark:bg-primary-800 focus:dark:bg-primary-700 focus:ring-0 text-dark w-full overflow-hidden rounded-xl border-1 border-primary-800"
                size="lg"
                disabled={$loading && $loading != "create"}
                on:click={() => requestJoinRoom($rejoinRoomCode)}
            >
                {#if $loading == "create"}
                    <Spinner
                        class="text-primary-350 dark:text-primary-200 me-3"
                        size="4"
                    />
                {/if}
                {$_("landing_page.connect_room.join_last_room")}
            </Button>
            {#if $create_error}
                <Helper class="mt-2">
                    <span class="text-red-900 dark:text-red-300 text-sm">
                        {$_(`error_messages.${$create_error}`, {
                            default: $_("error_messages.error_message", {
                                values: { error_message: $create_error },
                            })})}
                    </span>
                </Helper>
            {/if}
        </div>
        <div class="w-full text-center mb-4">
            {$_("landing_page.connect_room.or")}
        </div>
    {/if}
    <div class="mb-4">
        <ButtonGroup
            class="w-full overflow-hidden rounded-xl border-1 border-primary-800"
            size="sm"
        >
            <InputAddon
                class="w-10 text-center bg-primary-400 dark:bg-primary-800"
            >
                {#if $loading == "join"}
                    <div class="w-full flex justify-center">
                        <Spinner
                            class="text-primary-350 dark:text-primary-200"
                            size="4.2"
                        />
                    </div>
                {:else}
                    <span class="w-full"> # </span>
                {/if}
            </InputAddon>
            <input
                class="bg-primary-200 px-4 focus:bg-primary-100 placeholder:text-gray-500 dark:placeholder:text-gray-300 mr-md dark:bg-primary-700 dark:focus:bg-primary-600 w-full border-0 focus:outline-none focus:ring-0 h-12"
                placeholder={$_("landing_page.connect_room.enter_room_code")}
                class:cursor-not-allowed={$loading}
                class:opacity-50={$loading}
                disabled={!!$loading}
                bind:this={inputRef}
                bind:value={joinRoomId}
                on:input={formatInput}
                on:keydown={handleKeyDown}
            />
        </ButtonGroup>
        {#if $join_error}
            <Helper class="mt-2">
                <span class="text-red-900 dark:text-red-300 text-sm">
                    {$_(`error_messages.${$join_error}`, {
                        default: $_("error_messages.error_message", {
                            values: { error_message: $join_error },
                        })}
                    )}
                </span>
            </Helper>
        {/if}
    </div>
    <div class="w-full text-center mb-4">
        {$_("landing_page.connect_room.or")}
    </div>
    <div class="mb-6">
        <Button
            class="bg-primary-400 focus:bg-primary-100 dark:bg-primary-800 focus:dark:bg-primary-700 focus:ring-0 text-dark w-full overflow-hidden rounded-xl border-1 border-primary-800"
            size="lg"
            disabled={$loading && $loading != "create"}
            on:click={requestCreateRoom}
        >
            {#if $loading == "create"}
                <Spinner
                    class="text-primary-350 dark:text-primary-200 me-3"
                    size="4"
                />
            {/if}
            {$_("landing_page.connect_room.create_a_room")}
        </Button>
        {#if $create_error}
            <Helper class="mt-2">
                <span class="text-red-900 dark:text-red-300 text-sm">
                    {$_(`error_messages.${$create_error}`, {
                        default: $_("error_messages.error_message", {
                            values: { error_message: $create_error },
                        })})}
                </span>
            </Helper>
        {/if}
    </div>
</div>
