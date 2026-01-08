import { defineConfig } from 'wxt'

// See https://wxt.dev/api/config.html
export default defineConfig({
  modules: ['@wxt-dev/module-vue'],
  manifest: ({ mode }) => {
    return {
      host_permissions: mode === 'development' ? ['http://0.0.0.0:3000/*'] : [],
      permissions: ['storage', 'webRequest', 'activeTab'],
      web_accessible_resources: [
        {
          resources: ['mailx.svg'],
          matches: ['<all_urls>'],
        },
      ],
    }
  },
})