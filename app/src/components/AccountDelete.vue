<template>
    <div class="flex flex-col bg-white border border-red-600 shadow-sm rounded-xl p-5 pb-4 my-8">
        <h1 class="text-xl font-bold text-gray-800 mb-4">Delete Account</h1>
        <p class="text-sm text-gray-500 mb-3">
            Are you sure you want to delete your account? This action cannot be undone.
        </p>
        <div class="mb-4 max-w-xs">
            <label class="block text-gray-500 text-sm font-semibold mb-3" for="account-password">
                Password:
            </label>
            <input
                v-model="password"
                v-bind:class="{ 'border-red-600': passwordError }"
                class="appearance-none outline-none border-2 rounded-md w-full py-3 px-4 text-gray-700 leading-tight focus:border-bluish-500 mb-2"
                id="account-password" type="password">
            <p v-if="passwordError" class="text-red-600 text-sm mb-2">Required</p>
        </div>
        <div class="mb-3 max-w-xs">
            <button
                @click="promptDeleteAccount"
                class="bg-red-600 hover:bg-red-700 text-white font-medium text-sm py-2 px-3 rounded-md focus:outline-none focus:shadow-outline"
                type="submit">
                Delete Account
            </button>
        </div>
        <p v-if="error" class="text-red-600 text-sm mb-3">Error: {{ error }}</p>
    </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { userApi } from '../api/user.ts'
import axios from 'axios'

const password = ref('')
const passwordError = ref(false)
const error = ref('')

const validatePassword = () => {
    passwordError.value = !password.value
    return !passwordError.value
}

const promptDeleteAccount = () => {
    if (!validatePassword()) return
    if (!confirm('Are you sure you want to delete your account? This action cannot be undone.')) return
    deleteAccount()
}

const deleteAccount = async () => {
    const req = {
        password: password.value
    }

    try {
        await userApi.delete(req)
        alert('Account is deleted successfully. You will be logged out.')
        userApi.clearSession()
    } catch (err) {
        if (axios.isAxiosError(err)) {
            error.value = err.response?.data.error || err.message
        }
    }
}
</script>