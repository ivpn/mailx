<template>
    <div class="mb-5">
        <h2>Change Primary Email</h2>
        <div v-if="!success">
            <div class="mb-4 max-w-xs">
                <label for="new-email">
                    Set a new primary email address:
                </label>
                <input
                    v-model="email"
                    v-bind:class="{ 'error': emailError }"
                    id="new-email"
                    type="email"
                >
            </div>
            <div class="mb-3 max-w-xs">
                <button
                    @click="changeEmail"
                    class="cta">
                    Change Email
                </button>
            </div>
        </div>
        <p v-if="emailError" class="error mb-3">Error: {{ emailError }}</p>
        <p v-if="error" class="error mb-3">Error: {{ error }}</p>
        <p v-if="success" class="success text-sm mb-3">{{ success }}</p>
    </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { userApi } from '../api/user.ts'
import axios from 'axios'
import events from '../events.ts'

const email = ref('')
const emailError = ref('')
const error = ref('')
const success = ref('')

const validateEmail = () => {
    success.value = ''
    emailError.value = ''

    if (!email.value) {
        emailError.value = 'Required'
    }

    if (email.value == localStorage.getItem('email')) {
        emailError.value = 'The new email matches the current one'
    }

    return !emailError.value
}

const changeEmail = async () => {
    if (!validateEmail()) return

    const req = {
        email: email.value
    }

    try {
        const res = await userApi.changeEmail(req)
        localStorage.setItem('email', req.email)
        events.emit('user.update', { email: req.email })
        success.value = res.data.message
        error.value = ''
        email.value = ''
    } catch (err) {
        if (axios.isAxiosError(err)) {
            success.value = ''
            error.value = err.response?.data.error || err.message

            if (err.response?.status === 429) {
                error.value = 'Too many requests, please try again later.'
            }
        }
    }
}
</script>