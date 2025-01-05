<template>
    <nav>
        <div class="flex flex-row items-center gap-8">
            <router-link v-bind:class="{ 'text-gray-800 dark:text-gray-100 border-bluish-500 border-b-2': route == '/', 'text-gray-500 dark:text-gray-400 border-b-2 border-white dark:border-neutral-800': route != '/' }"
                class="py-6 font-medium hover:text-gray-800 dark:hover:text-gray-100" to="/">Dashboard</router-link>
            <router-link v-bind:class="{ 'text-gray-800 dark:text-gray-100 border-bluish-500 border-b-2': route == '/aliases', 'text-gray-500 dark:text-gray-400 border-b-2 border-white dark:border-neutral-800': route != '/aliases' }"
                class="py-6 font-medium hover:text-gray-800 dark:hover:text-gray-100" to="/aliases">Aliases</router-link>
            <router-link v-bind:class="{ 'text-gray-800 dark:text-gray-100 border-bluish-500 border-b-2': route == '/recipients', 'text-gray-500 dark:text-gray-400 border-b-2 border-white dark:border-neutral-800': route != '/recipients' }"
                class="py-6 font-medium hover:text-gray-800 dark:hover:text-gray-100" to="/recipients">Recipients</router-link>
            <router-link v-bind:class="{ 'text-gray-800 dark:text-gray-100 border-bluish-500 border-b-2': route == '/stats', 'text-gray-500 dark:text-gray-400 border-b-2 border-white dark:border-neutral-800': route != '/stats' }"
                class="py-6 font-medium hover:text-gray-800 dark:hover:text-gray-100" to="/stats">Stats</router-link>
            <router-link v-bind:class="{ 'text-gray-800 dark:text-gray-100 border-bluish-500 border-b-2': route == '/settings', 'text-gray-500 dark:text-gray-400 border-b-2 border-white dark:border-neutral-800': route != '/settings' }"
                class="py-6 font-medium hover:text-gray-800 dark:hover:text-gray-100" to="/settings">Settings</router-link>
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