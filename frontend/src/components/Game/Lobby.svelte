<script lang="ts">
    import RenamePlayer from "./../RenamePlayer.svelte";
    import {
        Table,
        TableBody,
        TableBodyCell,
        TableBodyRow,
        TableHead,
        TableHeadCell,
        TableSearch,
        Badge,
        Button,
        Modal,
        Popover,
        Tooltip,
    } from "flowbite-svelte";
    import {
        CircleArrowOutUpLeft,
        Copy,
        AlertCircle,
        UserX,
        Play,
        TextCursorInput,

        Gamepad2

    } from "lucide-svelte";
    import { _ } from "svelte-i18n";
    import { GameState, sessionStore } from "../../stores/sessionStore";
    import { toggleLobbyOverlay } from '../../stores/gameStore';

    let copied = false;
    let showLeaveModal = false;

    $: players = $sessionStore.players;

    let searchQuery = "";
    let rename_player = "";
    let showRenameModal = false;
    let kick_player = "";
    let showKickModal = false;

    function filteredPlayers() {
        return players.filter((player) =>
            player.Username.toLowerCase().includes(searchQuery.toLowerCase()),
        );
    }

    function insert(str: string, index: number, value: string) {
        return str.slice(0, index) + value + str.slice(index);
    }

    function copyGameCodeToClipboard() {
        navigator.clipboard
            .writeText(insert($sessionStore.joinCode || "000000", 3, "-"))
            .then(() => {
                copied = true;
                setTimeout(() => (copied = false), 2000);
            });
    }

    function copyGameLinkToClipboard() {
        navigator.clipboard
            .writeText(
                `${window.location.origin}/Game?join=${$sessionStore.joinCode}`,
            )
            .then(() => {
                copied = true;
                setTimeout(() => (copied = false), 2000);
            });
    }

    function leaveRoom() {
        sessionStore.leaveRoom();
        showLeaveModal = false;
    }
</script>

<!-- Modal: Confirm Leave Room -->
<Modal
    bind:open={showLeaveModal}
    size="md"
    backdropClass="fixed inset-0 z-40 bg-gray-900 bg-black/50 dark:bg-black/80 backdrop-opacity-50"
    autoclose
    outsideclose
>
    <div class="text-center">
        <AlertCircle
            class="mx-auto mb-4 text-gray-400 w-12 h-12 dark:text-gray-200"
        />
        <h3 class="mb-5 text-lg font-normal text-gray-500 dark:text-gray-400">
            {$_("lobby.confirm_leave_message")}
        </h3>
        <Button
            on:click={() => (showLeaveModal = false)}
            color="alternative"
            class="hover:text-dark hover:bg-gray-100"
            >{$_("lobby.cancel")}</Button
        >
        <Button on:click={leaveRoom} color="red" class="me-2"
            >{$_("lobby.confirm_leave")}</Button
        >
    </div>
</Modal>

<!-- Modal: Rename Player -->
<Modal
    bind:open={showRenameModal}
    size="md"
    backdropClass="fixed inset-0 z-40 bg-gray-900 bg-black/50 dark:bg-black/80 backdrop-opacity-50"
    autoclose
    outsideclose
>
    <div class="text-center">
        <TextCursorInput
            class="mx-auto mb-4 text-gray-400 w-12 h-12 dark:text-gray-200"
        />
        <RenamePlayer playerId={rename_player} />
    </div>
</Modal>

<!-- Modal: Confirm Kick Player -->
<Modal
    bind:open={showKickModal}
    size="md"
    backdropClass="fixed inset-0 z-40 bg-gray-900 bg-black/50 dark:bg-black/80 backdrop-opacity-50"
    autoclose
    outsideclose
>
    <div class="text-center">
        <UserX
            class="mx-auto mb-4 text-gray-400 w-12 h-12 dark:text-gray-200"
        />
        <h3 class="mb-5 text-lg font-normal text-gray-500 dark:text-gray-400">
            {$_("lobby.confirm_kick_player_message", {
                values: { player_name: sessionStore.getUser(kick_player)?.Username || 'Name not found' },
            })}
        </h3>
        <Button
            on:click={() => (showLeaveModal = false)}
            color="alternative"
            class="hover:text-dark hover:bg-gray-100"
            >{$_("lobby.cancel")}</Button
        >
        <Button
            on:click={() => {
                sessionStore.kickPlayer(kick_player);
            }}
            color="red"
            class="me-2">{$_("lobby.confirm_kick_player")}</Button
        >
    </div>
</Modal>

<!-- Leave Room Button -->
<Button
    color="none"
    class="sm:absolute m-2 border-2 border-gray-500 dark:border-gray-300 hover:bg-gray-500 dark:hover:bg-gray-300 hover:text-white dark:hover:text-black rounded-full text-gray-500 dark:text-gray-300"
    on:click={() => {
        showLeaveModal = true;
    }}
>
    <CircleArrowOutUpLeft class="mr-2" />
    <span>{$_("lobby.leave_game")}</span>
</Button>

<!-- Return to game Button -->
{#if sessionStore.getState().gameState !== GameState.Lobby}
<Button
    color="none"
    class="sm:absolute m-2 right-0 border-2 border-gray-500 dark:border-gray-300 hover:bg-gray-500 dark:hover:bg-gray-300 hover:text-white dark:hover:text-black rounded-full text-gray-500 dark:text-gray-300"
    on:click={() => {
        toggleLobbyOverlay();
    }}
>
    <span>{$_("lobby.return_to_game")}</span>
    <Gamepad2 class="ml-2" />
</Button>
{/if}

<!-- Game Status -->
<div class="text-center p-6 w-full">
    <span>{$_("game_status.game_status")}
        <Badge color="dark">
            {$_(`game_status.${GameState[$sessionStore.gameState].toLowerCase()}`)}
        </Badge>
        </span
    >
</div>

<!-- TODO Grid not fully responsive -->
<div class="grid md:grid-flow-col grid-flow-row justify-center mt-6 mb-2 gap-4">
    <!-- Rename (This) Player -->
    {#if sessionStore.isConnected()}
        <RenamePlayer />
    {/if}

    <!-- Copy Join Code Button -->
     <!-- TODO add Streamer mode (hide room code) here -->
    {#if sessionStore.getState().gameState === GameState.Lobby}
    <Button
        id="b1"
        type="button"
        class="w-xs mx-auto text-dark bg-primary-200 dark:bg-primary-900 hover:bg-primary-200 dark:hover:bg-primary-900 backdrop-blur-lg border border-black/20 dark:border-white/20 shadow-lg p-4 rounded-2xl flex items-center justify-between transition-all cursor-pointer"
        on:click={() => {
            copyGameCodeToClipboard();
        }}
    >
        <div class="grid justify-items-start">
            <span class="text-sm">{$_("lobby.room_join_code")}</span>
            <div class="relative">
                <span
                    class="text-xl font-semibold tracking-widest select-none transition-opacity duration-300 opacity-0"
                    class:opacity-100={!copied}
                    >{insert($sessionStore.joinCode || "000000", 3, "-")}
                </span>
                <span
                    class="absolute left-0 text-xl font-semibold tracking-widest select-none transition-opacity duration-300 opacity-0"
                    class:opacity-100={copied}
                >
                    {$_("lobby.copied")}
                </span>
            </div>
        </div>
        <Copy />
    </Button>
    <Popover
        class="text-sm max-w-screen font-light z-100"
        triggeredBy="#b1"
        placement="bottom"
    >
        <div class="grid gap-2">
            <Button
                on:click={() => {
                    copyGameCodeToClipboard();
                }}
            >
                {$_("lobby.copy_code")}
            </Button>
            <Button
                on:click={() => {
                    copyGameLinkToClipboard();
                }}
            >
                {$_("lobby.copy_join_link")}
            </Button>
            {#if sessionStore.getPlayerPermissions().isHost}
                <Button on:click={() => {}}>
                    {$_("lobby.regenerate_join_code")}
                </Button>
            {/if}
        </div>
    </Popover>
    {/if}

    <!-- Start game button -->
    {#if sessionStore.getPlayerPermissions().isHost && sessionStore.getState().gameState === GameState.Lobby}
        <Button
            class="w-xs mx-auto text-dark bg-green-200 dark:bg-green-900 hover:bg-green-200 dark:hover:bg-green-900 backdrop-blur-lg border border-black/20 dark:border-white/20 shadow-lg p-4 rounded-2xl flex items-center justify-between transition-all cursor-pointer"
            on:click={() => {
                sessionStore.startGame();
            }}
        >
            <div class="grid justify-items-start">
                <div class="relative">
                    <span
                        class="text-xl font-semibold tracking-widest select-none transition-opacity duration-300"
                        >{$_("lobby.start_game")}
                    </span>
                </div>
            </div>
            <Play />
        </Button>
    {/if}
</div>

{#if players.length > 5}
    <!-- Search Bar -->
    <TableSearch
        bind:inputValue={searchQuery}
        placeholder={$_("lobby.search_player")}
    />
{/if}

<!-- Players Table -->
<Table striped hoverable noborder class="mb-16">
    <TableHead>
        <TableHeadCell class="cursor-pointer flex items-center">
            {$_("lobby.player_name")}
        </TableHeadCell>
        <TableHeadCell>{$_("lobby.status")}</TableHeadCell>
    </TableHead>

    <TableBody tableBodyClass="divide-y">
        {#each filteredPlayers() as player}
            <TableBodyRow class="!bg-black/2 hover:!bg-black/4 dark:!bg-white/20 dark:hover:!bg-white/30">
                <TableBodyCell>
                    {player.Username}
                    {#if sessionStore.isCurrentPlayer(player.PlayerId)}
                        <Badge color="purple" class="ml-1"
                            >{$_("lobby.you")}</Badge
                        >
                    {/if}
                    {#if sessionStore.getPlayerPermissions(player.PlayerId).isHost}
                        <Badge color="blue" class="ml-1"
                            >{$_("lobby.host")}</Badge
                        >
                    {/if}
                </TableBodyCell>
                <TableBodyCell>
                    {#if player.IsConnected}
                        <Badge color="green"
                            >{$_(`player_status.connected`)}</Badge
                        >
                    {:else}
                        <Badge color="yellow"
                            >{$_(`player_status.disconnected`)}</Badge
                        >
                    {/if}
                </TableBodyCell>
                <!-- Can kick and rename player -->
                {#if sessionStore.getPlayerPermissions().isHost}
                    <TableBodyCell>
                        <!-- kick player -->
                        <Button
                            outline={true}
                            color="alternative"
                            class="p-2! text-red-800 hover:bg-red-500"
                            size="lg"
                            on:click={() => {
                                showKickModal = true;
                                kick_player = player.PlayerId;
                            }}
                        >
                            <UserX class="w-7 h-7" />
                        </Button>
                        <Tooltip type="auto">{$_("lobby.kick_player")}</Tooltip>
                        <!-- rename player -->
                        <Button
                            outline={true}
                            color="alternative"
                            class="p-2! text-blue-800 hover:bg-blue-500"
                            size="lg"
                            on:click={() => {
                                showRenameModal = true;
                                rename_player = player.PlayerId;
                            }}
                        >
                            <TextCursorInput class="w-7 h-7" />
                        </Button>
                        <Tooltip type="auto"
                            >{$_("lobby.rename_player")}</Tooltip
                        >
                    </TableBodyCell>
                {/if}
            </TableBodyRow>
        {/each}
    </TableBody>
</Table>
