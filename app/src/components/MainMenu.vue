<template>
    <nav>
        <div class="flex flex-col items-center">
            <router-link v-bind:class="{ 'active': route == '/' }" to="/">Aliases</router-link>
            <router-link v-bind:class="{ 'active': route == '/recipients' }" to="/recipients">Recipients</router-link>
            <router-link v-bind:class="{ 'active': route == '/stats' }" to="/stats">Stats</router-link>
            <router-link v-bind:class="{ 'active': route == '/settings' }" to="/settings">Settings</router-link>
            <router-link v-bind:class="{ 'active': route == '/account' }" to="/account">Account</router-link>
        </div>
    </nav>
</template>

<script setup lang="ts">
import { ref, watch } from 'vue'
import { useRoute } from 'vue-router'

const route = ref('/')
const currentRoute = useRoute()

watch(currentRoute, (newRoute) => {
    route.value = newRoute.path
}, { immediate: true })
</script>