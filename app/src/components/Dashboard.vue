<template>
    <div class="page">
        <div class="w-64 flex-shrink-0 fixed h-full">
            <Sidebar />
        </div>
        <div class="ml-64 flex-1">
            <SubscriptionStatus />
            <router-view></router-view>
        </div>
    </div>
</template>

<script setup lang="ts">
import Sidebar from './Sidebar.vue'
import SubscriptionStatus from './AccountSubscriptionStatus.vue'
</script>