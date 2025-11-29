export default defineContentScript({
  matches: ['*://*.mailx.net/*'],
  main() {
    // console.log('Hello content script', { id: browser.runtime.id })
  },
})
