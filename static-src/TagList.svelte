<script>
    import Toolbar from "./TagList/Toolbar.svelte";

    export let params
    export let browseURL = ""
    export let tagListURL = ""

</script>

<Toolbar title="{params.Title}" browseURL="{browseURL}" tagListURL="{tagListURL}" onToggleFavoriteFilter=""></Toolbar>