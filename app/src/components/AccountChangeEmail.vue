<template>
    <div class="mb-5">
        <h2 class="text-2xl font-semibold text-gray-800 dark:text-gray-100 mb-5">Change Email</h2>
        <div v-if="!success">
            <div class="mb-4 max-w-xs">
                <label class="block text-gray-500 dark:text-gray-400 mb-3" for="new-email">
                    New email:
                </label>
                <input v-model="email"
                    v-bind:class="{ 'border-gray-500 dark:border-neutral-400': !emailError, 'border-red-600 dark:border-red-600': emailError }"
                    class="appearance-none outline-none border w-full py-3 px-4 text-gray-500 bg-white dark:text-gray-300 dark:bg-neutral-800 leading-tight focus:border-bluish-500 mb-2"
                    id="new-email" type="email">
            </div>
            <div class="mb-3 max-w-xs">
                <button @click="changeEmail"
                    class="bg-bluish-500 hover:bg-bluish-600 text-white font-medium py-2 px-3 focus:outline-none focus:shadow-outline"
                    type="submit">
                    Change Email
                </button>
            </div>
        </div>
        <p v-if="emailError" class="text-red-600 text-sm mb-3">Error: {{ emailError }}</p>
        <p v-if="error" class="text-red-600 text-sm mb-3">Error: {{ error }}</p>
        <p v-if="success" class="text-emerald-600 dark:text-emerald-500 text-sm mb-3">{{ success }}</p>
    </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { userApi } from '../api/user.ts'
import axios from 'axios'
import events from '../events.ts'

const email = ref('')
const emailError = ref('')
const error = ref('')
const success = ref('')

const validateEmail = () => {
    success.value = ''
    emailError.value = ''

    if (!email.value) {
        emailError.value = 'Required'
    }

    return !emailError.value
}

const changeEmail = async () => {
    if (!validateEmail()) return

    const req = {
        email: email.value
    }

    try {
        const res = await userApi.changeEmail(req)
        localStorage.setItem('email', req.email)
        events.emit('user.update', { email: req.email })
        success.value = res.data.message
        error.value = ''
        email.value = ''
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