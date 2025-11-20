<template>
    <div v-if="res.id" class="mb-5">
        <h2>2-Factor Authentication</h2>
        <p class="text-sm">
            <span v-if="res.totp_enabled" class="badge success">Enabled</span>
            <span v-if="!res.totp_enabled" class="badge">Disabled</span>
        </p>
        <p>
            When enabled, 2-factor authentication will be required when you log in with your password.<br>
        </p>
        <div class="mb-3 max-w-xs">
            <AccountTotpEnable v-if="!res.totp_enabled" />
            <AccountTotpDisable v-if="res.totp_enabled" />
        </div>
        <p v-if="error" class="error">Error: {{ error }}</p>
    </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { userApi } from '../api/user.ts'
import axios from 'axios'
import AccountTotpEnable from './AccountTotpEnable.vue'
import AccountTotpDisable from './AccountTotpDisable.vue'
import events from '../events.ts'

const res = ref({ 
    id: '',
    totp_enabled: false
})
const error = ref('')

const getUser = async () => {
    try {
        const response = await userApi.get()
        res.value = response.data
    } catch (err) {
        if (axios.isAxiosError(err)) {
            error.value = err.message
        }
    }
}

onMounted(() => {
    getUser()
    events.on('totp.enable', getUser)
    events.on('totp.disable', getUser)
})
</script>