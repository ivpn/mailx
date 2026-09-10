<template>
    <!-- No visible trigger: this overlay is opened programmatically, but Preline only registers
         an instance for elements with a data-hs-overlay trigger, so a hidden one is required -->
    <button type="button" data-hs-overlay="#modal-stepup-auth" class="hidden" aria-hidden="true" tabindex="-1"></button>
    <div id="modal-stepup-auth" class="hs-overlay hidden">
        <div>
            <div>
                <header>
                    <button @click="cancel" class="close">
                        <i class="icon arrow-left-line icon-primary"></i>
                    </button>
                    <h4>VERIFY IT'S YOU</h4>
                </header>
                <article>
                    <div class="mb-5">
                        <p>
                            For your security, please confirm your identity to continue.
                        </p>
                    </div>
                    <div v-if="stepUpMethods.passkey && passkeySupported" class="mb-5">
                        <button :disabled="isLoading" @click="verifyWithPasskey" class="cta full">
                            Use Passkey
                        </button>
                    </div>
                    <div v-if="stepUpMethods.password" class="mb-7">
                        <label for="stepup-password">
                            Password:
                        </label>
                        <input
                            v-model="password"
                            v-bind:class="{ 'error': passwordError }"
                            id="stepup-password"
                            type="password"
                            autocomplete="current-password"
                        >
                        <p v-if="passwordError" class="error">Required</p>
                    </div>
                </article>
                <footer>
                    <nav>
                        <button v-if="stepUpMethods.password" :disabled="isLoading" @click="verifyWithPassword" class="cta">
                            Confirm
                        </button>
                        <button @click="cancel" class="cancel">
                            Cancel
                        </button>
                    </nav>
                    <p v-if="error" class="error px-5">Error: {{ error }}</p>
                </footer>
            </div>
        </div>
    </div>
</template>

<script setup lang="ts">
import { ref, watch, onMounted } from 'vue'
import { userApi } from '../api/user.ts'
import axios from 'axios'
import overlay from '@preline/overlay'
import { startAuthentication, browserSupportsWebAuthn } from '@simplewebauthn/browser'
import { stepUpIsOpen, stepUpMethods, completeStepUp, cancelStepUp } from '../stepup.ts'

const password = ref('')
const passwordError = ref(false)
const error = ref('')
const isLoading = ref(false)
const passkeySupported = ref(false)

const modalEl = () => document.querySelector('#modal-stepup-auth') as HTMLElement

const reset = () => {
    password.value = ''
    passwordError.value = false
    error.value = ''
}

const cancel = () => {
    reset()
    cancelStepUp()
}

const validatePassword = () => {
    passwordError.value = !password.value
    return !passwordError.value
}

const verifyWithPassword = async () => {
    error.value = ''
    if (!validatePassword()) return

    isLoading.value = true
    try {
        await userApi.stepUpPassword({ password: password.value })
        reset()
        completeStepUp()
    } catch (err) {
        if (axios.isAxiosError(err)) {
            error.value = err.response?.data.error || err.message

            if (err.response?.status === 429) {
                error.value = 'Too many requests, please try again later.'
            }
        }
    } finally {
        isLoading.value = false
    }
}

const verifyWithPasskey = async () => {
    error.value = ''
    isLoading.value = true
    try {
        const res = await userApi.stepUpPasskeyBegin()
        const creds = await startAuthentication({ optionsJSON: res.data['publicKey'] })
        await userApi.stepUpPasskeyFinish(creds)
        reset()
        completeStepUp()
    } catch (err) {
        if (axios.isAxiosError(err)) {
            error.value = err.response?.data.error || err.message

            if (err.response?.status === 429) {
                error.value = 'Too many requests, please try again later.'
            }
        } else {
            error.value = 'The operation was aborted or failed.'
        }
    } finally {
        isLoading.value = false
    }
}

watch(stepUpIsOpen, (open) => {
    if (open) {
        overlay.open(modalEl())
    } else {
        overlay.close(modalEl())
    }
})

const addEvents = () => {
    const modal = overlay.getInstance('#modal-stepup-auth' as any, true) as any
    // Covers dismissal via backdrop click / Escape, not just our own Cancel button
    modal.element.on('close', () => {
        cancel()
    })
}

onMounted(() => {
    overlay.autoInit()
    passkeySupported.value = browserSupportsWebAuthn()
    addEvents()
})
</script>
