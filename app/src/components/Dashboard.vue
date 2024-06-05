<template>
    <div class="min-h-screen">
        <Header />
        <div class="container mx-auto max-w-screen-lg px-5 pb-5">
            <router-view></router-view>
        </div>
    </div>
</template>

<script setup lang="ts">
import Header from './Header.vue';
</script>