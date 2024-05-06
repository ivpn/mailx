<template>
    <AccountConfirm :dashboard="true" />
    <Aliases :dashboard="true" />
    <Recipients :dashboard="true" />
</template>

<script setup lang="ts">
import AccountConfirm from './AccountConfirm.vue'
import Aliases from './Aliases.vue'
import Recipients from './Recipients.vue'
</script>