<template>
    <h2 class="text-2xl font-semibold text-gray-800 dark:text-gray-100 mb-5">Delete Account</h2>
    <p class="text-gray-500 mb-3">
        Are you sure you want to delete your account? This action cannot be undone.
    </p>
    <div class="mb-4 max-w-xs">
        <label class="block text-gray-500 mb-3" for="account-password">
            Password:
        </label>
        <input v-model="password" v-bind:class="{ 'border-gray-500': !passwordError, 'border-red-600': passwordError }"
            class="appearance-none outline-none border w-full py-3 px-4 text-gray-500 leading-tight focus:border-bluish-500 mb-2"
            id="account-password" type="password">
        <p v-if="passwordError" class="text-red-600 text-sm mb-2">Required</p>
    </div>
    <div class="mb-3 max-w-xs">
        <button @click="promptDeleteAccount"
            class="bg-red-600 hover:bg-red-700 text-white font-medium py-2 px-3 focus:outline-none focus:shadow-outline"
            type="submit">
            Delete Account
        </button>
    </div>
    <p v-if="error" class="text-red-600 text-sm mb-3">Error: {{ error }}</p>
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

            if (err.response?.status === 429) {
                error.value = 'Too many requests, please try again later'
            }
        }
    }
}
</script>