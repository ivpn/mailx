<template>
    <div></div>
</template>

<script setup lang="ts">
import { onMounted } from 'vue'
import { userApi } from '../api/user.ts'

onMounted(async () => {
    try {
        await userApi.logout()
    } catch { }
    userApi.clearSession()
})
</script>
