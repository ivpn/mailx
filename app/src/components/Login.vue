<template>
    <div class="flex flex-col items-center justify-center min-h-screen bg-gray-100 dark:bg-neutral-900">
        <h1 class="text-3xl text-gray-800 dark:text-gray-100 font-semibold my-8">MailX</h1>
        <h2 class="text-2xl text-gray-800 dark:text-gray-100 mb-8">
            Email forwarding service operated by
            <a class="text-bluish-500 hover:text-bluish-600" href="https://www.ivpn.net/">IVPN</a>
        </h2>
        <h2 class="text-xl text-gray-500 dark:text-gray-400 mb-8">Log In</h2>
        <form class="w-full max-w-sm bg-white dark:bg-neutral-800 px-8 pt-6 pb-8 mb-4" @submit.prevent="">
            <div v-if="passkeySupported" class="border-b border-gray-200 dark:border-neutral-600">
                <nav class="flex gap-x-1" aria-label="Tabs" role="tablist" aria-orientation="horizontal">
                    <button
                        type="button"
                        @click="onTabChange"
                        class="hs-tab-active:font-semibold hs-tab-active:border-bluish-500 hs-tab-active:text-bluish-500 pt-2 pb-4 px-1 text-center basis-0 grow inline-flex justify-center items-center gap-x-2 border-b-2 border-transparent whitespace-nowrap text-gray-500 hover:text-bluish-500 focus:outline-none focus:text-bluish-500 disabled:opacity-50 disabled:pointer-events-none dark:text-neutral-400 dark:hover:text-bluish-500 active"
                        id="tabs-with-underline-item-1" aria-selected="true" data-hs-tab="#tabs-with-underline-1"
                        aria-controls="tabs-with-underline-1" role="tab">
                        Passkey
                    </button>
                    <button
                        type="button"
                        @click="onTabChange"
                        class="hs-tab-active:font-semibold hs-tab-active:border-bluish-500 hs-tab-active:text-bluish-500 pt-2 pb-4 px-1 text-center basis-0 grow inline-flex justify-center items-center gap-x-2 border-b-2 border-transparent whitespace-nowrap text-gray-500 hover:text-bluish-500 focus:outline-none focus:text-bluish-500 disabled:opacity-50 disabled:pointer-events-none dark:text-neutral-400 dark:hover:text-bluish-500"
                        id="tabs-with-underline-item-2" aria-selected="false" data-hs-tab="#tabs-with-underline-2"
                        aria-controls="tabs-with-underline-2" role="tab">
                        Email & Password
                    </button>
                </nav>
            </div>
            <div v-bind:class="{ 'mt-6': passkeySupported }">
                <div v-if="passkeySupported" id="tabs-with-underline-1" role="tabpanel"
                    aria-labelledby="tabs-with-underline-item-1">
                    <div v-if="!isLoggedIn()">
                        <div class="mb-4">
                            <label class="block text-gray-500 dark:text-gray-400 mb-2" for="email_authn">
                                Email Address
                            </label>
                            <input v-model="emailAuthn"
                                v-bind:class="{ 'border-gray-500': !emailAuthnError, 'border-red-600 dark:border-red-600': emailAuthnError }"
                                class="appearance-none outline-none border w-full py-3 px-4 leading-tight focus:border-bluish-500 mb-2 dark:bg-neutral-800 dark:border-neutral-500 dark:text-gray-400"
                                id="email_authn" type="email" autocomplete="email_authn">
                            <p v-if="emailAuthnError" class="text-red-600 text-sm">Required</p>
                        </div>
                        <div class="flex items-center w-full">
                            <button :disabled="isLoading"
                                class="w-full bg-bluish-500 hover:bg-bluish-600 text-white font-medium py-3 px-4 focus:outline-none focus:shadow-outline"
                                type="button" @click="loginWithPasskey">
                                Log In with Passkey
                            </button>
                        </div>
                        <p v-if="error" class="text-red-600 text-sm mt-6">Error: {{ error }}</p>
                    </div>
                </div>
                <div id="tabs-with-underline-2" v-bind:class="{ 'hidden': passkeySupported }" role="tabpanel"
                    aria-labelledby="tabs-with-underline-item-2">
                    <div v-if="!isLoggedIn()">
                        <div class="mb-4">
                            <label class="block text-gray-500 dark:text-gray-400 mb-2" for="email">
                                Email Address
                            </label>
                            <input v-model="email"
                                v-bind:class="{ 'border-gray-500': !emailError, 'border-red-600 dark:border-red-600': emailError }"
                                class="appearance-none outline-none border w-full py-3 px-4 leading-tight focus:border-bluish-500 mb-2 dark:bg-neutral-800 dark:border-neutral-500 dark:text-gray-400"
                                id="email" type="email" autocomplete="email">
                            <p v-if="emailError" class="text-red-600 text-sm">Required</p>
                        </div>
                        <div class="mb-6">
                            <label class="block text-gray-500 dark:text-gray-400 mb-2" for="password">
                                Password
                            </label>
                            <input v-model="password"
                                v-bind:class="{ 'border-gray-500': !passwordError, 'border-red-600 dark:border-red-600': passwordError }"
                                class="appearance-none outline-none border w-full py-3 px-4 leading-tight focus:border-bluish-500 mb-2 dark:bg-neutral-800 dark:border-neutral-500 dark:text-gray-400"
                                id="password" type="password" autocomplete="current-password">
                            <p v-if="passwordError" class="text-red-600 text-sm mb-2">Required</p>
                        </div>
                        <div v-if="otpRequired" class="mb-6">
                            <label class="block text-gray-500 dark:text-gray-400 mb-2" for="password">
                                Two-factor authentication token:
                            </label>
                            <input v-model="otp"
                                v-bind:class="{ 'border-gray-500': !otpError, 'border-red-600 dark:border-red-600': otpError }"
                                class="appearance-none outline-none border w-full py-3 px-4 leading-tight focus:border-bluish-500 mb-2 dark:bg-neutral-800 dark:border-neutral-500 dark:text-gray-400"
                                id="otp" type="text">
                            <p v-if="otpError" class="text-red-600 text-sm mb-2">Required</p>
                        </div>
                        <div class="flex items-center w-full">
                            <button :disabled="isLoading"
                                class="w-full bg-bluish-500 hover:bg-bluish-600 text-white font-medium py-3 px-4 focus:outline-none focus:shadow-outline"
                                type="button" @click="login">
                                Log In
                            </button>
                        </div>
                        <p v-if="error" class="text-red-600 text-sm mt-6">Error: {{ error }}</p>
                    </div>
                </div>
            </div>
            <div v-if="isLoggedIn()" class="pb-2">
                <p class="text-gray-500 dark:text-gray-400 mb-6">You are logged in</p>
                <a href="/"
                    class="bg-bluish-500 hover:bg-bluish-600 text-white font-medium py-3 px-4 focus:outline-none focus:shadow-outline">
                    Go to Dashboard
                </a>
            </div>
        </form>
        <p class="text-gray-500 dark:text-gray-400 my-3">
            <a class="text-bluish-500 hover:text-bluish-600" href="/reset/password/initiate">Forgot Your Password?</a>
        </p>
        <Footer />
    </div>
</template>

<script setup lang="ts">
import { onMounted, onUpdated, ref } from 'vue'
import axios from 'axios'
import { userApi } from '../api/user.ts'
import { startAuthentication, browserSupportsWebAuthn } from '@simplewebauthn/browser'
import tabs from '@preline/tabs'
import Footer from './Footer.vue'

const email = ref('')
const emailAuthn = ref('')
const password = ref('')
const otp = ref('')
const emailError = ref(false)
const emailAuthnError = ref(false)
const passwordError = ref(false)
const otpError = ref(false)
const otpRequired = ref(false)
const error = ref('')
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

const validateOtp = () => {
    otpError.value = otpRequired.value && !otp.value
    return !otpError.value
}

const validate = () => {
    const validEmail = validateEmail()
    const validPass = validatePassword()
    const validotp = validateOtp()
    return validEmail && validPass && validotp
}

const login = async () => {
    if (!validate()) return

    isLoading.value = true // Start loading

    const data = {
        email: email.value,
        password: password.value,
        otp: otp.value
    }

    try {
        const response = await userApi.login(data)
        error.value = ''
        if (response.status === 200) {
            // Redirect to the dashboard
            localStorage.setItem('email', data.email)
            window.location.href = '/'
        }
    } catch (err) {
        if (axios.isAxiosError(err)) {
            error.value = err.response?.data.error || err.message

            if (err.response?.status === 429) {
                error.value = 'Too many requests, please try again later'
            }

            if (err.response?.data.code === 70001) {
                error.value = ''
                otpRequired.value = true
            }
        }
    } finally {
        isLoading.value = false // End loading
    }
}

const loginWithPasskey = async () => {
    if (!validateEmailAuthn()) return

    isLoading.value = true // Start loading

    const data = {
        email: emailAuthn.value
    }

    try {
        var res = await userApi.loginBegin(data)
        startAuth(data, res)
    } catch (err) {
        if (axios.isAxiosError(err)) {
            error.value = err.response?.data.error || err.message

            if (err.response?.status === 429) {
                error.value = 'Too many requests, please try again later'
            }
        }
    } finally {
        isLoading.value = false // End loading
    }
}

const startAuth = async (data: any, res: any) => {
    try {
        const creds = await startAuthentication(res.data['publicKey'])
        res = await userApi.loginFinish(creds)
        error.value = ''
        if (res.status === 200) {
            // Redirect to the dashboard
            localStorage.setItem('email', data.email)
            window.location.href = '/'
        }
    } catch (err: Error) {
        if (axios.isAxiosError(err)) {
            error.value = err.response?.data.error || err.message

            if (err.response?.status === 429) {
                error.value = 'Too many requests, please try again later'
            }
        } else {
            error.value = 'The operation was aborted or failed'
        }
    } finally {
        isLoading.value = false // End loading
    }
}

const isLoggedIn = () => {
    return localStorage.getItem('email') != '' && localStorage.getItem('email') != null
}

const onTabChange = () => {
    otpRequired.value = false
}

onMounted(() => {
    if (isLoggedIn()) {
        window.location.href = '/'
    }

    passkeySupported.value = browserSupportsWebAuthn()
    tabs.autoInit()
})

onUpdated(() => {
    tabs.autoInit()
})
</script>
