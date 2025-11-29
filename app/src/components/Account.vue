<template>
    <div class="card-container">
        <header class="head">
            <h2>Account Settings</h2>
        </header>
        <div class="card-primary">
            <AccountVerify />
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
            <AccountAccessKeys />
            <hr>
            <AccountAliasExport />
            <hr>
            <AccountDelete />
            <hr>
            <h2>Support</h2>
            <p>
                Support:
                <a href="mailto:mailx@ivpn.net">Email</a> /
                <a href="/faq">FAQ</a>
            </p>
        </div>
    </div>
</template>

<script setup lang="ts">
import AccountVerify from './AccountVerify.vue'
import AccountSubscription from './AccountSubscription.vue'
import AccountChangeEmail from './AccountChangeEmail.vue'
import AccountChangePassword from './AccountChangePassword.vue'
import AccountTotp from './AccountTotp.vue'
import AccountPasskeys from './AccountPasskeys.vue'
import AccountAccessKeys from './AccountAccessKeys.vue'
import AccountAliasExport from './AccountAliasExport.vue'
import AccountDelete from './AccountDelete.vue'
</script>