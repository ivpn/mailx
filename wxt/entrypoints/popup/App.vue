<template>
  <LoginScreen />
</template>

<script lang="ts" setup>
import { onMounted } from 'vue'
import { api } from '@/lib/api'
import LoginScreen from '@/components/popup/LoginScreen.vue'

onMounted(async () => {
  try {
    const res = await api.livez()
    console.log('API response:', res)
  } catch (error) {
    console.error('Failed to check API status:', error)
  }
})
</script>

<style scoped></style>
