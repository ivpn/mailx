<template>
    <div class="page page-bg">
        <div class="max-sm:hidden [@media(max-height:530px)]:hidden w-[224px] flex-shrink-0 fixed h-full">
            <Sidebar />
        </div>
        <div class="mb-[65px] sm:mb-0 sm:ml-[224px] [@media(max-height:530px)]:ml-0 [@media(max-height:530px)]:mb-[65px] grow overflow-x-hidden">
            <div class="mobile">
                <MobileHeader />
            </div>
            <router-view></router-view>
            <div class="mobile">
                <Tabbar />
            </div>
        </div>
    </div>
</template>

<script setup lang="ts">
import Sidebar from './Sidebar.vue'
import Tabbar from './Tabbar.vue';
import MobileHeader from './MobileHeader.vue';
</script>