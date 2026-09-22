import dropdown from '@preline/dropdown'
import overlay from '@preline/overlay'
import select from '@preline/select'
import tooltip from '@preline/tooltip'

// Preline's autoInit() attaches a *new anonymous* document/window listener on every call and
// never removes it, and its internal dedupe is an O(elements x collection) linear scan. Calling
// it per row therefore leaks hundreds of global handlers and costs O(n^2). These wrappers run it
// only when the DOM actually contains widgets Preline has not bound yet.

type Collection = { element?: { el?: Element } }[] | undefined

const collections = {
    dropdown: () => window.$hsDropdownCollection,
    overlay: () => window.$hsOverlayCollection,
    select: () => window.$hsSelectCollection,
    tooltip: () => window.$hsTooltipCollection
}

declare global {
    interface Window {
        $hsDropdownCollection?: Collection
        $hsOverlayCollection?: Collection
        $hsSelectCollection?: Collection
        $hsTooltipCollection?: Collection
    }
}

// Preline never drops entries for removed elements, so the collection grows without bound as
// rows are paginated, retaining detached DOM nodes and slowly re-inflating autoInit's cost.
const prune = (collection: Collection) => {
    if (!collection?.length) return

    for (let i = collection.length - 1; i >= 0; i--) {
        const el = collection[i]?.element?.el
        if (el && !document.contains(el)) collection.splice(i, 1)
    }
}

const hasUninitialized = (selector: string, collection: Collection) => {
    const elements = document.querySelectorAll(selector)
    if (!elements.length) return false
    if (!collection?.length) return true

    const known = new Set(collection.map(item => item?.element?.el))
    return Array.from(elements).some(el => !known.has(el))
}

const guard = (selector: string, getCollection: () => Collection, autoInit: () => void) => () => {
    const collection = getCollection()
    prune(collection)
    if (hasUninitialized(selector, collection)) autoInit()
}

export const initDropdowns = guard('.hs-dropdown', collections.dropdown, () => dropdown.autoInit())
export const initOverlays = guard('[data-hs-overlay]', collections.overlay, () => overlay.autoInit())
export const initSelects = guard('[data-hs-select]', collections.select, () => select.autoInit())
export const initTooltips = guard('.hs-tooltip', collections.tooltip, () => tooltip.autoInit())

// Needed by handlers that use @click.stop, which suppresses Preline's own window-level close.
// Preline only clears its internal `animationInProcess` flag from the menu's transitionend
// handler, so a menu that gets display:none'd mid-transition leaves the dropdown permanently
// stuck. Callers that hide the dropdown's container must pass withTransition = false.
export const closeDropdowns = (withTransition = true) => dropdown.closeCurrentlyOpened(null, withTransition)
