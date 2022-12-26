import Browse from './Browse.svelte'

const params = document.currentScript.getAttribute('data-params')

const app = new Browse({
    target: document.getElementById('app'),
    props: {
        params: JSON.parse(params)
    },
})

export default app;