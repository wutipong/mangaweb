<script>
    import Toolbar from "./TagList/Toolbar.svelte";
    import Item from "./TagList/Item.svelte";
    import AboutDialog from "./Common/AboutDialog.svelte";

    export let params;

    let favoriteOnly = false;
    let aboutDialog;

    function toggleFavoriteOnly() {
        favoriteOnly = !favoriteOnly;
    }

    function onAboutClick() {
        aboutDialog.show();
    }
</script>

<Toolbar
    title={params.Title}
    browseURL={params.BrowseURL}
    tagListURL={params.TagListURL}
    onFilterFavorite={toggleFavoriteOnly}
    {favoriteOnly}
    {onAboutClick}
/>

<div class="container-fluid" style="padding-top:100px;">
    <div class="grid-container">
        {#each params.Tags as tag}
            {#if !favoriteOnly || (favoriteOnly && tag.Favorite)}
                <Item
                    name={tag.Name}
                    favorite={tag.Favorite}
                    id={tag.ID}
                    url={tag.URL}
                    thumbnailURL={tag.ThumbnailURL}
                />
            {/if}
        {/each}
    </div>
</div>

<AboutDialog bind:this={aboutDialog} version={params.Version} />

<nav
    aria-label="Move to top navigation"
    class="position-fixed bottom-0 end-0 p-3"
>
    <a class="btn btn-secondary" href="#top">
        <i class="bi bi-chevron-double-up" />
        <span class="d-none d-sm-block">Top</span>
    </a>
</nav>
