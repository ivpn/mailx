<template>
    <header class="bg-secondary fixed bottom-0 left-0 right-0 z-10">
        <nav class="mobile">
            <div class="flex flex-row items-center">
                <router-link v-bind:class="{ 'active': route == '/account' && !route.startsWith('/account/') }" to="/account" title="Aliases">
                    <i class="icon at icon-primary"></i>
                </router-link>
                <router-link v-bind:class="{ 'active': route == '/account/wildcard' }" to="/account/wildcard" title="Wildcard">
                    <i class="icon scan icon-primary"></i>
                </router-link>
                <router-link v-bind:class="{ 'active': route == '/account/recipients' }" to="/account/recipients" title="Recipients">
                    <i class="icon inbox icon-primary"></i>
                </router-link>
                <router-link v-bind:class="{ 'active': route == '/account/stats' }" to="/account/stats" title="Stats">
                    <i class="icon chart icon-primary"></i>
                </router-link>
                <router-link v-bind:class="{ 'active': route == '/account/diagnostics' }" to="/account/diagnostics" title="Diagnostics">
                    <i class="icon alert icon-primary"></i>
                </router-link>
                <router-link v-bind:class="{ 'active': route == '/account/settings' }" to="/account/settings" title="Settings">
                    <i class="icon settings icon-primary"></i>
                </router-link>
                <router-link v-bind:class="{ 'active': route == '/account/profile' }" to="/account/profile" title="Account">
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