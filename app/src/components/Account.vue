<template>
    <AccountVerify />
    <div class="flex flex-col p-5 my-8 bg-white">
        <h1 class="text-3xl font-semibold text-gray-800 mb-6">Account Settings</h1>
        <AccountSubscription />
        <hr class="mb-8">
        <AccountChangePassword />
        <hr class="mb-8">
        <AccountTotp />
        <hr class="mb-8">
        <AccountDelete />
    </div>
</template>

<script setup lang="ts">
import AccountVerify from './AccountVerify.vue'
import AccountSubscription from './AccountSubscription.vue';
import AccountChangePassword from './AccountChangePassword.vue';
import AccountDelete from './AccountDelete.vue'
import AccountTotp from './AccountTotp.vue'
</script>