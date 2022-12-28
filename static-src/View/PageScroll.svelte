<script lang="ts">
    export let PageCount = 0
    export let onValueChange
    export let Current = 0
</script>

<nav aria-label='Move to top navigation' class='position-fixed bottom-0 end-0 p-3' style="z-index: 12;">
    <button class="btn btn-primary" type="button" data-bs-toggle="offcanvas" data-bs-target="#offcanvasBottom"
            aria-controls="offcanvasBottom">
        <i class='bi bi-chevron-double-up'></i>
        <span class='d-none d-sm-block'>Move to page</span>
    </button>
</nav>

<div class="offcanvas offcanvas-bottom" tabindex="-1" id="offcanvasBottom" aria-labelledby="offcanvasBottomLabel">
    <div class="offcanvas-header">
        Move to page
    </div>
    <div class="offcanvas-body small">
        <input type="range" class="form-range" min="0" max="{PageCount-1}" bind:value={Current}
               on:change={e => onValueChange(Current)}/>
        {Current + 1} of {PageCount}
    </div>
</div>