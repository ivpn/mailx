<template>
    <div class="page center pt-10">
        <h1>MailX</h1>
        <h2 class="mb-10">
            Email forwarding service operated by
            <a href="https://www.ivpn.net/">IVPN</a>
        </h2>
        <h3>Sign Up</h3>
        <p>Have an account? <router-link to="/login">Log In</router-link></p>
        <form class="card center" @submit.prevent="">
            <div v-if="passkeySupported" class="border-b border-gray-200 dark:border-neutral-600">
                <nav class="flex gap-x-1" aria-label="Tabs" role="tablist" aria-orientation="horizontal">
                    <button
                        class="active"
                        id="tabs-with-underline-item-1" aria-selected="true" data-hs-tab="#tabs-with-underline-1"
                        aria-controls="tabs-with-underline-1" role="tab">
                        Passkey
                    </button>
                    <button
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
                            <label for="email_authn">
                                Email Address
                            </label>
                            <input
                                v-model="emailAuthn"
                                v-bind:class="{ 'error': emailAuthnError }"
                                placeholder="name@example.net"
                                id="email_authn"
                                type="email"
                            >
                            <p v-if="emailAuthnError" class="error">Required</p>
                        </div>
                        <div class="flex items-center w-full">
                            <button @click="registerWithPasskey" :disabled="isLoading" class="cta full">
                                Sign Up with Passkey
                            </button>
                        </div>
                        <p v-if="apiError" class="error mt-6">Error: {{ apiError }}</p>
                    </div>
                </div>
                <div
                    id="tabs-with-underline-2"
                    v-bind:class="{ 'hidden': passkeySupported }"
                    role="tabpanel"
                    aria-labelledby="tabs-with-underline-item-2">
                    <div v-if="!apiSuccess">
                        <div class="mb-4">
                            <label for="email">
                                Email Address
                            </label>
                            <input
                                v-model="email"
                                v-bind:class="{ 'error': emailError }"
                                placeholder="name@example.net"
                                id="email"
                                type="email"
                            >
                            <p v-if="emailError" class="error">Required</p>
                        </div>
                        <div class="mb-6">
                            <label for="password">
                                Password
                            </label>
                            <input
                                v-model="password"
                                v-bind:class="{ 'error': passwordError }"
                                id="password"
                                type="password"
                            >
                            <p v-if="passwordError" class="error">Required</p>
                            <p class="text-sm mb-2">Must be 12+ characters and contain uppercase, lowercase, number, and special character (e.g. !@#$%^&*(),;.?":{}|<>)</p>
                        </div>
                        <div class="flex items-center w-full">
                            <button @click="register" :disabled="isLoading" class="cta full">
                                Sign Up
                            </button>
                        </div>
                        <p v-if="apiError" class="error mt-6">Error: {{ apiError }}</p>
                    </div>
                </div>
            </div>
            <div v-if="apiSuccess">
                <p class="success mb-6">{{ apiSuccess }}</p>
                <router-link to="/login" tag="button" class="cta full">
                    Proceed to Log In
                </router-link>
            </div>
        </form>
        <Footer />
    </div>
</template>

<script setup lang="ts">
import { ref, onMounted, onUpdated } from 'vue'
import { useRoute } from 'vue-router'
import axios from 'axios'
import { userApi } from '../api/user.ts'
import { startRegistration, browserSupportsWebAuthn } from '@simplewebauthn/browser'
import tabs from '@preline/tabs'
import Footer from './Footer.vue'

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
const subid = ref('')

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
        password: password.value,
        subid: subid.value
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
        email: emailAuthn.value,
        subid: subid.value
    }

    try {
        var res = await userApi.registerBegin(data)
        const creds = await startRegistration({ optionsJSON: res.data['publicKey'] })
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

const parseSubid = () => {
    const route = useRoute()
    subid.value = route.params.subid as string
    if (!subid.value || !subid.value.match(/^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$/)) {
        window.location.href = '/login'
    }
}

onMounted(() => {
    parseSubid()
    passkeySupported.value = browserSupportsWebAuthn()
    tabs.autoInit()
})

onUpdated(() => {
    tabs.autoInit()
})
</script>
