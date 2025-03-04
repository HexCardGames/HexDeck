<script lang="ts">
    import { Button, Card } from "flowbite-svelte";
    import { MoveLeft } from "lucide-svelte";
    import { _ } from "svelte-i18n";
    import { url } from '@roxi/routify'

    function goBack() {
        window.history.pushState({}, "", "/");
    }

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

     <Card class="max-w-lg mx-auto dark:text-gray-200 rounded-xl text-center space-y-2">
        <h1 class="text-2xl font-bold">{$_("404.404_page_not_found")}</h1>
        <h1 class="">{$_("404.page_not_found", { values: { page: $url('$leaf')  }})}</h1>
    </Card>
</div>
