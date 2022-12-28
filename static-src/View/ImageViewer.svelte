<script lang="ts">
    import {onMount} from "svelte";
    import * as bootstrap from "bootstrap"

    export let ImageURLs = []
    export let onIndexChange
    let carousel

    onMount(async () => {
        let carouselControl = document.querySelector('#carouselControl')
        carousel = new bootstrap.Carousel(carouselControl, {
            interval: false
        })

        carouselControl.addEventListener('slide.bs.carousel', e => {
            if (onIndexChange) onIndexChange(e.to)
        })

    })

    export function advance(n) {
        carousel.to(n)
    }
</script>

<div class='carousel slide w-100 h-100' id='carouselControl'>
    <div class='carousel-inner w-100 h-100' id='carousel' style='width:100%; height:100%;'>
        {#each ImageURLs as url, index}
            <div class='carousel-item w-100 h-100'
                 class:active={index === 0}>
                <div class='w-100 h-100 d-flex flex-col'>
                    <img class='ms-auto me-auto' loading='lazy' alt='page {index}'
                         src='{url}'
                         style='object-fit:contain;max-width:100%;max-height:100%'>
                </div>
            </div>
        {/each}
    </div>
    <button class='carousel-control-prev' data-bs-target='#carouselControl' data-bs-slide='prev'>
        <span class='carousel-control-prev-icon' aria-hidden='true'></span>
        <span class='visually-hidden'>Previous</span>
    </button>
    <button class='carousel-control-next' data-bs-target='#carouselControl' data-bs-slide='next'>
        <span class='carousel-control-next-icon' aria-hidden='true'></span>
        <span class='visually-hidden'>Next</span>
    </button>
</div>