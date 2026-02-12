import { defineConfig } from 'wxt'

// See https://wxt.dev/api/config.html
export default defineConfig({
  modules: ['@wxt-dev/module-vue'],
  manifest: () => {
    const apiUrl = import.meta.env.VITE_API_URL
    console.log('API URL:', apiUrl)
    return {
      host_permissions: apiUrl ? [apiUrl + '/*'] : [],
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