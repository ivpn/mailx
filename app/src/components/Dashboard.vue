<template>
    <div class="page">
        <div class="w-[223px] flex-shrink-0 fixed h-full">
            <Sidebar />
        </div>
        <div class="ml-[223px] flex-1">
            <SubscriptionStatus />
            <router-view></router-view>
        </div>
    </div>
</template>

<script setup lang="ts">
import Sidebar from './Sidebar.vue'
import SubscriptionStatus from './AccountSubscriptionStatus.vue'
</script>