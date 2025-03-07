<script lang="ts">
    import { Card } from "flowbite-svelte";
    import { _ } from "svelte-i18n";
    import { onMount } from "svelte";
    import { Spinner } from "flowbite-svelte";
    import { sessionStore } from "../stores/sessionStore";

    onMount(async () => {
        await sessionStore.leaveRoom();
    });
</script>

<div class="container mx-auto p-6">
    <Card class="max-w-lg mx-auto dark:text-gray-200 rounded-xl text-center space-y-2 flex">
        <Spinner class="text-primary-350 dark:text-primary-200 w-12" size="4.2" />
        <h1 class="text-2xl font-bold">{$_("leave.leaving_game")}</h1>
    </Card>
</div>
