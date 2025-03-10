<template>
    <div class="page center pt-10">
        <h1>Reset Password</h1>
        <form class="card center" @submit.prevent="initiatePasswordReset">
            <div v-if="!apiSuccess">
                <p>Please enter your registered email address. You will be sent instructions on how to reset your password.</p>
                <div class="mb-4">
                    <label for="email">
                        Email Address
                    </label>
                    <input
                        v-model="email"
                        v-bind:class="{ 'error': emailError }"
                        id="email"
                        type="email"
                        autocomplete="email"
                    >
                    <p v-if="emailError" class="error">Required</p>
                </div>
                <div class="flex items-center justify-between">
                    <button :disabled="isLoading" class="cta full">
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