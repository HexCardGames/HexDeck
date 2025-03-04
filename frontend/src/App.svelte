<script context="module">
  import { Router, createRouter } from "@roxi/routify";
  import routes from "../.routify/routes.default.js";
  import MetaTitle from "./components/meta-title.svelte";

  export const router = createRouter({ routes });

  router.ready().then(() => {
    console.log("Routify router is ready!");
  });
</script>

<MetaTitle />
<Router {router} />
