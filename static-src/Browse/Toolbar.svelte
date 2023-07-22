<script lang="ts">
    export let Title = "";
    export let BrowseURL = "";
    export let TagListURL = "";
    export let SortBy = "";
    export let SortOrder = "";
    export let FavoriteOnly = false;
    export let Tag = "";
    export let TagFavorite = false;

    export let changeSort;
    export let changeOrder;
    export let onFilterFavorite;
    export let rescanLibrary;
    export let onTagFavorite;
    export let onSearchClick;
    export let onAboutClick;

    export let SearchText = "";
</script>

<nav class="navbar navbar-dark bg-dark fixed-top navbar-expand-lg">
    <div class="container-fluid">
        <span class="navbar-brand text-truncate">{Title}</span>

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
                        aria-haspopup="true"
                        aria-expanded="false"
                    >
                        Browse
                    </button>
                    <div
                        class="dropdown-menu"
                        aria-labelledby="navbarBrowseDropdown"
                    >
                        <a class="dropdown-item" type="button" href={BrowseURL}>
                            <i class="bi bi-list-ul me-3" /> All items
                        </a>

                        <a
                            class="dropdown-item"
                            type="button"
                            href={TagListURL}
                        >
                            <i class="bi bi-tags-fill me-3" /> Tag list
                        </a>
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
                        Sort by
                    </button>
                    <div class="dropdown-menu" aria-labelledby="navbarDropdown">
                        <button
                            class="dropdown-item"
                            class:active={SortBy === "name"}
                            type="button"
                            on:click={(e) => changeSort("name")}
                        >
                            <i class="bi bi-type me-3" /> Name
                        </button>

                        <button
                            class="dropdown-item"
                            class:active={SortBy === "createTime"}
                            type="button"
                            on:click={(e) => changeSort("createTime")}
                        >
                            <i class="bi bi-clock me-3" /> Added date
                        </button>

                        <div class="dropdown-divider" />

                        <button
                            class="dropdown-item"
                            class:active={SortOrder === "ascending"}
                            type="button"
                            on:click={(e) => changeOrder("ascending")}
                        >
                            <i class="bi bi-sort-down-alt me-3" />
                            Ascending
                        </button>

                        <button
                            class="dropdown-item"
                            class:active={SortOrder === "descending"}
                            type="button"
                            on:click={(e) => changeOrder("descending")}
                        >
                            <i class="bi bi-sort-down me-3" />
                            Descending
                        </button>
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
                        Filter
                    </button>
                    <div class="dropdown-menu" aria-labelledby="navbarDropdown">
                        <button
                            class="dropdown-item"
                            type="button"
                            class:active={FavoriteOnly}
                            id="filter-favorite"
                            on:click={onFilterFavorite}
                        >
                            <i class="bi bi-star-fill me-3" /> Favorite
                        </button>
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
                            type="button"
                            on:click={rescanLibrary}
                        >
                            <i class="bi bi-arrow-clockwise me-3" /> Re-scan library
                        </button>
                    </div>
                </li>
                <li class="nav-item">
                    <button class="nav-link" on:click={onAboutClick}
                        >About</button
                    >
                </li>
            </ul>
            <ul
                class="navbar-nav ms-lg-2 mb-2 mb-lg-0"
                class:d-none={Tag === ""}
            >
                <li class="nav-item">
                    <button
                        id="favorite-btn"
                        class="btn"
                        class:btn-pink={TagFavorite}
                        class:active={TagFavorite}
                        class:btn-outline-pink={!TagFavorite}
                        on:click={onTagFavorite}
                    >
                        <i class="bi bi-star-fill me-3" /> Favorite tag
                    </button>
                </li>
            </ul>
            <div class="d-flex ms-lg-2">
                <div class="input-group">
                    <input
                        class="form-control"
                        type="search"
                        placeholder="Search"
                        aria-label="Search"
                        id="search-text"
                        bind:value={SearchText}
                        on:keydown={(e) => {
                            if (e.key === "Enter") onSearchClick(SearchText);
                        }}
                    />
                    <button
                        class="btn btn-outline-success"
                        type="button"
                        id="search-button"
                        on:click={(e) => onSearchClick(SearchText)}
                    >
                        <i class="bi bi-search" />
                    </button>
                </div>
            </div>
        </div>
    </div>
</nav>
