<template>
    <AccountConfirm />
    <div class="flex flex-col p-5 pb-4 my-8 bg-white">
        <h1 class="text-3xl font-semibold text-gray-800 mb-6">Account Settings</h1>
        <Subscription />
        <DeleteAccount />
    </div>
</template>

<script setup lang="ts">
import AccountConfirm from './AccountVerify.vue'
import Subscription from './Subscription.vue';
import DeleteAccount from './AccountDelete.vue'
</script>