<template>
    <div class="page center pt-10">
        <div></div>
        <form class="card-tertiary center" @submit.prevent="" autocomplete="off">
            <article>
                <div>
                    <div v-if="passkeySupported" id="tabs-with-underline-1" role="tabpanel" aria-labelledby="tabs-with-underline-item-1">
                        <h1 class="text-center text-accent mb-8">MailX</h1>
                        <h4 class="text-center mb-8">Log in with Passkey</h4>
                        <div v-if="!isLoggedIn()">
                            <div class="mb-5">
                                <input
                                    v-model="emailAuthn"
                                    v-bind:class="{ 'error': emailAuthnError }"
                                    id="email_authn"
                                    type="email"
                                    placeholder="Email Address"
                                    class="email"
                                    autocomplete="false"
                                >
                                <p v-if="emailAuthnError" class="error">Required</p>
                            </div>
                            <div class="flex items-center w-full">
                                <button :disabled="isLoading" @click="loginWithPasskey" class="cta lg full">
                                    Log in with Passkey
                                </button>
                            </div>
                            <p v-if="error" class="error mt-6">Error: {{ error }}</p>
                        </div>
                    </div>
                    <div id="tabs-with-underline-2" v-bind:class="{ 'hidden': passkeySupported }" role="tabpanel"
                        aria-labelledby="tabs-with-underline-item-2">
                        <div v-if="!isLoggedIn()">
                            <h1 class="text-center text-accent mb-8">MailX</h1>
                            <h4 class="text-center mb-8">Log in with email and password</h4>
                            <div class="mb-3">
                                <input
                                    v-model="email"
                                    v-bind:class="{ 'error': emailError }"
                                    id="email"
                                    type="email"
                                    autocomplete="email"
                                    placeholder="Email address"
                                    class="email"
                                >
                                <p v-if="emailError" class="error">Required</p>
                            </div>
                            <div class="mb-5">
                                <input
                                    v-model="password"
                                    v-bind:class="{ 'error': passwordError }"
                                    id="password"
                                    type="password"
                                    autocomplete="current-password"
                                    placeholder="Password"
                                    class="password"
                                >
                                <p v-if="passwordError" class="error mb-2">Required</p>
                                <p class="text-right">
                                    <router-link to="/reset/password/initiate">
                                        <button class="plain-alt">Forgot password?</button>
                                    </router-link>
                                </p>
                            </div>
                            <div v-if="otpRequired" class="mb-6">
                                <label for="password">
                                    Two-factor authentication token:
                                </label>
                                <input
                                    v-model="otp"
                                    v-bind:class="{ 'error': otpError }"
                                    id="otp"
                                    type="text"
                                >
                                <p v-if="otpError" class="error">Required</p>
                            </div>
                            <div class="flex items-center w-full">
                                <button :disabled="isLoading" @click="login" class="cta lg full">
                                    Log in
                                </button>
                            </div>
                            <p v-if="error" class="error mt-6">Error: {{ error }}</p>
                        </div>
                    </div>
                </div>
                <nav v-if="passkeySupported" aria-label="Tabs" role="tablist" aria-orientation="horizontal" class="tabs-router">
                    <button
                        @click="onTabChange"
                        class="active"
                        id="tabs-with-underline-item-1" aria-selected="true" data-hs-tab="#tabs-with-underline-1"
                        aria-controls="tabs-with-underline-1" role="tab">
                        Or use Passkey
                    </button>
                    <button
                        @click="onTabChange"
                        id="tabs-with-underline-item-2" aria-selected="false" data-hs-tab="#tabs-with-underline-2"
                        aria-controls="tabs-with-underline-2" role="tab">
                        Or use password
                    </button>
                </nav>
                <div v-if="isLoggedIn()" class="pb-2">
                    <p>You are logged in</p>
                    <router-link to="/" tag="button" class="cta lg full">
                        Go to Dashboard
                    </router-link>
                </div>
            </article>
            <div class="footer">
                <div>
                    <i class="icon info icon-primary"></i>
                </div>
                <div>
                    <h4>Here to try MailX? You need an active IVPN account.</h4>
                    <p>Sign up or log in on <a href="https://www.ivpn.net/account/">ivpn.net</a> and look for "Email Beta" in your account settings.</p>
                </div>
            </div>
        </form>
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
        const creds = await startAuthentication({ optionsJSON: res.data['publicKey'] })
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
