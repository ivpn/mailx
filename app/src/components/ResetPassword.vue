<template>
    <div class="flex flex-col items-center justify-center min-h-screen bg-gray-100">
        <h1 class="text-3xl text-gray-800 font-semibold mb-8">Set New Password</h1>
        <form class="w-full max-w-sm bg-white px-8 pt-6 pb-8 mb-4" @submit.prevent="resetPassword">
            <div v-if="!apiSuccess">
                <div class="mb-6">
                    <label class="block text-gray-500 mb-2" for="password-new">
                        New Password
                    </label>
                    <input
                        v-model="password"
                        v-bind:class="{ 'border-gray-500': !passwordError, 'border-red-600': passwordError }"
                        class="appearance-none outline-none border w-full py-3 px-4 leading-tight focus:border-bluish-500 mb-2"
                        id="password-new" type="password">
                </div>
                <div class="mb-6">
                    <label class="block text-gray-500 mb-2" for="password-new-conmfirm">
                        Confirm
                    </label>
                    <input
                        v-model="passwordConfirm"
                        v-bind:class="{ 'border-gray-500': !passwordError, 'border-red-600': passwordError }"
                        class="appearance-none outline-none border w-full py-3 px-4 leading-tight focus:border-bluish-500 mb-2"
                        id="password-new-conmfirm" type="password">
                </div>
                <div class="flex items-center justify-between">
                    <button :disabled="isLoading"
                        class="bg-bluish-500 hover:bg-bluish-600 text-white font-medium py-3 px-4 focus:outline-none focus:shadow-outline"
                        type="submit">
                        Change Password
                    </button>
                </div>
                <p v-if="passwordError" class="text-red-600 text-sm mt-6">Error: {{ passwordError }}</p>
                <p v-if="apiError" class="text-red-600 text-sm mt-6">Error: {{ apiError }}</p>
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
        <p class="text-gray-500 my-5"><a class="text-bluish-500 hover:text-bluish-600"
            href="/login">Back to Log In</a></p>
    </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import axios from 'axios'
import { userApi } from '../api/user.ts'

const password = ref('')
const passwordConfirm = ref('')
const passwordError = ref('')
const apiError = ref('')
const apiSuccess = ref('')
const isLoading = ref(false)

const validatePassword = () => {
    apiSuccess.value = ''
    passwordError.value = ''

    if (!password.value || !passwordConfirm.value) {
        passwordError.value = 'Please fill required fields'
    }

    if (password.value !== passwordConfirm.value) {
        passwordError.value = 'Passwords do not match'
    }

    return !passwordError.value
}

const resetPassword = async () => {
    if (!validatePassword()) return

    isLoading.value = true // Start loading

    const req = {
        password: password.value
    }

    try {
        const res = await userApi.resetPassword(req)
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