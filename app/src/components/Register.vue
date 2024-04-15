<template>
    <div class="flex flex-col items-center justify-center min-h-screen bg-gray-100">
        <h1 class="text-2xl font-semibold mb-2">Sign up</h1>
        <p class="text-slate-500 text-sm mb-8">Have an account? <a class="text-blue-500 hover:text-blue-700" href="#">log in</a></p>
        <form class="w-full max-w-sm bg-white rounded px-8 pt-6 pb-8 mb-4" @submit.prevent="register">
            <div class="mb-4">
                <label class="block text-gray-700 text-sm font-bold mb-2" for="email">
                    Email
                </label>
                <input v-model="email"
                    class="appearance-none border rounded w-full py-3 px-4 text-gray-700 leading-tight focus:border-blue-500"
                    id="email" type="email" autocomplete="email">
            </div>
            <div class="mb-4">
                <label class="block text-gray-700 text-sm font-bold mb-2" for="password">
                    Password
                </label>
                <input v-model="password"
                    class="appearance-none border rounded w-full py-3 px-4 text-gray-700 mb-3 leading-tight focus:border-blue-500"
                    id="password" type="password" autocomplete="current-password">
            </div>
            <div class="flex items-center justify-between">
                <button
                    class="bg-blue-500 hover:bg-blue-700 text-white font-bold py-3 px-4 rounded focus:outline-none focus:shadow-outline"
                    type="submit">
                    Sign Up
                </button>
            </div>
        </form>
    </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { userApi } from '../api/user.ts'

const email = ref('')
const password = ref('')

const register = async () => {
    const data = {
        email: email.value,
        password: password.value
    }

    try {
        const response = await userApi.register(data)
        console.log(response)
    } catch (error) {
        console.error(error)
    }
}
</script>
