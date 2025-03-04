<template>
    <div class="flex flex-col items-center justify-center min-h-screen bg-gray-100 dark:bg-neutral-900">
        <h1>Reset Password</h1>
        <form class="w-full max-w-sm bg-white dark:bg-neutral-800 px-8 pt-6 pb-8 mb-4" @submit.prevent="initiatePasswordReset">
            <div v-if="!apiSuccess">
                <p>Please enter your registered email address. You will be sent instructions on how to reset your password.</p>
                <div class="mb-4">
                    <label for="email">
                        Email Address
                    </label>
                    <input
                        v-model="email"
                        v-bind:class="{ 'border-gray-500': !emailError, 'border-red-600 dark:border-red-600': emailError }"
                        class="appearance-none outline-none border w-full py-3 px-4 leading-tight focus:border-bluish-500 mb-2 text-gray-500 bg-white dark:text-gray-300 dark:bg-neutral-800"
                        id="email" type="email" autocomplete="email">
                    <p v-if="emailError" class="error">Required</p>
                </div>
                <div class="flex items-center justify-between">
                    <button :disabled="isLoading"
                        class="bg-bluish-500 hover:bg-bluish-600 text-white font-medium py-3 px-4 focus:outline-none focus:shadow-outline"
                        type="submit">
                        Send Reset Instructions
                    </button>
                </div>
                <p v-if="apiError" class="error mt-6">Error: {{ apiError }}</p>
            </div>
            <div v-if="apiSuccess">
                <p>If an account with the specified email address exists we will send an email with further instructions on how to reset your password.</p>
            </div>
        </form>
        <p class="text-gray-500 my-5">
            <router-link to="/login">Back to Log In</router-link>
        </p>
    </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import axios from 'axios'
import { userApi } from '../api/user.ts'

const email = ref('')
const emailError = ref(false)
const apiError = ref('')
const apiSuccess = ref('')
const isLoading = ref(false)

const validateEmail = () => {
    emailError.value = !email.value
    return !emailError.value
}

const initiatePasswordReset = async () => {
    if (!validateEmail()) return

    isLoading.value = true // Start loading

    const req = {
        email: email.value
    }

    try {
        const res = await userApi.initiatePasswordReset(req)
        apiSuccess.value = res.data.message
        apiError.value = ''
    } catch (err) {
        apiSuccess.value = ''
        if (axios.isAxiosError(err)) {
            apiError.value = err.response?.data.error || err.message

            if (err.response?.status === 429) {
                apiError.value = 'Too many requests, please try again later'
            }
        }
    } finally {
        isLoading.value = false // End loading
    }
}

</script>