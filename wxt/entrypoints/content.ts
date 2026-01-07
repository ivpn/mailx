export default defineContentScript({
  matches: ['<all_urls>'],
  runAt: 'document_idle',
  main() {
    // console.log('Hello content script', { id: browser.runtime.id })
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

  observer.observe(document.body, {
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
  host.style.zIndex = '9999'

  const shadow = host.attachShadow({ mode: 'closed' })

  const button = document.createElement('button')
  button.textContent = '📨'
  button.title = 'Generate email alias'

  Object.assign(button.style, {
    cursor: 'pointer',
    border: 'none',
    background: 'transparent',
    fontSize: '16px',
  })

  button.addEventListener('click', (e) => {
    e.preventDefault()
    e.stopPropagation()
    generateAliasFor(input)
  })

  shadow.appendChild(button)
  parent.appendChild(host)
}

function generateAliasFor(input: HTMLInputElement) {
  const alias = generateAlias()
  input.value = alias

  // Trigger input events so frameworks (React/Vue) notice
  input.dispatchEvent(new Event('input', { bubbles: true }))
  input.dispatchEvent(new Event('change', { bubbles: true }))
}

function generateAlias() {
  const site = location.hostname.replace('www.', '')
  const random = Math.random().toString(36).slice(2, 8)
  return `user+${site}-${random}@example.com`
}
