<script>
    import Toolbar from "./Browse/Toolbar.svelte";
    import {onMount} from "svelte";

    export let params

    let searchText = ""
    console.log(params)

    function changeSort(s){
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

    onMount(async() =>{
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