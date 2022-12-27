<script>
    import Toolbar from "./View/Toolbar.svelte";
    import ImageViewer from "./View/ImageViewer.svelte";

    export let params
    console.log(params)

    let favorite = params.Favorite
    let name = params.Name
    let tags = params.TagsV2
    let browseURL = params.BrowseURL

    let current = 0

    function downloadManga() {
        browser.downloads.download({
            url: params.DownloadURL
        })
    }

    function downloadPage() {
        browser.downloads.download({
            url: params.DownloadPageURLs[current]
        })
    }

    function toggleFavorite() {
        favorite = !favorite
    }

    function updateCover() {

    }

    function onIndexChange(i) {
        current = i
    }

</script>

<div class='fullscreen' style='padding-top:80px;'>
    <ImageViewer ImageURLs={params.ImageURLs}
                 onIndexChange={onIndexChange}/>
</div>

<Toolbar Tags={tags}
         Name={name}
         Favorite={favorite}
         BrowseURL={browseURL}
         onDownloadManga={downloadManga}
         onDownloadPage={downloadPage}
         toggleFavorite={toggleFavorite}
         updateCover={updateCover}/>