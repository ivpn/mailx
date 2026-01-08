import { store } from '@/lib/store'
import { Defaults } from '@/lib/types'

let apiToken: string | undefined
let defaults: Defaults | undefined

export default defineContentScript({
  matches: ['<all_urls>'],
  runAt: 'document_start',
  async main() {
    apiToken = await store.getApiToken()
    defaults = await store.getDefaults()
    observeEmailInputs()
  },
})

function observeEmailInputs() {
  const observer = new MutationObserver(() => {
    const inputs = document.querySelectorAll<HTMLInputElement>(
      'input[type="email"]:not([data-alias-injected])'
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
  input.dataset.aliasInjected = 'true'

  // Ensure parent is positioned
  const parent = input.parentElement
  if (!parent) return

  if (getComputedStyle(parent).position === 'static') {
    parent.style.position = 'relative'
  }

  const host = document.createElement('div')
  host.style.position = 'absolute'
  host.style.right = '8px'
  host.style.top = '50%'
  host.style.transform = 'translateY(-50%)'
  host.style.zIndex = '9999'
  host.style.width = '24px'
  host.style.height = '24px'

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
    backgroundSize: '18px',
    backgroundRepeat: 'no-repeat',
    backgroundPosition: 'center',
    backgroundColor: '#2c2c2c',
  })

  button.addEventListener('click', (e) => {
    e.preventDefault()
    e.stopPropagation()
    generateAliasFor(input)
  })

  shadow.appendChild(button)
  parent.appendChild(host)
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
