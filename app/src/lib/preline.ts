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

// Tooltips are the one widget that must not be bound up front. Preline builds a Popper instance
// the moment it binds a tooltip, and Popper registers a passive `scroll` listener per instance
// that re-measures the page on every scroll event. An alias row carries seven tooltips, so a
// 100-alias page installed ~700 of them: each scroll frame stalled the main thread while the
// compositor, which does not wait for it, scrolled on and left rows blank until it caught up.
// Binding on first hover keeps the behaviour identical - the tooltip opens on the same pointer
// movement it always did - but only the few tooltips actually pointed at ever build a Popper, and
// Preline drops their scroll listeners again on the mouseleave that hides them.
const boundTooltips = new WeakSet<Element>()
let tooltipsDelegated = false

const bindTooltip = (event: MouseEvent) => {
    const target = event.target
    if (!(target instanceof Element)) return

    const el = target.closest<HTMLElement>('.hs-tooltip')
    if (!el || boundTooltips.has(el)) return

    // The trigger can be a narrower element than the .hs-tooltip box. Waiting until the pointer is
    // actually on it guarantees the mouseleave that hides the tooltip - and with it the teardown of
    // Popper's scroll listeners - will follow.
    const toggle = el.querySelector<HTMLElement>('.hs-tooltip-toggle') ?? el
    if (!toggle.contains(target)) return

    // Several tooltip bodies are rendered conditionally. Preline drops an .hs-tooltip with no
    // content on the floor, so leave this one unmarked and look again once it has something to say.
    if (!el.querySelector('.hs-tooltip-content')) return

    boundTooltips.add(el)
    if (!window.$hsTooltipCollection) window.$hsTooltipCollection = []
    const collection = collections.tooltip()
    prune(collection)
    // Preline runs its own autoInit on window `load`, so anything on screen by then is already bound.
    if (collection?.some(item => item?.element?.el === el)) return

    new tooltip(el)
    // Preline opens a tooltip from a mouseenter on the toggle, which fired before the instance
    // existed, so open it here.
    tooltip.show(el)
}

// Idempotent: there is no per-render work left to do, only the one delegated listener to install.
export const initTooltips = () => {
    if (tooltipsDelegated) return
    document.addEventListener('mouseover', bindTooltip, true)
    tooltipsDelegated = true
}

// Importing @preline/tooltip registers `window.addEventListener('load', () => autoInit())`, which
// would bind every tooltip on screen at that moment - a whole page of alias rows, if the list won
// the race with the last subresource. Redirecting autoInit to the lazy binder closes that door and
// leaves the components that call it directly working exactly as before.
tooltip.autoInit = initTooltips

// Needed by handlers that use @click.stop, which suppresses Preline's own window-level close.
// Preline only clears its internal `animationInProcess` flag from the menu's transitionend
// handler, so a menu that gets display:none'd mid-transition leaves the dropdown permanently
// stuck. Callers that hide the dropdown's container must pass withTransition = false.
export const closeDropdowns = (withTransition = true) => dropdown.closeCurrentlyOpened(null, withTransition)
