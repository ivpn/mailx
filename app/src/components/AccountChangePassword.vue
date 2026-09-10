<template>
    <div v-if="res.id" class="mb-5">
        <h2>{{ res.has_password ? 'Change Password' : 'Set Password' }}</h2>
        <div v-if="res.has_password" class="mb-4 max-w-xs">
            <label for="old-password">
                Old password:
            </label>
            <input
                v-model="oldPassword"
                v-bind:class="{ 'error': passwordError }"
                id="old-password"
                type="password"
            >
        </div>
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
                {{ res.has_password ? 'Change Password' : 'Set Password' }}
            </button>
        </div>
        <p v-if="passwordError" class="error">Error: {{ passwordError }}</p>
        <p v-if="error" class="error">Error: {{ error }}</p>
        <p v-if="success" class="success">{{ success }}</p>
    </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { userApi } from '../api/user.ts'
import axios from 'axios'

const res = ref({
    id: '',
    has_password: false
})
const oldPassword = ref('')
const password = ref('')
const passwordConfirm = ref('')
const passwordError = ref('')
const error = ref('')
const success = ref('')

const getUser = async () => {
    try {
        const response = await userApi.get()
        res.value = response.data
    } catch (err) {
        if (axios.isAxiosError(err)) {
            error.value = err.message
        }
    }
}

const validatePassword = () => {
    success.value = ''
    passwordError.value = ''

    if (res.value.has_password && !oldPassword.value) {
        passwordError.value = 'Please fill required fields'
    }

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
        old_password: oldPassword.value,
        password: password.value
    }

    try {
        const response = await userApi.changePassword(req)
        success.value = response.data.message
        error.value = ''
        oldPassword.value = ''
        password.value = ''
        passwordConfirm.value = ''
        getUser()
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

onMounted(() => {
    getUser()
})
</script>
