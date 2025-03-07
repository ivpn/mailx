<template>
    <div class="page">
        <Header />
        <div class="container mx-auto max-w-screen-lg px-5 pb-5">
            <SubscriptionStatus />
            <router-view></router-view>
        </div>
    </div>
</template>

<script setup lang="ts">
import Header from './Header.vue'
import SubscriptionStatus from './AccountSubscriptionStatus.vue'
</script>