<template>
    <header class="bg-secondary fixed bottom-0 left-0 right-0 z-10">
        <nav class="mobile">
            <div class="flex flex-row items-center">
                <router-link v-bind:class="{ 'active': route == '/' }" to="/" title="Aliases">
                    <i class="icon at icon-primary"></i>
                </router-link>
                <router-link v-bind:class="{ 'active': route == '/wildcard' }" to="/wildcard" title="Wildcard">
                    <i class="icon scan icon-primary"></i>
                </router-link>
                <router-link v-bind:class="{ 'active': route == '/recipients' }" to="/recipients" title="Recipients">
                    <i class="icon inbox icon-primary"></i>
                </router-link>
                <router-link v-bind:class="{ 'active': route == '/stats' }" to="/stats" title="Stats">
                    <i class="icon chart icon-primary"></i>
                </router-link>
                <router-link v-bind:class="{ 'active': route == '/settings' }" to="/settings" title="Settings">
                    <i class="icon settings icon-primary"></i>
                </router-link>
                <router-link v-bind:class="{ 'active': route == '/account' }" to="/account" title="Account">
                    <i class="icon user icon-primary"></i>
                </router-link>
            </div>
        </nav>
    </header>
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