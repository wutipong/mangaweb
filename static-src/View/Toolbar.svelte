<script lang="ts">
    import { Button, Icon } from "sveltestrap/src";

    export let Tags = [];
    export let Name = "";
    export let Favorite = false;
    export let BrowseURL = "";

    export let onDownloadManga;
    export let onDownloadPage;
    export let toggleFavorite;
    export let updateCover;
    export let onAboutClick;
</script>

<nav class="navbar navbar-dark bg-dark fixed-top navbar-expand-lg">
    <div class="container-fluid">
        <span
            class="navbar-brand mr-auto text-truncate h1 d-none d-sm-block"
            style="max-width: 50%;">{Name}</span
        >
        <span
            class="navbar-brand mr-auto text-truncate h1 d-block d-sm-none"
            style="max-width: 50%;">View</span
        >

        <button
            class="navbar-toggler"
            type="button"
            data-bs-toggle="collapse"
            data-bs-target="#navbarSupportedContent"
            aria-controls="navbarSupportedContent"
            aria-expanded="false"
            aria-label="Toggle navigation"
        >
            <span class="navbar-toggler-icon" />
        </button>

        <div class="collapse navbar-collapse" id="navbarSupportedContent">
            <ul class="navbar-nav mb-2 me-auto">
                <li class="nav-item dropdown">
                    <button
                        class="nav-link dropdown-toggle"
                        id="navbarDropdown"
                        data-bs-toggle="dropdown"
                        aria-haspopup="true"
                        aria-expanded="false"
                    >
                        Tags
                    </button>
                    <div class="dropdown-menu" aria-labelledby="navbarDropdown">
                        {#each Tags as tag}
                            <a class="dropdown-item" href={tag.URL}>
                                <Icon name="tag" class="me-3" />{tag.Name}
                            </a>
                        {/each}
                    </div>
                </li>
                <li class="nav-item dropdown">
                    <button
                        class="nav-link dropdown-toggle"
                        id="navbarDropdown"
                        data-bs-toggle="dropdown"
                        aria-haspopup="true"
                        aria-expanded="false"
                    >
                        Tools
                    </button>
                    <div class="dropdown-menu" aria-labelledby="navbarDropdown">
                        <button
                            class="dropdown-item"
                            id="download-btn"
                            on:click={onDownloadPage}
                        >
                            <Icon name="download" class="me-3" />
                            Download Current Page
                        </button>
                        <button
                            class="dropdown-item"
                            on:click={onDownloadManga}
                        >
                            <Icon name="download" class="me-3" />
                            Download Manga
                        </button>
                        <div class="dropdown-divider" />
                        <button
                            class="dropdown-item"
                            id="update-cover-btn"
                            on:click={updateCover}
                        >
                            <Icon name="journal-arrow-up" class="me-3" />
                            Update Cover
                        </button>
                    </div>
                </li>
                <li class="nav-item">
                    <button class="nav-link" on:click={onAboutClick}
                        >About</button
                    >
                </li>
            </ul>
            <ul class="navbar-nav">
                <li class="nav-item mb-2 mb-lg-0 ms-lg-2">
                    <Button
                        id="favorite-btn"
                        color="pink"
                        active={Favorite}
                        outline={!Favorite}
                        on:click={toggleFavorite}
                    >
                        <Icon name="star-fill" class="me-3" /> Favorite
                    </Button>
                </li>
                <li class="nav-item mb-2 mb-lg-0 ms-lg-2">
                    <Button id="close-btn" color="danger" href={BrowseURL}>
                        <Icon name="x-circle-fill" class="me-3" /> Close
                    </Button>
                </li>
            </ul>
        </div>
    </div>
</nav>
