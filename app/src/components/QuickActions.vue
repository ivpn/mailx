<template>
    <AccountConfirm :dashboard="true" />
    <Aliases :dashboard="false" />
</template>

<script setup lang="ts">
import AccountConfirm from './AccountVerify.vue'
import Aliases from './Aliases.vue'
</script>