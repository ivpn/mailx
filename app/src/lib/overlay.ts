// @preline/overlay >= 5 no longer default-exports HSOverlay from its entry file: index.js only
// registers the class as window.HSOverlay (and installs its load/resize handlers), while its
// typings still declare a default export. vite.config.ts aliases the bare '@preline/overlay'
// import to this module so `import overlay from '@preline/overlay'` keeps working everywhere,
// including vi.mock('@preline/overlay') in specs. Remove once the package exports the class again.
import '@preline/overlay/index.js'
import type HSOverlay from '@preline/overlay/index.js'

declare global {
    interface Window {
        HSOverlay: typeof HSOverlay
    }
}

export default window.HSOverlay
