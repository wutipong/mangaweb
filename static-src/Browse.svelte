<script>
    import Toolbar from "./Browse/Toolbar.svelte";
    import {onMount} from "svelte";
    import ModalDialog from "./Common/ModalDialog.svelte";
    import Item from "./Browse/Item.svelte";

    export let params

    let searchText = ""
    console.log(params)

    function changeSort(s) {
        console.log(s)
    }

    function changeOrder(o) {
        console.log(o)
    }

    function toggleFavoriteFilter() {

    }

    function rescanLibrary() {

    }

    function toggleFavorite() {

    }

    function onSearchClick(t) {

    }

    onMount(async () => {
        let url = window.location
        let searchParams = new URLSearchParams(url.search)

        searchText = searchParams.get('search')
    })
</script>

<Toolbar Title={params.Title}
         BrowseURL={params.BrowseURL}
         TagListURL={params.TagListURL}
         SortBy={params.SortBy}
         SortOrder={params.SortOrder}
         FavoriteOnly={params.FavoriteOnly}
         Tag={params.Tag}
         TagFavorite={params.TagFavorite}
         changeSort={changeSort}
         changeOrder={changeOrder}
         toggleFavoriteFilter="{toggleFavoriteFilter}"
         rescanLibrary="{rescanLibrary}"
         toggleFavorite="{toggleFavorite}"
         onSearchClick="{onSearchClick}"
         SearchText={searchText}/>

<div class='container-fluid' style='padding-top:100px;'>
    <div class='grid-container'>
        {#each params.Items as item}
            <Item Favorite={item.Favorite}
                  IsRead={item.IsRead}
                  ID={item.ID}
                  ViewURL={item.ViewURL}
                  ThumbnailURL={item.ThumbnailURL}
                  Name={item.Name}/>
        {/each}
    </div>
</div>

<ModalDialog Id="aboutModal" Tittle="About">
    <h5>MangaWeb</h5>
    <h6>Version {params.Version} </h6>
    <p>&copy; 2021-2022 Wutipong Wongsakuldej. All Right Reserved</p>
    <p>Licensed under MIT License</p>
    <p><a href='https://github.com/wutipong/mangaweb'>Homepage</a></p>
</ModalDialog>