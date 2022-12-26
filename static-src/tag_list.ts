import TagList from './TagList.svelte'

const params = document.currentScript.getAttribute('data-params')

const app = new TagList({
    target: document.getElementById('app'),
    props: {
        params: JSON.parse(params)
    },
})

export default app;