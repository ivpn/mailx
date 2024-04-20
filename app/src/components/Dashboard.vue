<template>
    <div class="container mx-auto max-w-screen-lg px-5">
        <Header />
        <router-view></router-view>
    </div>
</template>

<script setup lang="ts">
import Header from './Header.vue';
</script>