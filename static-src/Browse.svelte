<script>
    import Toolbar from "./Browse/Toolbar.svelte";
    import AboutDialog from "./Common/AboutDialog.svelte";
    import Item from "./Browse/Item.svelte";
    import Pagination from "./Common/Pagination.svelte";
    import PageItem from "./Common/PageItem.svelte";
    import Notification from "./Common/Notification.svelte";

    export let params;

    let toast;
    let tagFavorite = params.TagFavorite;
    let aboutDialog;

    function changeSort(sortBy) {
        let url = window.location;
        let searchParams = new URLSearchParams(url.search);
        searchParams.set("sort", sortBy);

        if (sortBy === "name") {
            searchParams.set("order", "ascending");
        } else if (sortBy === "createTime") {
            searchParams.set("order", "descending");
        }

        searchParams.delete("page");

        url.search = searchParams.toString();
    }

    function changeOrder(order) {
        let url = window.location;
        let searchParams = new URLSearchParams(url.search);

        searchParams.set("order", order);

        url.search = searchParams.toString();
    }

    function onFilterFavorite() {
        let url = window.location;
        let searchParams = new URLSearchParams(url.search);

        let isFavorite = params.FavoriteOnly;
        searchParams.set("favorite", (!isFavorite).toString());

        url.search = searchParams.toString();
    }

    async function rescanLibrary() {
        const url = params.RescanURL;
        await fetch(url);
        toast.show(
            "Re-scan Library",
            "Library re-scanning in progress. Please refresh after a few minutes."
        );
    }

    async function onTagFavorite() {
        tagFavorite = !tagFavorite;

        const urlSearchParams = new URLSearchParams();
        urlSearchParams.set("favorite", tagFavorite.toString());

        const url = new URL(params.SetTagFavoriteURL, window.location.origin);
        url.search = urlSearchParams.toString();

        const resp = await fetch(url);
        const json = await resp.json();

        if (json.favorite) {
            toast.show(
                "Favorite",
                `The tag "${params.Tag}" is now your favorite.`
            );
        } else {
            toast.show(
                "Favorite",
                `The tag "${params.Tag}" is no longer your favorite.`
            );
        }
    }

    function onSearchClick(t) {
        let searchText = t;
        let url = window.location;
        let searchParams = new URLSearchParams(url.search);
        searchParams.set("search", searchText);

        url.search = searchParams.toString();
    }

    function onAboutClick() {
        aboutDialog.show();
    }
</script>

<Toolbar
    Title={params.Title}
    BrowseURL={params.BrowseURL}
    TagListURL={params.TagListURL}
    SortBy={params.SortBy}
    SortOrder={params.SortOrder}
    FavoriteOnly={params.FavoriteOnly}
    Tag={params.Tag}
    TagFavorite={tagFavorite}
    {changeSort}
    {changeOrder}
    {onFilterFavorite}
    {rescanLibrary}
    {onTagFavorite}
    {onSearchClick}
    SearchText={params.SearchText}
    {onAboutClick}
/>

<div class="container-fluid" style="padding-top:100px;">
    <div class="grid-container">
        {#each params.Items as item}
            <Item
                Favorite={item.Favorite}
                IsRead={item.IsRead}
                ID={item.ID}
                ViewURL={item.ViewURL}
                ThumbnailURL={item.ThumbnailURL}
                Name={item.Name}
            />
        {/each}
    </div>
</div>
<div style="height: 100px;" />

<Pagination>
    {#each params.Pages as page}
        <PageItem
            IsActive={page.IsActive}
            IsEnabled={page.IsEnabled}
            IsHiddenOnSmall={page.IsHiddenOnSmall}
            URL={page.LinkURL}
            Content={page.Content}
        />
    {/each}
</Pagination>

<AboutDialog bind:this={aboutDialog} version={params.Version} />

<Notification bind:this={toast} />

<nav
    aria-label="Move to top navigation"
    class="position-fixed bottom-0 end-0 p-3"
>
    <a class="btn btn-secondary" href="#top">
        <i class="bi bi-chevron-double-up" />
        <span class="d-none d-sm-block">Top</span>
    </a>
</nav>
