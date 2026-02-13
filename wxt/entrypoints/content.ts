import { store } from '@/lib/store'
import { Defaults, Preferences } from '@/lib/types'

let apiToken: string | undefined
let defaults: Defaults | undefined
let preferences: Preferences | undefined

// Track button hosts and position update intervals for cleanup
const buttonHosts = new WeakMap<HTMLInputElement, HTMLDivElement>()
const updateIntervals = new WeakMap<HTMLInputElement, ReturnType<typeof setInterval>>()
const updateFunctions = new WeakMap<HTMLInputElement, () => void>()

export default defineContentScript({
  matches: ['<all_urls>'],
  runAt: 'document_start',
  async main() {
    apiToken = await store.getApiToken()
    defaults = await store.getDefaults()
    preferences = await store.getPreferences()
    if (preferences.input_button && apiToken) {
      observeEmailInputs()
    }
  },
})

function observeEmailInputs() {
  const observer = new MutationObserver((mutations) => {
    // Handle added nodes
    const inputs = document.querySelectorAll<HTMLInputElement>(
      `
      input[type="email"]:not([data-alias-injected]),
      input[type="text"][name*="email" i]:not([data-alias-injected]),
      input[type="text"][id*="email" i]:not([data-alias-injected]),
      input[type="email"][name*="email" i]:not([data-alias-injected]),
      input[type="email"][id*="email" i]:not([data-alias-injected])
      `
    )

    inputs.forEach(injectButton)

    // Handle removed nodes - cleanup buttons for removed inputs
    mutations.forEach((mutation) => {
      mutation.removedNodes.forEach((node) => {
        if (node.nodeType === Node.ELEMENT_NODE) {
          const element = node as Element
          // Check if removed node is a tracked input
          if (element instanceof HTMLInputElement && buttonHosts.has(element)) {
            cleanupButton(element)
          }
          // Check descendants
          if (element.querySelectorAll) {
            const trackedInputs = element.querySelectorAll<HTMLInputElement>('input[data-alias-injected="true"]')
            trackedInputs.forEach(cleanupButton)
          }
        }
      })
    })
  })

  const root = document.documentElement ?? document
  observer.observe(root, {
    childList: true,
    subtree: true,
  })
}

function injectButton(input: HTMLInputElement) {
  if (input.dataset.aliasInjected === 'true') return
  if (!isValidEmailInput(input)) return
  input.dataset.aliasInjected = 'true'

  // Button host - appended to body with absolute positioning
  const host = document.createElement('div')
  Object.assign(host.style, {
    position: 'absolute',
    left: '0',
    top: '0',
    width: '24px',
    height: '24px',
    pointerEvents: 'auto',
    zIndex: '10',
  })

  const icon = browser.runtime.getURL('/mailx.svg')
  const shadow = host.attachShadow({ mode: 'closed' })

  const button = document.createElement('button')
  button.title = 'Create Mailx alias'

  Object.assign(button.style, {
    width: '24px',
    height: '24px',
    border: 'none',
    padding: '0',
    cursor: 'pointer',
    borderRadius: '50%',
    backgroundImage: `url(${icon})`,
    backgroundSize: '16px',
    backgroundRepeat: 'no-repeat',
    backgroundPosition: 'center',
    backgroundColor: '#2c2c2c',
    transition: 'transform 0.2s ease',
  })

  button.addEventListener('click', (e) => {
    e.preventDefault()
    e.stopPropagation()
    generateAliasFor(input)
  })

  button.addEventListener('mouseenter', () => {
    button.style.transform = 'scale(1.15)';
  })

  button.addEventListener('mouseleave', () => {
    button.style.transform = 'scale(1)';
  })

  shadow.appendChild(button)
  document.body.appendChild(host)

  // Store reference and start position tracking
  buttonHosts.set(input, host)
  positionButtonRelativeToInput(input, host)
}

function positionButtonRelativeToInput(input: HTMLInputElement, host: HTMLDivElement) {
  function updatePosition() {
    // Check if input still exists in DOM
    if (!document.body.contains(input)) {
      cleanupButton(input)
      return
    }

    // Get input dimensions and position
    const inputRect = input.getBoundingClientRect()
    const inputStyle = getComputedStyle(input)
    
    // Check if input is still visible
    if (inputRect.width === 0 || inputRect.height === 0) {
      host.style.display = 'none'
      return
    }
    
    host.style.display = 'block'

    // Calculate button position (middle-right of input)
    const buttonWidth = 24
    const buttonHeight = 24
    const rightOffset = 10 // 10px from right edge of input
    
    // Use absolute positioning with page offsets
    const left = inputRect.right + window.pageXOffset - buttonWidth - rightOffset
    const top = inputRect.top + window.pageYOffset + (inputRect.height - buttonHeight) / 2

    host.style.left = `${left}px`
    host.style.top = `${top}px`
  }

  // Store update function for cleanup
  updateFunctions.set(input, updatePosition)

  // Update immediately
  updatePosition()

  // Update periodically (handles dynamic changes)
  const intervalId = setInterval(updatePosition, 200)
  updateIntervals.set(input, intervalId)
  
  // Add scroll and resize listeners for immediate updates
  window.addEventListener('scroll', updatePosition, true)
  window.addEventListener('resize', updatePosition)
}

function cleanupButton(input: HTMLInputElement) {
  // Remove event listeners
  const updateFunction = updateFunctions.get(input)
  if (updateFunction) {
    window.removeEventListener('scroll', updateFunction, true)
    window.removeEventListener('resize', updateFunction)
    updateFunctions.delete(input)
  }

  // Clear position update interval
  const intervalId = updateIntervals.get(input)
  if (intervalId !== undefined) {
    clearInterval(intervalId)
    updateIntervals.delete(input)
  }

  // Remove button host from DOM
  const host = buttonHosts.get(input)
  if (host && host.parentNode) {
    host.parentNode.removeChild(host)
  }
  buttonHosts.delete(input)

  // Clear injected marker
  input.dataset.aliasInjected = 'false'
}

function isValidEmailInput(element: HTMLInputElement): boolean {
  const style = getComputedStyle(element);

  return (
    // visible & interactive
    style.visibility !== "hidden" &&
    style.display !== "none" &&
    style.opacity !== "0" &&
    style.pointerEvents !== "none" &&

    // usable
    !element.disabled &&
    !element.readOnly &&

    // avoid false positives (checkboxes, radios, etc.)
    (element.type === "email" || element.type === "text")
  )
}

async function generateAliasFor(input: HTMLInputElement) {
  const alias = await postAlias()
  if (!alias) return
  input.value = alias

  // Trigger input events so frameworks (React/Vue) notice
  input.dispatchEvent(new Event('input', { bubbles: true }))
  input.dispatchEvent(new Event('change', { bubbles: true }))
}

type CreateAliasResponse =
  | { ok: true; result: { alias: { name: string } } }
  | { ok: false; error: string }

async function postAlias(): Promise<string | undefined> {
  if (!defaults || !apiToken) {
    console.warn('[CS] Missing defaults or apiToken')
    return
  }

  const alias = {
    domain: defaults.domain,
    recipients: defaults.recipient,
    format: defaults.alias_format,
    enabled: true,
  }

  let res: CreateAliasResponse | undefined

  try {
    res = await browser.runtime.sendMessage({
      type: 'CREATE_ALIAS',
      payload: { apiToken, alias },
    })
  } catch (err) {
    // This catches channel-level failures (SW died, extension reloaded, etc.)
    console.error('[CS] Message send failed:', err)
    return
  }

  if (!res) {
    console.error('[CS] No response from background')
    return
  }

  if (!res.ok) {
    console.error('[CS] Create alias error:', res.error)
    return
  }

  if (!res.result.alias.name) {
    console.error('[CS] Invalid success response:', res)
    return
  }

  return res.result.alias.name
}
