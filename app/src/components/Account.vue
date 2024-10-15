<template>
    <AccountVerify />
    <div class="flex flex-col p-5 my-8 bg-white dark:bg-neutral-800">
        <h1 class="text-3xl font-semibold text-gray-800 dark:text-gray-100 mb-6">Account Settings</h1>
        <AccountSubscription />
        <hr class="mb-8 dark:border-neutral-600">
        <AccountChangePassword />
        <hr class="mb-8 dark:border-neutral-600">
        <AccountTotp />
        <hr class="mb-8 dark:border-neutral-600">
        <AccountPasskeys />
        <hr class="mb-8 dark:border-neutral-600">
        <AccountDelete />
    </div>
</template>

<script setup lang="ts">
import AccountVerify from './AccountVerify.vue'
import AccountSubscription from './AccountSubscription.vue'
import AccountChangePassword from './AccountChangePassword.vue'
import AccountDelete from './AccountDelete.vue'
import AccountTotp from './AccountTotp.vue'
import AccountPasskeys from './AccountPasskeys.vue'
</script>