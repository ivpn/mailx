import { store } from '@/lib/store'
import { Defaults, Preferences } from '@/lib/types'

let apiToken: string | undefined
let defaults: Defaults | undefined
let preferences: Preferences | undefined

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
  const observer = new MutationObserver(() => {
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

  // Measure the input
  // const rect = input.getBoundingClientRect()
  // const height = rect.height
  // const width = rect.width

  const parent = input.parentNode
  if (!parent) return

  // Create wrapper
  const wrapper = document.createElement('div')
  Object.assign(wrapper.style, {
    position: 'relative',
  })

  // Move input into wrapper
  parent.insertBefore(wrapper, input)
  wrapper.appendChild(input)

  // Button host
  const host = document.createElement('div')
  Object.assign(host.style, {
    position: 'absolute',
    top: '50%',
    right: '8px',
    transform: 'translateY(-50%)',
    width: '24px',
    height: '24px',
    pointerEvents: 'auto',
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
  wrapper.appendChild(host)
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
