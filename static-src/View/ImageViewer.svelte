<script lang="ts">
    import { onMount } from "svelte";
    import * as bootstrap from "bootstrap";
    import PageImage from "./PageImage.svelte";

    export let imageURLs = [];
    export let onIndexChange: Function;

    let carousel: bootstrap.Carousel;
    let imgs: PageImage[] = [];

    onMount(async () => {
        let carouselControl = document.querySelector("#carouselControl");
        carousel = new bootstrap.Carousel(carouselControl, {
            interval: false,
        });

        carouselControl.addEventListener("slide.bs.carousel", (e) => {
            if (onIndexChange) {
                onIndexChange(e.to);
            }

            const next = Math.min(e.to + 1, imgs.length - 1);
            imgs[next]?.forceLoad();

            const prev = Math.min(e.to - 1, imgs.length - 1);
            imgs[prev]?.forceLoad();
        });

        imgs[0]?.forceLoad();
        imgs[1]?.forceLoad();
        imgs[imgs.length -1]?.forceLoad();
    });

    export function advance(n) {
        carousel.to(n);
    }
</script>

<div class="carousel slide w-100 h-100" id="carouselControl">
    <div
        class="carousel-inner w-100 h-100"
        id="carousel"
        style="width:100%; height:100%;"
    >
        {#each imageURLs as url, index}
            <div class="carousel-item w-100 h-100" class:active={index === 0}>
                <div class="w-100 h-100 d-flex flex-col">
                    <PageImage
                        alt="page {index}"
                        src={url}
                        bind:this={imgs[index]}
                    />
                </div>
            </div>
        {/each}
    </div>
    <button
        class="carousel-control-prev"
        data-bs-target="#carouselControl"
        data-bs-slide="prev"
    >
        <span class="carousel-control-prev-icon" aria-hidden="true" />
        <span class="visually-hidden">Previous</span>
    </button>
    <button
        class="carousel-control-next"
        data-bs-target="#carouselControl"
        data-bs-slide="next"
    >
        <span class="carousel-control-next-icon" aria-hidden="true" />
        <span class="visually-hidden">Next</span>
    </button>
</div>
