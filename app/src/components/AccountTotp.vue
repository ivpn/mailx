<template>
    <div v-if="res.id" class="mb-5">
        <h2 class="text-2xl font-semibold text-gray-800 mb-5">2-Factor Authentication</h2>
        <p class="text-sm text-gray-500 mb-5">
            <span v-if="res.totp_enabled"
                class="inline-flex items-center py-1.5 px-2 rounded-md text-xs font-medium bg-emerald-100 text-emerald-800">Enabled</span>
            <span v-if="!res.totp_enabled"
                class="inline-flex items-center py-1.5 px-2 rounded-md text-xs font-medium bg-gray-100 text-gray-500">Disabled</span>
        </p>
        <p class="text-gray-500 mb-5">
            When enabled, 2-factor authentication will be required when you log in.<br>
        </p>
        <div class="mb-3 max-w-xs">
            <AccountTotpEnable v-if="!res.totp_enabled" @onTotpEnable="reload" />
            <AccountTotpDisable v-if="res.totp_enabled" />
        </div>
        <p v-if="error" class="text-red-600 text-sm mb-3">Error: {{ error }}</p>
    </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { userApi } from '../api/user.ts'
import axios from 'axios'
import AccountTotpEnable from './AccountTotpEnable.vue'
import AccountTotpDisable from './AccountTotpDisable.vue'

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

const reload = () => {
    getUser()
}

onMounted(() => {
    reload()
})
</script>