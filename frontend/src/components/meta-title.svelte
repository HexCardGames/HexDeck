<script lang="ts">
    import { _, waitLocale } from "svelte-i18n";
    import { onMount } from "svelte";

    let title = "Loading...";

    onMount(async () => {
        // Ensure translations are loaded
        await waitLocale();
        title = $_("page_name");
    });
</script>

<svelte:head>
    <title>{title}</title>
</svelte:head>
