if (window.matchMedia('(prefers-color-scheme: dark)').matches) {
    document.head.insertAdjacentHTML(
        'beforeend',
        '<meta name="color-scheme" content="dark">'
    )
    document.documentElement.classList.add('dark')
} else {
    document.documentElement.classList.remove('dark')
}