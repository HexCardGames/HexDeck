<script lang="ts">
    import { Button, Card, Skeleton } from "flowbite-svelte";
    import { MoveLeft, RefreshCcw } from "lucide-svelte";
    import { _ } from "svelte-i18n";
    import Markdown from "svelte-exmarkdown";
    import { onMount } from "svelte";
    import { goto } from "@roxi/routify";

    // Reactive stores for better state management
    let md: string = "";
    let loading: boolean = true;
    let error: string | null = null;

    function goBack() {
        window.history.back();
    }

    async function getImprintMd() {
        loading = true;
        error = null;

        try {
            const controller = new AbortController();
            const timeout = setTimeout(() => controller.abort(), 5000);

            const response = await fetch(`/api/imprint`, {
                method: "GET",
                headers: {
                    "Content-Type": "application/json",
                },
                signal: controller.signal,
            });

            clearTimeout(timeout);

            if (!response.ok) throw new Error("Server error");

            const data = await response.json();
            if (!data?.Content) throw new Error("Empty response");

            md = data.Content;
        } catch (err: any) {
            if (err.name === "AbortError") {
                error = "imprint.timeout_while_loading";
            } else {
                error = "imprint.something_went_wrong";
            }
        } finally {
            loading = false;
        }
    }

    onMount(getImprintMd);
</script>

<div class="container mx-auto p-6">
    <Button
        color="none"
        class="border-2 border-gray-500 dark:border-gray-300 hover:bg-gray-500 dark:hover:bg-gray-300 hover:text-white dark:hover:text-black rounded-full text-gray-500 dark:text-gray-300 mb-12"
        on:click={goBack}
    >
        <MoveLeft class="mr-2" />
        <span>{$_("imprint.go_back")}</span>
    </Button>

    <Card class="max-w-lg mx-auto dark:text-gray-200 rounded-xl">
        <h1 class="text-2xl font-bold mb-6">{$_("imprint.title")}</h1>

        {#if loading}
            <div class="w-full">
                <Skeleton size="lg" />
            </div>
        {:else if error}
            <div class="text-red-400 text-lg font-semibold grid">
                {$_(error)}

                <Button
                    color="none"
                    class="border-2 border-gray-500 dark:border-gray-300 hover:bg-gray-500 dark:hover:bg-gray-300 hover:text-white dark:hover:text-black rounded-full text-gray-500 dark:text-gray-300 mt-4"
                    on:click={getImprintMd}
                >
                    <RefreshCcw class="mr-2" />
                    <span>{$_("imprint.retry")}</span>
                </Button>
            </div>
        {:else}
            <Markdown {md} />
        {/if}
    </Card>
</div>
