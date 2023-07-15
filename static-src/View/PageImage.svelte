<script lang="ts">
    export let alt: string;
    export let src: string;

    let loading: "lazy" | "eager" = "lazy";

    export function forceLoad() {
        loading = "eager";
    }
</script>

<img
    class="ms-auto me-auto"
    {loading}
    {alt}
    {src}
    style="object-fit:contain;max-width:100%;max-height:100%"
/>
