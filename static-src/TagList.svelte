<script>
    import Toolbar from "./TagList/Toolbar.svelte";
    import Item from "./TagList/Item.svelte";
    import ModalDialog from "./Common/ModalDialog.svelte";

    export let params

    let favoriteOnly = false
    let aboutDialog

    function toggleFavoriteOnly() {
        favoriteOnly = !favoriteOnly
    }

    function onAboutClick() {
        aboutDialog.show()
    }

</script>

<Toolbar
        title={params.Title}
        browseURL={params.BrowseURL}
        tagListURL={params.TagListURL}
        onFilterFavorite={toggleFavoriteOnly}
        favoriteOnly={favoriteOnly}
        onAboutClick={onAboutClick}
/>

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

<ModalDialog Id="aboutModal" Title="About" bind:this={aboutDialog}>
    <h5>MangaWeb</h5>
    <h6>Version {params.Version} </h6>
    <p>&copy; 2021-2023 Wutipong Wongsakuldej. All Right Reserved</p>
    <p>Licensed under MIT License</p>
    <p><a href='https://github.com/wutipong/mangaweb'>Homepage</a></p>
</ModalDialog>