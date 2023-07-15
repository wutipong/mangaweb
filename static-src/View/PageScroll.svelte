<script lang="ts">
    export let PageCount = 0
    export let onValueChange
    export let Current = 0
</script>

<nav aria-label='Move to top navigation' class='position-fixed bottom-0 vw-100' style="z-index: 12;">
    <button class="btn btn-transparent px-5 vw-100" type="button" data-bs-toggle="offcanvas" data-bs-target="#offcanvasBottom"
            aria-controls="offcanvasBottom">
        <i class='bi bi-chevron-double-up'></i>
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

<style>
    button {
    --bs-btn-color: #808080;
    --bs-btn-bg: transparent;
    --bs-btn-border-color: transparent;
    --bs-btn-hover-color: #fff;
    --bs-btn-hover-bg: #2f4d6c;
    --bs-btn-hover-border-color: #2c4866;
    --bs-btn-focus-shadow-rgb: 85, 115, 146;
    --bs-btn-active-color: #fff;
    --bs-btn-active-bg: #2c4866;
    --bs-btn-active-border-color: #29445f;
    --bs-btn-active-shadow: inset 0 3px 5px #00000020;
    --bs-btn-disabled-color: #fff;
    --bs-btn-disabled-bg: #375a7f;
    --bs-btn-disabled-border-color: #375a7f;
}
</style>