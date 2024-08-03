<template>
    <div class="mb-5">
        <h2 class="text-2xl font-semibold text-gray-800 dark:text-gray-100 mb-5">Change Password</h2>
        <div class="mb-4 max-w-xs">
            <label class="block text-gray-500 mb-3" for="new-password">
                New password:
            </label>
            <input v-model="password"
                v-bind:class="{ 'border-gray-500': !passwordError, 'border-red-600': passwordError }"
                class="appearance-none outline-none border w-full py-3 px-4 text-gray-500 leading-tight focus:border-bluish-500 mb-2"
                id="new-password" type="password">
        </div>
        <div class="mb-4 max-w-xs">
            <label class="block text-gray-500 mb-3" for="new-password-confirm">
                Confirm new password:
            </label>
            <input v-model="passwordConfirm"
                v-bind:class="{ 'border-gray-500': !passwordError, 'border-red-600': passwordError }"
                class="appearance-none outline-none border w-full py-3 px-4 text-gray-500 leading-tight focus:border-bluish-500 mb-2"
                id="new-password-confirm" type="password">
        </div>
        <div class="mb-3 max-w-xs">
            <button @click="changePassword"
                class="bg-bluish-500 hover:bg-bluish-600 text-white font-medium py-2 px-3 focus:outline-none focus:shadow-outline"
                type="submit">
                Change Password
            </button>
        </div>
        <p v-if="passwordError" class="text-red-600 text-sm mb-3">Error: {{ passwordError }}</p>
        <p v-if="error" class="text-red-600 text-sm mb-3">Error: {{ error }}</p>
        <p v-if="success" class="text-emerald-600 dark:text-emerald-500 text-sm mb-3">{{ success }}</p>
    </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { userApi } from '../api/user.ts'
import axios from 'axios'

const password = ref('')
const passwordConfirm = ref('')
const passwordError = ref('')
const error = ref('')
const success = ref('')

const validatePassword = () => {
    success.value = ''
    passwordError.value = ''

    if (!password.value || !passwordConfirm.value) {
        passwordError.value = 'Please fill required fields'
    }

    if (password.value !== passwordConfirm.value) {
        passwordError.value = 'Passwords do not match'
    }

    return !passwordError.value
}

const changePassword = async () => {
    if (!validatePassword()) return

    const req = {
        password: password.value
    }

    try {
        const res = await userApi.changePassword(req)
        success.value = res.data.message
        error.value = ''
        password.value = ''
        passwordConfirm.value = ''
    } catch (err) {
        if (axios.isAxiosError(err)) {
            success.value = ''
            error.value = err.response?.data.error || err.message

            if (err.response?.status === 429) {
                error.value = 'Too many requests, please try again later'
            }
        }
    }
}
</script>