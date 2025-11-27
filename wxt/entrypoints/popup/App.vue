<template>
  <LoginScreen v-if="!apiToken" />
  <AliasesScreen v-else :apiToken="apiToken" />
</template>

<script lang="ts" setup>
import { ref, onMounted } from 'vue'
import { store } from '@/lib/store'
import LoginScreen from '@/components/popup/LoginScreen.vue'
import AliasesScreen from '@/components/popup/AliasesScreen.vue'

const apiToken = ref<string | undefined>()

onMounted(async () => {
  apiToken.value = await store.getApiToken()

  store.onApiTokenChange((newToken) => {
    console.log('API token changed:', newToken)
    apiToken.value = newToken
  })
})
</script>
