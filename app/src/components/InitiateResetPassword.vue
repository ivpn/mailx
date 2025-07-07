<template>
    <div class="page center">
        <div></div>
        <form class="card-tertiary center" @submit.prevent="initiatePasswordReset">
            <article>
                <h1 class="flex justify-center text-accent mb-8">
                    <span class="logo"></span>
                </h1>
                <h4 class="text-center mb-8">Reset password</h4>
                <div v-if="!apiSuccess">
                    <div class="mb-3">
                        <input
                            v-model="email"
                            v-bind:class="{ 'error': emailError }"
                            id="email"
                            type="email"
                            autocomplete="email"
                            placeholder="Email Address"
                            class="email"
                        >
                        <p v-if="emailError" class="error">Required</p>
                    </div>
                    <p class="text-sm mb-5">
                        Please enter your registered email address. You will be sent instructions on how to reset your password.
                    </p>
                    <div class="flex items-center justify-between">
                        <button :disabled="isLoading" class="cta full">
                            Send reset instructions
                        </button>
                    </div>
                    <p v-if="apiError" class="error mt-6">Error: {{ apiError }}</p>
                </div>
                <div v-if="apiSuccess">
                    <p>If an account with the specified email address exists we will send an email with further instructions on how to reset your password.</p>
                </div>
                <nav role="tablist" class="tabs-router">
                    <router-link to="/login" custom v-slot="{ navigate }">
                        <button @click="navigate">Back to Log In</button>
                    </router-link>
                </nav>
            </article>
        </form>
        <Footer />
    </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import axios from 'axios'
import { userApi } from '../api/user.ts'
import Footer from './Footer.vue'

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
                apiError.value = 'Too many requests, please try again later.'
            }
        }
    } finally {
        isLoading.value = false // End loading
    }
}

</script>