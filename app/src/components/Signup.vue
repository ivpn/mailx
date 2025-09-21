<template>
    <div class="page center">
        <div></div>
        <form class="card-tertiary center" @submit.prevent="">
            <article>
                <div>
                    <div v-if="passkeySupported" id="tabs-with-underline-1" role="tabpanel" aria-labelledby="tabs-with-underline-item-1">
                        <h1 class="flex justify-center text-accent mb-8">
                            <span class="logo"></span>
                        </h1>
                        <h4 class="text-center mb-8">Sign up with Passkey</h4>
                        <div v-if="!apiSuccess">
                            <div class="mb-5">
                                <input
                                    v-model="emailAuthn"
                                    v-bind:class="{ 'error': emailAuthnError }"
                                    placeholder="Email Address"
                                    id="email_authn"
                                    type="email"
                                    class="email"
                                    @keypress.enter.prevent
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
                        <h1 class="flex justify-center text-accent mb-8">
                            <span class="logo"></span>
                        </h1>
                        <h4 class="text-center mb-8">Sign up with email and password</h4>
                        <div v-if="!apiSuccess">
                            <div class="mb-5">
                                <input
                                    v-model="email"
                                    v-bind:class="{ 'error': emailError }"
                                    placeholder="Email Address"
                                    id="email"
                                    type="email"
                                    class="email"
                                    @keypress.enter.prevent
                                >
                                <p v-if="emailError" class="error">Required</p>
                            </div>
                            <div class="mb-3">
                                <input
                                    v-model="password"
                                    v-bind:class="{ 'error': passwordError }"
                                    placeholder="Password"
                                    id="password"
                                    type="password"
                                    class="password"
                                    @keypress.enter.prevent
                                >
                                <p v-if="passwordError" class="error">Required</p>
                            </div>
                            <p class="text-sm mb-5">Must be 12+ characters and contain uppercase, lowercase, number, and special character (e.g. -_+=~!@#$%^&*(),;.?":{}|<>)</p>
                            <div class="flex items-center w-full">
                                <button @click="register" :disabled="isLoading" class="cta full">
                                    Sign Up
                                </button>
                            </div>
                            <p v-if="apiError" class="error mt-5">Error: {{ apiError }}</p>
                        </div>
                    </div>
                </div>
                <nav v-if="passkeySupported" aria-label="Tabs" role="tablist" aria-orientation="horizontal" class="tabs-router">
                    <button
                        class="active"
                        id="tabs-with-underline-item-1" aria-selected="true" data-hs-tab="#tabs-with-underline-1"
                        aria-controls="tabs-with-underline-1" role="tab">
                        Use Passkey instead
                    </button>
                    <button
                        id="tabs-with-underline-item-2" aria-selected="false" data-hs-tab="#tabs-with-underline-2"
                        aria-controls="tabs-with-underline-2" role="tab">
                        Use Password instead
                    </button>
                </nav>
                <div v-if="apiSuccess">
                    <p class="success mb-6">{{ apiSuccess }}</p>
                    <router-link to="/login" tag="button" class="cta full">
                        Proceed to Log In
                    </router-link>
                </div>
            </article>
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
const preauthid = ref('')
const preauthtokenhash = ref('')
const service = ref('')

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
        subid: subid.value,
        preauthid: preauthid.value,
        preauthtokenhash: preauthtokenhash.value,
        service: service.value
    }

    try {
        await userApi.register(data)
        apiError.value = ''
        window.location.href = '/signup-complete'
    } catch (err) {
        apiSuccess.value = ''
        if (axios.isAxiosError(err)) {
            apiError.value = err.response?.data.error || err.message

            if (err.response?.status === 429) {
                apiError.value = 'Too many requests, please try again later.'
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
        subid: subid.value,
        preauthid: preauthid.value,
        preauthtokenhash: preauthtokenhash.value,
        service: service.value
    }

    try {
        var res = await userApi.registerBegin(data)
        const creds = await startRegistration({ optionsJSON: res.data['publicKey'] })
        res = await userApi.registerFinish(creds)
        apiError.value = ''
        localStorage.setItem('email', data.email)
        window.location.href = '/'
    } catch (err) {
        if (axios.isAxiosError(err)) {
            apiError.value = err.response?.data.error || err.message

            if (err.response?.status === 429) {
                apiError.value = 'Too many requests, please try again later.'
            }
        }
    } finally {
        isLoading.value = false // End loading
    }
}

const parseParams = () => {
    const route = useRoute()
    const q = route.query
    const first = (v: unknown) => typeof v === 'string' ? v : Array.isArray(v) ? v[0] : ''
    subid.value = first(q.subid) || (route.params.subid as string) || ''
    preauthid.value = first(q.preauthid) || (route.params.preauthid as string) || ''
    preauthtokenhash.value = first(q.preauthtokenhash) || (route.params.preauthtokenhash as string) || ''
    service.value = first(q.service) || (route.params.service as string) || ''
    preauthtokenhash.value = preauthtokenhash.value.replace(/ /g, '+')

    if (!subid.value || !subid.value.match(/^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$/)) {
        console.error('Invalid or missing subid')
    }

    if (!preauthid.value || !preauthid.value.match(/^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$/)) {
        console.error('Invalid or missing preauthid')
    }

    if (!preauthtokenhash.value) {
        console.error('Invalid or missing preauthtokenhash')
    }

    if (!service.value) {
        console.error('Invalid or missing service')
    }
}

const isLoggedIn = (): boolean => {
    const email = localStorage.getItem('email')
    return email !== null && email.trim() !== ''
}

onMounted(() => {
    if (isLoggedIn()) {
        window.location.href = '/'
    }
    
    parseParams()
    passkeySupported.value = browserSupportsWebAuthn()
    tabs.autoInit()
})

onUpdated(() => {
    tabs.autoInit()
})
</script>
