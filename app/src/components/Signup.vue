<template>
    <div class="flex flex-col items-center justify-center min-h-screen bg-gray-100">
        <h1 class="text-gray-700 text-2xl font-semibold mb-2">Sign Up</h1>
        <p class="text-gray-500 text-sm mb-8">Have an account? <a class="text-blue-600 hover:text-blue-700" href="/login">Log In</a></p>
        <form class="w-full max-w-sm bg-white rounded-md px-8 pt-6 pb-8 mb-4" @submit.prevent="register">
            <div v-if="!apiSuccess">
                <div class="mb-4">
                    <label class="block text-gray-500 text-sm font-semibold mb-2" for="email">
                        Email Address
                    </label>
                    <input
                        v-model="email"
                        v-bind:class="{ 'border-red-500': emailError }"
                        class="appearance-none outline-none border-2 rounded-md w-full py-3 px-4 text-gray-700 leading-tight focus:border-blue-600 mb-2"
                        id="email" type="email" autocomplete="email">
                    <p v-if="emailError" class="text-red-500 text-sm">Required field</p>
                </div>
                <div class="mb-6">
                    <label class="block text-gray-500 text-sm font-semibold mb-2" for="password">
                        Password
                    </label>
                    <input
                        v-model="password"
                        v-bind:class="{ 'border-red-500': passwordError }"
                        class="appearance-none outline-none border-2 rounded-md w-full py-3 px-4 text-gray-700 leading-tight focus:border-blue-600 mb-2"
                        id="password" type="password" autocomplete="current-password">
                    <p v-if="passwordError" class="text-red-500 text-sm mb-2">Required field</p>
                    <p class="text-gray-500 text-sm mb-2">Must be at least 8 characters</p>
                </div>
                <div class="flex items-center justify-between">
                    <button
                        :disabled="isLoading"
                        class="bg-blue-600 hover:bg-blue-700 text-white font-medium py-3 px-4 rounded-md focus:outline-none focus:shadow-outline"
                        type="submit">
                        Sign Up
                    </button>
                </div>
                <p v-if="apiError" class="text-red-500 text-sm mt-6">{{ apiError }}</p>
            </div>
            <div v-if="apiSuccess">
                <p class="text-gray-500 text-sm mb-6">{{ apiSuccess }}</p>
                <a
                    href="/login"
                    class="bg-blue-600 hover:bg-blue-700 text-white font-medium py-3 px-4 rounded-md focus:outline-none focus:shadow-outline">
                    Proceed to Log In
                </a>
            </div>
        </form>
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
const apiSuccess = ref('')
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

const register = async () => {
    if (!validate()) {
        return
    }
    
    isLoading.value = true // Start loading
    const data = {
        email: email.value,
        password: password.value
    }

    try {
        const response = await userApi.register(data)
        apiSuccess.value = response.data.message
        apiError.value = ''
    } catch (error) {
        apiSuccess.value = ''
        if (axios.isAxiosError(error)) {
            apiError.value = error.response?.data.error
        }
    } finally {
        isLoading.value = false // End loading
    }
}
</script>
