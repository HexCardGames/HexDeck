<script>
    import Footer from "../components/Footer.svelte";
    import ConnectRoom from "../components/ConnectRoom.svelte";
    import StatsContainer from "../components/StatsContainer.svelte";
    import { GradientButton } from "flowbite-svelte";
    import { _ } from "svelte-i18n";
    import options from "../stores/pageoptions";
</script>

<div class="overflow-auto size-full pb-18">
    <div class="flex justify-center mb-8">
        <div>
            <!-- Top Title container -->
            <div class="flex items-center space-x-2 my-6">
                <div>
                    <svelte:component this={options.page_icon} size="5.2rem" />
                </div>
                <div>
                    <h2 class="text-4xl">{$_("page_name")}</h2>
                    <h3 class="text-xl">
                        {$_("landing_page.sub_title")}
                    </h3>
                </div>
            </div>

            <!-- Join or create rooms -->
            <ConnectRoom />
        </div>
    </div>
    <div class="flex justify-center">
        <div class="grid grid-cols-1 sm:grid-cols-2 gap-8">
            <!-- Open source container -->
            <div class="p-4 bg-primary-50 dark:bg-primary-950 rounded-xl grid content-start justify-items-center w-3xs text-center space-y-2 border-1 border-primary-200 dark:border-primary-800">
                <svg xmlns="http://www.w3.org/2000/svg" width="48" height="48" viewBox="0 0 24 24"
                    ><path
                        fill="none"
                        stroke="currentColor"
                        stroke-linecap="round"
                        stroke-linejoin="round"
                        stroke-width="1.5"
                        d="M15.157 20.136c.211.51.8.757 1.284.492a9.25 9.25 0 1 0-8.882 0c.484.265 1.073.018 1.284-.492l1.358-3.28c.212-.51-.043-1.086-.478-1.426a3.7 3.7 0 1 1 4.554 0c-.435.34-.69.916-.478 1.426z"
                    /></svg
                >
                <h4 class="text-xl font-semibold">
                    {$_("landing_page.open_source_container.title")}
                </h4>
                <span>{$_("landing_page.open_source_container.content")}</span>
                <GradientButton color="purpleToBlue" class="opacity-75 focus:opacity-100 focus:ring-0" href="https://github.com/HexCardGames/HexDeck" target="_blank">
                    {$_("landing_page.open_source_container.github")}
                </GradientButton>
            </div>

            <!-- stats container -->
            <StatsContainer />
        </div>
    </div>
</div>
<Footer />
