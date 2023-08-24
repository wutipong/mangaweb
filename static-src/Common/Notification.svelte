<script lang="ts">
    import { Toast, ToastHeader, ToastBody } from "sveltestrap/src";

    let message = "";
    let header = "";

    let isOpen = false;
    export function show(h: string, m: string) {
        message = m;
        header = h;

        isOpen = true;
    }

    function toggle() {
        isOpen = !isOpen;
    }
</script>

<div
    class="toast-container position-fixed bottom-0 end-0 p-3"
    style="z-index: 11"
>
    <Toast autohide {isOpen} on:close={() => (isOpen = false)}>
        <ToastHeader {toggle}>{header}</ToastHeader>
        <ToastBody>{message}</ToastBody>
    </Toast>
</div>
