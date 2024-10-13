<template>
    <div class="flex flex-col items-center justify-center min-h-screen bg-gray-100 dark:bg-neutral-900">
        <h1 class="text-3xl text-gray-800 dark:text-gray-100 font-semibold mb-2">Sign Up</h1>
        <p class="text-gray-500 dark:text-gray-400 mb-8">Have an account? <a class="text-bluish-500 hover:text-bluish-600" href="/login">Log In</a></p>
        <form class="w-full max-w-sm bg-white dark:bg-neutral-800 px-8 pt-6 pb-8 mb-4" @submit.prevent="">
            <div v-if="passkeySupported" class="border-b border-gray-200 dark:border-neutral-600">
                <nav class="flex gap-x-1" aria-label="Tabs" role="tablist" aria-orientation="horizontal">
                    <button type="button"
                        class="hs-tab-active:font-semibold hs-tab-active:border-bluish-500 hs-tab-active:text-bluish-500 pt-2 pb-4 px-1 text-center basis-0 grow inline-flex justify-center items-center gap-x-2 border-b-2 border-transparent whitespace-nowrap text-gray-500 hover:text-bluish-500 focus:outline-none focus:text-bluish-500 disabled:opacity-50 disabled:pointer-events-none dark:text-neutral-400 dark:hover:text-bluish-500 active"
                        id="tabs-with-underline-item-1" aria-selected="true" data-hs-tab="#tabs-with-underline-1"
                        aria-controls="tabs-with-underline-1" role="tab">
                        Passkey
                    </button>
                    <button type="button"
                        class="hs-tab-active:font-semibold hs-tab-active:border-bluish-500 hs-tab-active:text-bluish-500 pt-2 pb-4 px-1 text-center basis-0 grow inline-flex justify-center items-center gap-x-2 border-b-2 border-transparent whitespace-nowrap text-gray-500 hover:text-bluish-500 focus:outline-none focus:text-bluish-500 disabled:opacity-50 disabled:pointer-events-none dark:text-neutral-400 dark:hover:text-bluish-500"
                        id="tabs-with-underline-item-2" aria-selected="false" data-hs-tab="#tabs-with-underline-2"
                        aria-controls="tabs-with-underline-2" role="tab">
                        Email & Password
                    </button>
                </nav>
            </div>
            <div v-bind:class="{ 'mt-6': passkeySupported }">
                <div v-if="passkeySupported" id="tabs-with-underline-1" role="tabpanel" aria-labelledby="tabs-with-underline-item-1">
                    <div v-if="!apiSuccess">
                        <div class="mb-4">
                            <label class="block text-gray-500 dark:text-gray-400 mb-2" for="email_authn">
                                Email Address
                            </label>
                            <input v-model="emailAuthn"
                                v-bind:class="{ 'border-gray-500': !emailAuthnError, 'border-red-600 dark:border-red-600': emailAuthnError }"
                                placeholder="name@example.net"
                                class="appearance-none outline-none border w-full py-3 px-4 leading-tight focus:border-bluish-500 mb-2 dark:bg-neutral-800 dark:border-gray-500 dark:text-gray-400"
                                id="email_authn" type="email">
                            <p v-if="emailAuthnError" class="text-red-600 text-sm">Required</p>
                        </div>
                        <div class="flex items-center w-full">
                            <button :disabled="isLoading"
                                class="w-full bg-bluish-500 hover:bg-bluish-600 text-white font-medium py-3 px-4 focus:outline-none focus:shadow-outline"
                                type="button" @click="registerWithPasskey">
                                Sign Up with Passkey
                            </button>
                        </div>
                        <p v-if="apiError" class="text-red-600 text-sm mt-6">Error: {{ apiError }}</p>
                    </div>
                </div>
                <div
                    id="tabs-with-underline-2"
                    v-bind:class="{ 'hidden': passkeySupported }"
                    role="tabpanel"
                    aria-labelledby="tabs-with-underline-item-2">
                    <div v-if="!apiSuccess">
                        <div class="mb-4">
                            <label class="block text-gray-500 dark:text-gray-400 mb-2" for="email">
                                Email Address
                            </label>
                            <input v-model="email"
                                v-bind:class="{ 'border-gray-500': !emailError, 'border-red-600 dark:border-red-600': emailError }"
                                placeholder="name@example.net"
                                class="appearance-none outline-none border w-full py-3 px-4 leading-tight focus:border-bluish-500 mb-2 dark:bg-neutral-800 dark:border-gray-500 dark:text-gray-400"
                                id="email" type="email">
                            <p v-if="emailError" class="text-red-600 text-sm">Required</p>
                        </div>
                        <div class="mb-6">
                            <label class="block text-gray-500 dark:text-gray-400 mb-2" for="password">
                                Password
                            </label>
                            <input v-model="password"
                                v-bind:class="{ 'border-gray-500': !passwordError, 'border-red-600 dark:border-red-600': passwordError }"
                                class="appearance-none outline-none border w-full py-3 px-4 leading-tight focus:border-bluish-500 mb-2 dark:bg-neutral-800 dark:border-gray-500 dark:text-gray-400"
                                id="password" type="password">
                            <p v-if="passwordError" class="text-red-600 text-sm mb-2">Required</p>
                            <p class="text-gray-500 dark:text-gray-400 text-sm mb-2">Must be 12+ characters and contain uppercase, lowercase, number, and special character (e.g. !@#$%^&*(),;.?":{}|<>)</p>
                        </div>
                        <div class="flex items-center w-full">
                            <button :disabled="isLoading"
                                class="w-full bg-bluish-500 hover:bg-bluish-600 text-white font-medium py-3 px-4 focus:outline-none focus:shadow-outline"
                                type="button" @click="register">
                                Sign Up
                            </button>
                        </div>
                        <p v-if="apiError" class="text-red-600 text-sm mt-6">Error: {{ apiError }}</p>
                    </div>
                </div>
            </div>
            <div v-if="apiSuccess">
                <p class="text-emerald-600 dark:text-emerald-500 text-sm mb-6">{{ apiSuccess }}</p>
                <a href="/login"
                    class="bg-bluish-500 hover:bg-bluish-600 text-white font-medium py-3 px-4 focus:outline-none focus:shadow-outline">
                    Proceed to Log In
                </a>
            </div>
        </form>
    </div>
</template>

<script setup lang="ts">
import { ref, onMounted, onUpdated } from 'vue'
import axios from 'axios'
import { userApi } from '../api/user.ts'
import { startRegistration, browserSupportsWebAuthn } from '@simplewebauthn/browser'
import tabs from '@preline/tabs'

const email = ref('')
const emailAuthn = ref('')
const password = ref('')
const emailError = ref(false)
const emailAuthnError = ref(false)
const passwordError = ref(false)
const apiSuccess = ref('')
const apiError = ref('')
const isLoading = ref(false)
const passkeySupported = ref(false)

const validateEmail = () => {
    emailError.value = !email.value
    return !emailError.value
}

const validateEmailAuthn = () => {
    emailAuthnError.value = !emailAuthn.value
    return !emailAuthnError.value
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
    if (!validate()) return

    isLoading.value = true // Start loading
    const data = {
        email: email.value,
        password: password.value
    }

    try {
        const res = await userApi.register(data)
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

const registerWithPasskey = async () => {
    if (!validateEmailAuthn()) return

    isLoading.value = true // Start loading

    const data = {
        email: emailAuthn.value
    }

    try {
        var res = await userApi.registerBegin(data)
        const creds = await startRegistration(res.data['publicKey'])
        res = await userApi.registerFinish(creds)
        apiSuccess.value = res.data.message
        apiError.value = ''
        if (res.status === 200) {
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

onMounted(() => {
    passkeySupported.value = browserSupportsWebAuthn()
    tabs.autoInit()
})

onUpdated(() => {
    tabs.autoInit()
})
</script>
