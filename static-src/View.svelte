<script>
    import Toolbar from "./View/Toolbar.svelte";

    export let params
    console.log(params)

    let favorite = params.Favorite
    let name = params.Name
    let tags = params.Tags
    let browseURL = params.BrowseURL

    function downloadManga(){

    }

    function downloadPage() {

    }

    function toggleFavorite() {
        favorite = !favorite
    }

    function updateCover() {

    }

</script>

<Toolbar Tags="{tags}"
         Name={name}
         Favorite={favorite}
         BrowseURL={browseURL}
         onDownloadManga="{downloadManga}"
         onDownloadPage="{downloadPage}"
         toggleFavorite="{toggleFavorite}"
         updateCover="{updateCover}">

</Toolbar>