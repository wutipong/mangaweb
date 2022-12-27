import View from './View.svelte'

const params = document.currentScript.getAttribute('data-params')

const app = new View({
    target: document.getElementById('app'),
    props: {
        params: JSON.parse(params)
    },
})

export default app;