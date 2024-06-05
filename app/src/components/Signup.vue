<template>
    <div class="flex flex-col items-center justify-center min-h-screen bg-gray-100">
        <h1 class="text-3xl text-gray-800 font-semibold mb-2">Sign Up</h1>
        <p class="text-gray-500 mb-8">Have an account? <a class="text-bluish-500 hover:text-bluish-600" href="/login">Log In</a></p>
        <form class="w-full max-w-sm bg-white px-8 pt-6 pb-8 mb-4" @submit.prevent="register">
            <div v-if="!apiSuccess">
                <div class="mb-4">
                    <label class="block text-gray-500 mb-2" for="email">
                        Email Address
                    </label>
                    <input
                        v-model="email"
                        v-bind:class="{ 'border-red-600': emailError }"
                        placeholder="name@example.net"
                        class="appearance-none outline-none border border-gray-500 w-full py-3 px-4 leading-tight focus:border-bluish-500 mb-2"
                        id="email" type="email">
                    <p v-if="emailError" class="text-red-600 text-sm">Required</p>
                </div>
                <div class="mb-6">
                    <label class="block text-gray-500 mb-2" for="password">
                        Password
                    </label>
                    <input
                        v-model="password"
                        v-bind:class="{ 'border-red-600': passwordError }"
                        class="appearance-none outline-none border border-gray-500 w-full py-3 px-4 leading-tight focus:border-bluish-500 mb-2"
                        id="password" type="password">
                    <p v-if="passwordError" class="text-red-600 text-sm mb-2">Required</p>
                    <p class="text-gray-500 text-sm mb-2">Must be at least 8 characters</p>
                </div>
                <div class="flex items-center justify-between">
                    <button
                        :disabled="isLoading"
                        class="bg-bluish-500 hover:bg-bluish-600 text-white font-medium py-3 px-4 focus:outline-none focus:shadow-outline"
                        type="submit">
                        Sign Up
                    </button>
                </div>
                <p v-if="apiError" class="text-red-600 text-sm mt-6">{{ apiError }}</p>
            </div>
            <div v-if="apiSuccess">
                <p class="text-emerald-600 text-sm mb-6">{{ apiSuccess }}</p>
                <a
                    href="/login"
                    class="bg-bluish-500 hover:bg-bluish-600 text-white font-medium py-3 px-4 focus:outline-none focus:shadow-outline">
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
    } catch (err) {
        apiSuccess.value = ''
        if (axios.isAxiosError(err)) {
            apiError.value = err.response?.data.error || err.message
        }
    } finally {
        isLoading.value = false // End loading
    }
}
</script>
