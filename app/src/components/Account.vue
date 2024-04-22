<template>
    <div class="flex flex-col bg-white border border-gray-200 rounded-xl p-5 pb-4 my-8">
        <h1 class="text-lg font-bold text-gray-800 mb-4">Account</h1>
        <p class="text-sm text-gray-500 mb-3">
            {{ jwt().email }} user.name@example.net
        </p>
        <!-- <p class="text-sm text-gray-500 mb-3">
            <span class="inline-flex items-center py-1.5 px-2 rounded-md text-xs font-medium bg-teal-100 text-teal-800">Confirmed</span>
        </p> -->
        <p class="text-sm text-gray-500 mb-4">
            <span class="inline-flex items-center py-1.5 px-2 rounded-md text-xs font-medium bg-gray-100 text-gray-500">Not confirmed</span>
        </p>
        <div class="mb-4 max-w-xs">
            <label class="block text-gray-500 text-sm font-semibold mb-3" for="account-otp">
                Enter OTP code to confirm email
            </label>
            <input
                class="appearance-none outline-none border-2 rounded-md w-full py-3 px-4 text-gray-700 leading-tight focus:border-blue-600 mb-4"
                id="account-otp" type="email" autocomplete="account-otp">
            <!-- <p class="text-red-500 text-sm mb-3">Required field</p> -->
            <button
                class="bg-blue-600 hover:bg-blue-700 text-white font-bold py-2 px-3 rounded-md focus:outline-none focus:shadow-outline"
                type="submit">
                Send
            </button>
        </div>
    </div>
</template>

<script setup lang="ts">
import { jwt } from '../utils/jwt'
</script>