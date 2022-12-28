<script lang="ts">
    import {Toast} from "bootstrap";

    let message = ""
    let header = ""

    export function show(h: string, m: string) {
        message = m
        header = h

        const toast = new Toast(document.getElementById('toast'))
        toast.show()
    }
</script>


<div class='toast-container position-fixed bottom-0 end-0 p-3' style='z-index: 11'>
    <div id='toast' class='toast' role='alert' aria-live='assertive' aria-atomic='true'>
        <div class='toast-header'>
            <strong class='me-auto' id='toast-header'>{header}</strong>
            <button type='button' class='btn-close' data-bs-dismiss='toast' aria-label='Close'></button>
        </div>
        <div class='toast-body' id='toast-content'>
            {message}

        </div>
    </div>
</div>
