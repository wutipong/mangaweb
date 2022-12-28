<script lang="ts">
    export let IsActive = false
    export let IsEnabled = false
    export let IsHiddenOnSmall = false
    export let URL = ""
    export let Content = ""
</script>

<li class='page-item'
    class:active={IsActive}
    class:disabled={!IsEnabled}
    class:d-none={IsHiddenOnSmall}
    class:d-sm-block={IsHiddenOnSmall}
>
    <a class='page-link' href='{URL}'>
        {Content}
    </a>
</li>