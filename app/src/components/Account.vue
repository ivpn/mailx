<template>
    <AccountVerify />
    <div class="flex flex-col p-5 my-8 bg-white dark:bg-neutral-800">
        <h1>Account Settings</h1>
        <AccountSubscription />
        <hr>
        <AccountChangeEmail />
        <hr>
        <AccountChangePassword />
        <hr>
        <AccountTotp />
        <hr>
        <AccountPasskeys />
        <hr>
        <AccountDelete />
    </div>
</template>

<script setup lang="ts">
import AccountVerify from './AccountVerify.vue'
import AccountSubscription from './AccountSubscription.vue'
import AccountChangeEmail from './AccountChangeEmail.vue'
import AccountChangePassword from './AccountChangePassword.vue'
import AccountDelete from './AccountDelete.vue'
import AccountTotp from './AccountTotp.vue'
import AccountPasskeys from './AccountPasskeys.vue'
</script>