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
            <a class="mt-3 py-2 px-3 inline-flex justify-center items-center gap-x-2 text-sm font-medium rounded-md border border-transparent bg-violet-600 text-white hover:bg-violet-700"
                href="#">
                Create Alias
            </a>
            <a class="mt-3 py-2 px-3 inline-flex justify-center items-center gap-x-2 text-sm font-medium rounded-md border border-transparent bg-violet-600 text-white hover:bg-violet-700"
                href="#">
                Add Recipient
            </a>
        </div>
    </div>
</template>

<script setup lang="ts">
import AccountConfirm from './AccountConfirm.vue'
</script>