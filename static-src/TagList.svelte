<script>
    import Toolbar from "./TagList/Toolbar.svelte";
    import Item from "./TagList/Item.svelte";

    export let params
    export let browseURL = ""
    export let tagListURL = ""

    let favoriteOnly = false

    function toggleFavoriteOnly() {
        favoriteOnly = !favoriteOnly
    }
</script>

<Toolbar
        title="{params.Title}"
        browseURL="{browseURL}"
        tagListURL="{tagListURL}"
        onToggleFavoriteFilter="{toggleFavoriteOnly}"
        favoriteOnly={favoriteOnly}>
</Toolbar>

<div class='container-fluid' style='padding-top:100px;'>
    <div class='grid-container'>
        {#each params.Tags as tag}
            {#if !favoriteOnly || (favoriteOnly && tag.Favorite)}
                <Item
                        name={tag.Name}
                        favorite={tag.Favorite}
                        id={tag.ID}
                        url={tag.URL}
                        thumbnailURL={tag.ThumbnailURL}
                ></Item>
            {/if}
        {/each}
    </div>
</div>