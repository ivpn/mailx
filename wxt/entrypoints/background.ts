import { api } from '@/lib/api'

export default defineBackground(() => {
  browser.runtime.onMessage.addListener(
    (msg, _sender, sendResponse) => {
      if (msg?.type !== 'CREATE_ALIAS') return

      const { apiToken, alias } = msg.payload ?? {}

      if (!apiToken || !alias) {
        sendResponse({ ok: false, error: 'Invalid payload' })
        return true
      }

      ;(async () => {
        try {
          const res = await api.createAlias(apiToken, alias)
          console.log('[BG] Created alias:', res)
          sendResponse({ ok: true, result: res })
        } catch (err) {
          sendResponse({
            ok: false,
            error: err instanceof Error ? err.message : String(err),
          })
        }
      })()

      // 🔑 THIS LINE IS CRITICAL
      return true
    }
  )
})
