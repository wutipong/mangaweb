<script>
    import Toolbar from "./View/Toolbar.svelte";
    import ImageViewer from "./View/ImageViewer.svelte";
    import Toast from "./Common/Toast.svelte";
    import ModalDialog from "./Common/ModalDialog.svelte";

    export let params

    let favorite = params.Favorite
    let name = params.Name
    let tags = params.Tags
    let browseURL = params.BrowseURL

    let current = 0
    let toast
    let aboutDialog

    function downloadManga() {
        download(params.DownloadURL)
    }

    function downloadPage() {
        download(params.DownloadPageURLs[current])
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

    async function updateCover() {
        const url = new URL(params.UpdateCoverURLs[current], window.location.origin)

        await fetch(url)
        toast.show('Update Cover', 'The cover image is updated successfully.')
    }

    function onIndexChange(i) {
        current = i
    }

    function download(url) {
        let link = document.createElement("a");
        link.setAttribute('download', '');
        link.href = url;
        document.body.appendChild(link);

        link.click();
        link.remove();
    }

    function onAboutClick() {
        aboutDialog.show()
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
         updateCover={updateCover}
         onAboutClick={onAboutClick}
/>

<Toast bind:this={toast}/>

<ModalDialog Id="aboutModal" Title="About" bind:this={aboutDialog}>
    <h5>MangaWeb</h5>
    <h6>Version {params.Version} </h6>
    <p>&copy; 2021-2022 Wutipong Wongsakuldej. All Right Reserved</p>
    <p>Licensed under MIT License</p>
    <p><a href='https://github.com/wutipong/mangaweb'>Homepage</a></p>
</ModalDialog>