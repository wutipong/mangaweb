<script lang="ts">
    export let Favorite = false;
    export let IsRead = false;
    export let ID = "";
    export let ViewURL = "";
    export let ThumbnailURL = "";
    export let Name = "";
</script>

<div
    class="card mb-3"
    class:border={Favorite || !IsRead}
    class:border-pink={Favorite}
    class:border-yellow={!IsRead}
    id={ID}
>
    {#if Favorite}
        <span class="position-absolute bottom-0 end-0 p-2 badge bg-pink">
            <span><i class="bi bi-star-fill" /> Favorite </span>
        </span>
    {/if}

    {#if !IsRead}
        <span class="position-absolute bottom-0 end-0 p-2 badge bg-yellow">
            <span><i class="bi bi-lightning-fill" /> New </span>
        </span>
    {/if}

    <div class="row g-0">
        <div class="col-md-4">
            <a href={ViewURL} class="thumb-link">
                <img
                    class="img-fluid lazy rounded thumb-img"
                    loading="lazy"
                    src={ThumbnailURL}
                    alt={Name}
                />
            </a>
        </div>
        <div class="col-md-8">
            <div class="card-body">
                <h5 class="card-title">
                    <a href={ViewURL}>{Name}</a>
                </h5>
            </div>
        </div>
    </div>
</div>
