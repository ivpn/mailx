<template>
    <div class="flex flex-col bg-white border border-gray-200 rounded-xl p-5 pb-4 my-8">
        <h1 class="text-lg font-bold text-gray-800 mb-4">Confirm Your Email</h1>
        <p class="text-sm text-gray-500 mb-3">
            We have sent a 6-digit OTP code to your email address. Please enter the code below to confirm your email. Accounts with unconfirmed email address may be deleted after 7 days.
        </p>
        <div class="mb-4 max-w-xs">
            <label class="block text-gray-500 text-sm font-semibold mb-3" for="account-otp">
                6-digit OTP code:
            </label>
            <input
                class="appearance-none outline-none border-2 rounded-md w-full py-3 px-4 text-gray-700 leading-tight focus:border-blue-600 mb-4"
                id="account-otp" type="email" autocomplete="account-otp">
            <!-- <p class="text-red-500 text-sm mb-3">Required field</p> -->
            <div class="flex flex-row gap-4">
                <button
                    class="bg-blue-600 hover:bg-blue-700 text-white font-medium text-sm py-2 px-3 rounded-md focus:outline-none focus:shadow-outline"
                    type="submit">
                    Confirm
                </button>
                <button
                    class="text-gray-500 bg-gray-100 hover:bg-gray-200 font-medium text-sm py-2 px-3 rounded-md focus:outline-none focus:shadow-outline"
                    type="submit">
                    Resend OTP
                </button>
            </div>
        </div>
    </div>
</template>

<script setup lang="ts">
</script>