<template>
  <div class="max-w-[420px] min-w-[420px]">
    <LoginScreen v-if="!apiToken" />
    <AliasesScreen v-if="apiToken && defaults" :apiToken="apiToken" :defaults="defaults" />
    <header v-if="apiToken && defaults" class="bg-secondary fixed bottom-0 left-0 right-0 z-10">
      <nav>
          <div class="flex flex-row items-center">
              <a v-bind:class="{ 'active': true }" href="/">
                  Aliases
              </a>
              <a v-bind:class="{ 'active': false }" href="/settings">
                  Settings
              </a>
          </div>
      </nav>
    </header>
  </div>
</template>

<script lang="ts" setup>
import { ref, onMounted } from 'vue'
import { api } from '@/lib/api'
import { store } from '@/lib/store'
import { Defaults } from '@/lib/types'
import LoginScreen from '@/components/popup/LoginScreen.vue'
import AliasesScreen from '@/components/popup/AliasesScreen.vue'

const apiToken = ref<string | undefined>()
const defaults = ref<Defaults | undefined>()

onMounted(async () => {
  apiToken.value = await store.getApiToken()
  defaults.value = await store.getDefaults()

  store.onApiTokenChange((newToken) => {
    console.log('API token changed:', newToken)
    apiToken.value = newToken
  })

  store.onDefaultsChange((newDefaults) => {
    console.log('Defaults changed:', newDefaults)
    defaults.value = newDefaults
  })

  try {
    const res = await api.livez()
    console.log('/livez response:', res)
  } catch (err) {
    console.error('Error calling /livez:', err)
  }
})
</script>
