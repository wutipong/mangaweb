<script>
    import Toolbar from "./View/Toolbar.svelte";
    import ImageViewer from "./View/ImageViewer.svelte";
    import Toast from "./View/Toast.svelte";

    export let params
    console.log(params)

    let favorite = params.Favorite
    let name = params.Name
    let tags = params.TagsV2
    let browseURL = params.BrowseURL

    let current = 0
    let toast

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

    async function toggleFavorite() {
        favorite = !favorite

        const urlSearchParams = new URLSearchParams()
        urlSearchParams.set('favorite', favorite.toString())

        const url = new URL(params.SetFavoriteURL, window.location.origin)
        url.search = urlSearchParams.toString()

        await fetch(url)
        if (favorite) {
            toast.show('Favorite', 'The current manga is now your favorite.')
        } else {
            toast.show('Favorite', 'The current manga is no longer your favorite.')
        }
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

<Toast bind:this={toast}/>