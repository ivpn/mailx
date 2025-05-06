<template>
    <div class="mb-5">
        <h2>Change Password</h2>
        <div class="mb-4 max-w-xs">
            <label for="new-password">
                New password:
            </label>
            <input
                v-model="password"
                v-bind:class="{ 'error': passwordError }"
                id="new-password"
                type="password"
            >
        </div>
        <div class="mb-4 max-w-xs">
            <label for="new-password-confirm">
                Confirm new password:
            </label>
            <input
                v-model="passwordConfirm"
                v-bind:class="{ 'error': passwordError }"
                id="new-password-confirm"
                type="password"
            >
        </div>
        <p class="text-sm">
            Must be 12+ characters and contain uppercase, lowercase, number, and special character (e.g. -_+=~!@#$%^&*(),;.?":{}|<>)
        </p>
        <div class="mb-3 max-w-xs">
            <button
                @click="changePassword"
                class="cta">
                Change Password
            </button>
        </div>
        <p v-if="passwordError" class="error">Error: {{ passwordError }}</p>
        <p v-if="error" class="error">Error: {{ error }}</p>
        <p v-if="success" class="success">{{ success }}</p>
    </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { userApi } from '../api/user.ts'
import axios from 'axios'

const password = ref('')
const passwordConfirm = ref('')
const passwordError = ref('')
const error = ref('')
const success = ref('')

const validatePassword = () => {
    success.value = ''
    passwordError.value = ''

    if (!password.value || !passwordConfirm.value) {
        passwordError.value = 'Please fill required fields'
    }

    if (password.value !== passwordConfirm.value) {
        passwordError.value = 'Passwords do not match'
    }

    return !passwordError.value
}

const changePassword = async () => {
    if (!validatePassword()) return

    const req = {
        password: password.value
    }

    try {
        const res = await userApi.changePassword(req)
        success.value = res.data.message
        error.value = ''
        password.value = ''
        passwordConfirm.value = ''
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