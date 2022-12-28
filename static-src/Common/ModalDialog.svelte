<script lang="ts">
    import {Modal} from 'bootstrap'

    export let Title=""
    export let Id =""

    function show(){
        const dialog = new Modal(`#${Id}`);
        dialog.show()
    }
</script>

<div class='modal fade' id={Id} tabindex='-1'>
    <div class='modal-dialog'>
        <div class='modal-content'>
            <div class='modal-header'>
                <h5 class='modal-title'>{Title}</h5>
                <button type='button' class='btn-close' data-bs-dismiss='modal' aria-label='Close'></button>
            </div>
            <div class='modal-body'>
                <slot></slot>
            </div>
            <div class='modal-footer'>
                <button type='button' class='btn btn-secondary' data-bs-dismiss='modal'>Close</button>
            </div>
        </div>
    </div>
</div>