<template>
    <div class="p-4">
        <h2 class="m-0">Settings</h2>
        <hr class="my-5">
        <p class="text-sm mb-3">Refresh recipients, domains and defaults:</p>
        <button @click="refreshDefaults" class="cta sm">Refresh Defaults</button>
        <hr class="my-5">
        <p class="text-sm mb-3">Log out / delete session:</p>
        <button @click="logout" class="cta sm">Log Out</button>
        <p v-if="error" class="error my-5">Error: {{ error }}</p>
        <p v-if="success" class="success my-5">{{ success }}</p>
    </div>
</template>

<script lang="ts" setup>
import { api } from '@/lib/api'
import { store } from '@/lib/store'

const props = defineProps<{
    apiToken: string
}>()

const success = ref('')
const error = ref('')

const refreshDefaults = async () => {
    try {
        const res = await api.fetchDefaults(props.apiToken)
        console.log('Fetched defaults:', res)
        processResponse(res)
        success.value = 'Defaults refreshed successfully'
        error.value = ''
    } catch (err) {
        console.error('Fetch defaults error:', err)
        success.value = ''
        error.value = 'An unexpected error occurred'
    }
}

const logout = async () => {
    if (!confirm('Do you want to proceed with the logout?')) return

    try {
        await api.logout(props.apiToken)
        store.clearAll()
        console.log('Logged out successfully')
        error.value = ''
        alert('You have been logged out.')
    } catch (err) {
        console.error('Logout error:', err)
        success.value = ''
        error.value = 'An unexpected error occurred during logout'
    }
}

const processResponse = (res: any) => {
    const defaults = {
        domain: res.domain,
        domains: res.domains,
        recipient: res.recipient,
        recipients: res.recipients,
        alias_format: res.alias_format,
    }
    store.setDefaults(defaults)
}
</script>