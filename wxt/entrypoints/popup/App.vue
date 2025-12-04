<template>
  <LoginScreen v-if="!apiToken" />
  <AliasesScreen v-else-if="defaults" :apiToken="apiToken" :defaults="defaults" />
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
