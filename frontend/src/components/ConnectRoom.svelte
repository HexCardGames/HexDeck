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
    import { sessionStore } from "/stores/sessionStore";

    let loading: "join" | "create" | false = false;
    let rejoinRoomCode = "";
    let rejoinRoomSessionData = {
        sessionToken: "",
        userId: "",
    };
    let joinRoomId = "";
    let join_error: string | false = false;
    let create_error: string | false = false;
    let inputRef: HTMLInputElement | null = null;

    function formatInput(event: any) {
        let rawValue = event.target.value.replace(/\D/g, "");
        join_error = false;

        if (rawValue.length > 6) {
            rawValue = rawValue.slice(0, 6);
        }

        let formattedValue = rawValue
            .replace(/(\d{3})(\d{0,3})/, "$1-$2")
            .trim();

        joinRoomId = formattedValue;

        if (joinRoomId.length > 6) {
            requestJoinRoom();
        }
    }

    function handleKeyDown(event: any) {
        if (event.key === "Backspace") {
            join_error = false;

            let cursorPosition = event.target.selectionStart;

            if (cursorPosition === 4) {
                // If cursor is at "-" position, delete the number before it
                joinRoomId = joinRoomId.slice(0, 2);
                event.preventDefault();
            }
        }
    }

    async function requestJoinRoom(joinCode = joinRoomId) {
        if (loading) return;
        loading = "join";
        join_error = false;

        try {
            const controller = new AbortController();
            // 5s timeout
            const timeout = setTimeout(() => controller.abort(), 5000);

            const response = await fetch(`/api/room/join`, {
                method: "POST",
                body: JSON.stringify({
                    JoinCode: joinCode.replaceAll("-", ""),
                    UsernameProposal: "UsernameProposal",
                }),
                headers: {
                    "Content-Type": "application/json",
                },
                signal: controller.signal,
            });

            clearTimeout(timeout);

            if (!response.ok) {
                const data: { StatusCode: string; Message: string } =
                    await response.json();
                // TODO i18n here on StatusCode if not use Message
                if (["invalid_join_code"].includes(data?.StatusCode)) {
                    join_error = "no_room_found";
                } else if (data?.Message) {
                    join_error = data?.Message;
                } else {
                    throw new Error("Server error");
                }
                return;
            }

            const data: {
                SessionToken: string;
                PlayerId: string;
                Username: string;
                Permissions: any;
            } = await response.json();
            const SessionToken = data.SessionToken;
            const UserId = data.PlayerId;
            joinSession(SessionToken, UserId);
        } catch (error: any) {
            if (error.name === "AbortError") {
                join_error = "timeout";
            } else {
                join_error = "request_failed";
            }
            console.error("Error joining room: ", error);
        } finally {
            loading = false;
            setTimeout(() => {
                focusInput();
            }, 50);
        }
    }

    async function requestCreateRoom() {
        if (loading) return;
        loading = "create";
        create_error = false;

        try {
            const controller = new AbortController();
            // 5s timeout
            const timeout = setTimeout(() => controller.abort(), 5000);

            const response = await fetch(`/api/room/create`, {
                method: "POST",
                body: JSON.stringify({
                    UsernameProposal: "UsernameProposal",
                }),
                headers: {
                    "Content-Type": "application/json",
                },
                signal: controller.signal,
            });

            clearTimeout(timeout);

            if (!response.ok) {
                throw new Error("Server error");
            }

            const data:
                | {
                      SessionToken: string;
                      PlayerId: string;
                      Username: string;
                      Permissions: any;
                  }
                | { error: string } = await response.json();

            if (response.ok) {
                const SessionToken = data.SessionToken;
                const UserId = data.PlayerId;
                sessionStore.connect(SessionToken, UserId);
            } else {
                create_error = data.error || "room_creation_failed";
            }
        } catch (error: any) {
            if (error.name === "AbortError") {
                create_error = "timeout";
            } else {
                create_error = String(error);
            }
            console.error("Error creating room:", error);
        } finally {
            loading = false;
        }
    }

    function focusInput() {
        inputRef?.focus();
    }

    function joinSession(sessionToken: string, userId: string) {
        try {
            sessionStore.connect(sessionToken, userId);
        } catch (error: any) {
            join_error = "request_failed";
            console.error("Error joining room session: ", error);
        } finally {
            loading = false;
            setTimeout(() => {
                focusInput();
            }, 50);
        }
    }

    async function checkSessionToken(
        sessionToken: string | undefined,
    ): Promise<boolean> {
        if (!sessionToken) return false;
        const params = new URLSearchParams({
            sessionToken: sessionToken,
        });
        const res = await fetch(`/api/check/session?${params}`);
        return res.status == 200;
    }

    async function checkJoinCode(
        joinCode: string | undefined,
    ): Promise<boolean> {
        if (!joinCode) return false;
        const params = new URLSearchParams({
            JoinCode: joinCode,
        });
        const res = await fetch(`/api/check/joinCode?${params}`);
        return res.status == 200;
    }

    async function checkSessionData() {
        const currentSessionData: {
            sessionToken?: string;
            userId?: string;
            joinCode?: string;
        } = JSON.parse(localStorage.getItem("currentSessionIds") || "{}");
        if (await checkSessionToken(currentSessionData.sessionToken)) {
            // Session is still valid
            rejoinRoomSessionData = currentSessionData as any;
            return;
        }
        if (await checkJoinCode(currentSessionData.joinCode)) {
            // joinCode  is still valid
            rejoinRoomCode = currentSessionData.joinCode as string;
            return;
        }

        const lastSessionData: {
            sessionToken?: string;
            userId?: string;
            joinCode?: string;
        } = JSON.parse(localStorage.getItem("lastSessionIds") || "{}");

        if (await checkSessionToken(lastSessionData.sessionToken)) {
            // Session is still valid
            rejoinRoomSessionData = lastSessionData as any;
            return;
        }
        if (await checkJoinCode(lastSessionData.joinCode)) {
            // joinCode  is still valid
            rejoinRoomCode = lastSessionData.joinCode as string;
            return;
        }
    }

    // focus room code input on mount
    onMount(() => {
        focusInput();
        checkSessionData();
    });
</script>

<div class="w-full max-w-md">
    {#if rejoinRoomSessionData?.sessionToken}
        <div class="mb-6">
            <Button
                class="bg-primary-400 focus:bg-primary-100 dark:bg-primary-800 focus:dark:bg-primary-700 focus:ring-0 text-dark w-full overflow-hidden rounded-xl border-1 border-primary-800"
                size="lg"
                disabled={loading && loading != "create"}
                on:click={() =>
                    joinSession(
                        rejoinRoomSessionData?.sessionToken,
                        rejoinRoomSessionData?.userId,
                    )}
            >
                {#if loading == "create"}
                    <Spinner
                        class="text-primary-350 dark:text-primary-200 me-3"
                        size="4"
                    />
                {/if}
                {$_("landing_page.connect_room.rejoin_last_room")}
            </Button>
            {#if create_error}
                <Helper class="mt-2">
                    <span class="text-red-900 dark:text-red-300 text-sm">
                        {$_(`error_messages.${create_error}`, {
                            default: $_("error_messages.error_message", {
                                values: { error_message: create_error },
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
    {#if rejoinRoomCode}
        <div class="mb-6">
            <Button
                class="bg-primary-400 focus:bg-primary-100 dark:bg-primary-800 focus:dark:bg-primary-700 focus:ring-0 text-dark w-full overflow-hidden rounded-xl border-1 border-primary-800"
                size="lg"
                disabled={loading && loading != "create"}
                on:click={() => requestJoinRoom(rejoinRoomCode)}
            >
                {#if loading == "create"}
                    <Spinner
                        class="text-primary-350 dark:text-primary-200 me-3"
                        size="4"
                    />
                {/if}
                {$_("landing_page.connect_room.join_last_room")}
            </Button>
            {#if create_error}
                <Helper class="mt-2">
                    <span class="text-red-900 dark:text-red-300 text-sm">
                        {$_(`error_messages.${create_error}`, {
                            default: $_("error_messages.error_message", {
                                values: { error_message: create_error },
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
    <div class="mb-4">
        <ButtonGroup
            class="w-full overflow-hidden rounded-xl border-1 border-primary-800"
            size="sm"
        >
            <InputAddon
                class="w-10 text-center bg-primary-400 dark:bg-primary-800"
            >
                {#if loading == "join"}
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
                class="bg-primary-200 px-4 focus:bg-primary-100 text-black dark:text-white mr-md dark:bg-primary-700 dark:focus:bg-primary-600 w-full border-0 focus:outline-none h-12"
                placeholder={$_("landing_page.connect_room.enter_room_code")}
                class:cursor-not-allowed={loading}
                class:opacity-50={loading}
                disabled={!!loading}
                bind:this={inputRef}
                bind:value={joinRoomId}
                on:input={formatInput}
                on:keydown={handleKeyDown}
            />
        </ButtonGroup>
        {#if join_error}
            <Helper class="mt-2">
                <span class="text-red-900 dark:text-red-300 text-sm">
                    {$_(`error_messages.${join_error}`, {
                        default: $_("error_messages.error_message", {
                            values: { error_message: join_error },
                        }),
                    })}
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
            disabled={loading && loading != "create"}
            on:click={requestCreateRoom}
        >
            {#if loading == "create"}
                <Spinner
                    class="text-primary-350 dark:text-primary-200 me-3"
                    size="4"
                />
            {/if}
            {$_("landing_page.connect_room.create_a_room")}
        </Button>
        {#if create_error}
            <Helper class="mt-2">
                <span class="text-red-900 dark:text-red-300 text-sm">
                    {$_(`error_messages.${create_error}`, {
                        default: $_("error_messages.error_message", {
                            values: { error_message: create_error },
                        }),
                    })}
                </span>
            </Helper>
        {/if}
    </div>
</div>
