<template>
    <AccountConfirm :dashboard="true" />
    <Recipients :dashboard="true" />
    <Aliases :dashboard="true" />
</template>

<script setup lang="ts">
import AccountConfirm from './AccountConfirm.vue'
import Aliases from './Aliases.vue'
import Recipients from './Recipients.vue'
</script>