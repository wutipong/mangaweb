<script lang="ts">
    export let name = ""
    export let favorite = false
    export let id = ""
    export let url = ""
    export let thumbnailURL = ""
</script>

<div class='card mb-3'
     class:border={favorite}
     class:border-pink={favorite}
     id='{id}'>

    {#if favorite}
    <span class='position-absolute bottom-0 end-0 p-2 badge bg-pink'>
        <span><i class='bi bi-star-fill'></i> Favorite</span>
    </span>
    {/if}

    <div class='row g-0'>
        <div class='col-md-4'>
            <a href='{url}' class='thumb-link'>
                <img class='img-fluid lazy rounded thumb-img' loading='lazy' src='{thumbnailURL}'
                     alt='{name}'>
            </a>
        </div>
        <div class='col-md-8'>
            <div class='card-body'>
                <h5 class='card-title'><a href='{url}'>{name}</a></h5>
            </div>
        </div>
    </div>
</div>