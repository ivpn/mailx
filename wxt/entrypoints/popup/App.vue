<template>
  <LoginScreen />
</template>

<script lang="ts" setup>
import { onMounted } from 'vue'
import LoginScreen from '@/components/popup/LoginScreen.vue'
import { ref } from 'vue'
import { store } from '@/lib/store'

const apiToken = ref<string | undefined>()

onMounted(async () => {
  apiToken.value = await store.getApiToken()

  store.onApiTokenChange((newToken) => {
    apiToken.value = newToken
  })
})
</script>

<style scoped></style>
