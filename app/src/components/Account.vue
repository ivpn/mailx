<template>
    <AccountVerify />
    <div class="flex flex-col p-5 my-8 bg-white">
        <h1 class="text-3xl font-semibold text-gray-800 mb-6">Account Settings</h1>
        <AccountSubscription />
        <AccountChangePassword />
        <AccountDelete />
    </div>
</template>

<script setup lang="ts">
import AccountVerify from './AccountVerify.vue'
import AccountSubscription from './AccountSubscription.vue';
import AccountChangePassword from './AccountChangePassword.vue';
import AccountDelete from './AccountDelete.vue'
</script>