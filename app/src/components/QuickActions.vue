<template>
    <AccountConfirm />
    <div class="flex flex-col items-center p-4 text-center pt-20">
        <h3 class="text-lg font-bold text-gray-800">
            Quick actions
        </h3>
        <p class="my-2 text-gray-500">
            To get started, create an alias or add a recipient.
        </p>
        <div class="flex gap-4">
            <button
                class="mt-3 py-2 pl-2 pr-3 inline-flex justify-center items-center gap-x-2 text-sm font-medium rounded-md border border-transparent bg-violet-600 text-white hover:bg-violet-700">
                <svg class="flex-shrink-0 size-3.5" xmlns="http://www.w3.org/2000/svg" width="24" height="24"
                    viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="3" stroke-linecap="round"
                    stroke-linejoin="round">
                    <path d="M5 12h14"></path>
                    <path d="M12 5v14"></path>
                </svg>
                Add Alias
            </button>
            <button
                class="mt-3 py-2 pl-2 pr-3 inline-flex justify-center items-center gap-x-2 text-sm font-medium rounded-md border border-transparent bg-violet-600 text-white hover:bg-violet-700">
                <svg class="flex-shrink-0 size-3.5" xmlns="http://www.w3.org/2000/svg" width="24" height="24"
                    viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="3" stroke-linecap="round"
                    stroke-linejoin="round">
                    <path d="M5 12h14"></path>
                    <path d="M12 5v14"></path>
                </svg>
                Add Recipient
            </button>
        </div>
    </div>
</template>

<script setup lang="ts">
import AccountConfirm from './AccountConfirm.vue'
</script>