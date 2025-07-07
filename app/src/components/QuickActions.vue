<template>
    <AccountVerify :dashboard="true" />
    <SubscriptionStatus :dashboard="true" />
    <Aliases />
</template>

<script setup lang="ts">
import AccountVerify from './AccountVerify.vue'
import SubscriptionStatus from './AccountSubscriptionStatus.vue'
import Aliases from './Aliases.vue'
</script>