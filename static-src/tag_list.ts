import TagList from './TagList.svelte'

const params = document.currentScript.getAttribute('data-params')
const browseURL = document.currentScript.getAttribute('data-browse-url')
const tagListURL = document.currentScript.getAttribute('data-tag-list-url')

const app = new TagList({
    target: document.getElementById('app'),
    props: {
        params: JSON.parse(params),
        browseURL: browseURL,
        tagListURL: tagListURL,
    },
})

export default app;