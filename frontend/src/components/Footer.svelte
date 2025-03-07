<script>
    import { Banner, Button, Tooltip } from "flowbite-svelte";
    import { ArrowUpFromLine, Dot } from "lucide-svelte";
    import { slide } from "svelte/transition";
    import { quintOut } from "svelte/easing";
    import { _ } from "svelte-i18n";

    let bannerStatus = true;
    const params = { delay: 250, duration: 500, easing: quintOut };

    function showFooter() {
        bannerStatus = true;
    }
</script>

<Banner id="bottom-banner" position="absolute" bannerType="bottom" transition={slide} {params} bind:bannerStatus>
    <p class="flex items-center text-sm font-normal text-gray-500 dark:text-gray-400">
        <a href="/imprint">{$_("footer.imprint")}</a>
        <Dot />
        <a href="https://github.com/HexCardGames/HexDeck" target="_blank">{$_("footer.github")}</a>
    </p>
</Banner>

<div class="absolute right-0 bottom-0">
    <Button on:click={showFooter} class="!p-2 focus:ring-0 mt-2" color="none">
        <ArrowUpFromLine />
    </Button>
    <Tooltip type="auto">{$_("landing_page.show_footer")}</Tooltip>
</div>
