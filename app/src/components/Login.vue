<template>
    <div class="flex flex-col items-center justify-center min-h-screen bg-gray-100">
        <h1 class="text-3xl text-gray-800 font-semibold mb-2">Log In</h1>
        <p class="text-gray-500 mb-8">Need an account? <a class="text-bluish-500 hover:text-bluish-600"
                href="/signup">Sign Up</a></p>
        <form class="w-full max-w-sm bg-white px-8 pt-6 pb-8 mb-4" @submit.prevent="login">
            <div v-if="!isLoggedIn()">
                <div class="mb-4">
                    <label class="block text-gray-500 mb-2" for="email">
                        Email Address
                    </label>
                    <input
                        v-model="email"
                        v-bind:class="{ 'border-gray-500': !emailError, 'border-red-600': emailError }"
                        class="appearance-none outline-none border w-full py-3 px-4 leading-tight focus:border-bluish-500 mb-2"
                        id="email" type="email" autocomplete="email">
                    <p v-if="emailError" class="text-red-600 text-sm">Required</p>
                </div>
                <div class="mb-6">
                    <label class="block text-gray-500 mb-2" for="password">
                        Password
                    </label>
                    <input
                        v-model="password"
                        v-bind:class="{ 'border-gray-500': !passwordError, 'border-red-600': passwordError }"
                        class="appearance-none outline-none border w-full py-3 px-4 leading-tight focus:border-bluish-500 mb-2"
                        id="password" type="password" autocomplete="current-password">
                    <p v-if="passwordError" class="text-red-600 text-sm mb-2">Required</p>
                </div>
                <div class="flex items-center justify-between">
                    <button :disabled="isLoading"
                        class="bg-bluish-500 hover:bg-bluish-600 text-white font-medium py-3 px-4 focus:outline-none focus:shadow-outline"
                        type="submit">
                        Log In
                    </button>
                </div>
                <p v-if="apiError" class="text-red-600 text-sm mt-6">Error: {{ apiError }}</p>
            </div>
            <div v-if="isLoggedIn()" class="pb-2">
                <p class="text-gray-500 mb-6">You are logged in</p>
                <a href="/"
                    class="bg-bluish-500 hover:bg-bluish-600 text-white font-medium py-3 px-4 focus:outline-none focus:shadow-outline">
                    Go to Dashboard
                </a>
            </div>
        </form>
        <p class="text-gray-500 my-5"><a class="text-bluish-500 hover:text-bluish-600"
            href="/initiate-password-reset">Forgot Your Password?</a></p>
    </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import axios from 'axios'
import { userApi } from '../api/user.ts'

const email = ref('')
const password = ref('')
const emailError = ref(false)
const passwordError = ref(false)
const apiError = ref('')
const isLoading = ref(false)

const validateEmail = () => {
    emailError.value = !email.value
    return !emailError.value
}

const validatePassword = () => {
    passwordError.value = !password.value
    return !passwordError.value
}

const validate = () => {
    const validEmail = validateEmail()
    const validPass = validatePassword()
    return validEmail && validPass
}

const login = async () => {
    if (!validate()) return

    isLoading.value = true // Start loading
    const data = {
        email: email.value,
        password: password.value
    }

    try {
        const response = await userApi.login(data)
        apiError.value = ''
        if (response.status === 200) {
            // Redirect to the dashboard
            localStorage.setItem('email', data.email)
            window.location.href = '/'
        }
    } catch (err) {
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

const isLoggedIn = () => {
    return localStorage.getItem('email') != '' && localStorage.getItem('email') != null
}
</script>
