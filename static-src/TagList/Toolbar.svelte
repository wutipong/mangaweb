<script lang="ts">
    export let title = "";
    export let browseURL = "";
    export let tagListURL = "";
    export let onFilterFavorite;
    export let favoriteOnly = false;
    export let onAboutClick;
</script>

<nav class="navbar navbar-dark bg-dark fixed-top navbar-expand-lg">
    <div class="container-fluid">
        <span class="navbar-brand text-truncate">{title}</span>

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
            <ul class="navbar-nav me-auto mb-2 mb-lg-0">
                <li class="nav-item dropdown">
                    <button
                        class="nav-link dropdown-toggle"
                        id="navbarBrowseDropdown"
                        data-bs-toggle="dropdown"
                        aria-expanded="false"
                    >
                        Browse
                    </button>
                    <ul class="dropdown-menu">
                        <li>
                            <a class="dropdown-item" href={browseURL}>
                                <i class="bi bi-list-ul me-3" /> All items
                            </a>
                        </li>

                        <li>
                            <a class="dropdown-item" href={tagListURL}>
                                <i class="bi bi-tags-fill me-3" /> Tag list
                            </a>
                        </li>
                    </ul>
                </li>
                <li class="nav-item dropdown">
                    <button
                        class="nav-link dropdown-toggle"
                        id="navbarDropdown"
                        data-bs-toggle="dropdown"
                        aria-expanded="false"
                    >
                        Filter
                    </button>
                    <div class="dropdown-menu">
                        <button
                            class="dropdown-item"
                            type="button"
                            id="filter-favorite"
                            on:click={onFilterFavorite}
                            class:active={favoriteOnly}
                        >
                            <i class="bi bi-star-fill me-3" /> Favorite
                        </button>
                    </div>
                </li>
                <li class="nav-item">
                    <button class="nav-link" on:click={onAboutClick}
                        >About</button
                    >
                </li>
            </ul>
        </div>
    </div>
</nav>
