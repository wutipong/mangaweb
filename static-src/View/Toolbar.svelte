<script lang="ts">
    export let Tags = []
    export let Name = ""
    export let Favorite = false
    export let BrowseURL = ""

    export let onDownloadManga
    export let onDownloadPage
    export let toggleFavorite
    export let updateCover
</script>

<nav class='navbar navbar-dark bg-dark fixed-top navbar-expand-lg'>
    <div class='container-fluid'>
        <span class='navbar-brand mr-auto text-truncate h1 d-none d-sm-block' style='max-width: 50%;'>{Name}</span>
        <span class='navbar-brand mr-auto text-truncate h1 d-block d-sm-none' style='max-width: 50%;'>View</span>

        <button class='navbar-toggler' type='button' data-bs-toggle='collapse'
                data-bs-target='#navbarSupportedContent' aria-controls='navbarSupportedContent'
                aria-expanded='false' aria-label='Toggle navigation'>

            <span class='navbar-toggler-icon'></span>
        </button>

        <div class='collapse navbar-collapse' id='navbarSupportedContent'>
            <ul class='navbar-nav mb-2 me-auto'>
                <li class='nav-item dropdown'>
                    <a class='nav-link dropdown-toggle' href='#' id='navbarDropdown'
                       role='button' data-bs-toggle='dropdown' aria-haspopup='true' aria-expanded='false'>
                        <i class="bi bi-tags-fill"></i> Tags
                    </a>
                    <div class='dropdown-menu' aria-labelledby='navbarDropdown'>
                        {#each Tags as tag}
                            <a class='dropdown-item' href='{tag.URL}'>
                                <i class="bi bi-tag"></i>{tag.Name}
                            </a>
                        {/each}
                    </div>
                </li>
                <li class='nav-item dropdown'>
                    <a class='nav-link dropdown-toggle' href='#' id='navbarDropdown'
                            data-bs-toggle='dropdown' aria-haspopup='true' aria-expanded='false'>
                        <i class='bi bi-list'></i> Menu
                    </a>
                    <div class='dropdown-menu' aria-labelledby='navbarDropdown'>
                        <button class='dropdown-item' id='download-btn' download on:click={onDownloadPage}>
                            <i class='bi bi-download'></i>
                            Download Current Page
                        </button>
                        <button class='dropdown-item' download on:click={onDownloadManga}>
                            <i class='bi bi-download'></i>
                            Download Manga
                        </button>
                        <div class='dropdown-divider'></div>
                        <button class='dropdown-item' id='update-cover-btn' on:click={updateCover}>
                            <i class='bi bi-journal-arrow-up'></i>
                            Update Cover
                        </button>
                    </div>
                </li>
            </ul>
            <ul class='navbar-nav'>
                <li class='nav-item mb-2 mb-lg-0 ms-lg-2'>
                    <button id='favorite-btn'
                            class='btn'
                            class:btn-pink={Favorite}
                            class:active={Favorite}
                            class:btn-outline-pink={!Favorite}
                            on:click={toggleFavorite}
                    >
                        <i class='bi bi-star-fill'></i> Favorite
                    </button>
                </li>
                <li class='nav-item mb-2 mb-lg-0 ms-lg-2'>
                    <a id='close-btn' class='btn btn-danger' href='{BrowseURL}'>
                        <i class='bi bi-x-circle-fill'></i> Close
                    </a>
                </li>
            </ul>
        </div>
    </div>
</nav>