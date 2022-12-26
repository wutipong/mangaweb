<script>
    import Toolbar from "./TagList/Toolbar.svelte";
    import Item from "./TagList/Item.svelte";

    export let params
    export let browseURL = ""
    export let tagListURL = ""

    let favoriteOnly = false

    function toggleFavoriteFilter() {
        favoriteOnly = !favoriteOnly
        if (favoriteOnly) {
            $('.not-favorite').addClass('d-none')
            $('#filter-favorite').addClass('active')
        } else {
            $('.not-favorite').removeClass('d-none')
            $('#filter-favorite').removeClass('active')
        }
    }
</script>

<Toolbar title="{params.Title}" browseURL="{browseURL}" tagListURL="{tagListURL}" onToggleFavoriteFilter=""></Toolbar>

<div class='container-fluid' style='padding-top:100px;'>
    <div class='grid-container'>
        {#each params.Tags as tag}
            <Item name={tag.Name} favorite={tag.Favorite} id={tag.ID} url={tag.URL} thumbnailURL={tag.ThumbnailURL}></Item>
        {/each}
    </div>
</div>