<template>
    <div class="page page-bg">
        <div class="w-[224px] flex-shrink-0 fixed h-full">
            <Sidebar />
        </div>
        <div class="ml-[224px] grow">
            <SubscriptionStatus />
            <router-view></router-view>
        </div>
    </div>
</template>

<script setup lang="ts">
import Sidebar from './Sidebar.vue'
import SubscriptionStatus from './AccountSubscriptionStatus.vue'
</script>