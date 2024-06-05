<template>
    <nav>
        <div class="flex flex-row items-center gap-8">
            <a v-bind:class="{ 'text-bluish-500 border-bluish-500 border-b-2': route == '/', 'text-gray-500 border-b-2 border-white': route != '/' }"
                class="py-5 font-medium hover:text-bluish-500" href="/">Dashboard</a>
            <a v-bind:class="{ 'text-bluish-500 border-bluish-500 border-b-2': route == '/aliases', 'text-gray-500 border-b-2 border-white': route != '/aliases' }"
                class="py-5 font-medium hover:text-bluish-500" href="/aliases">Aliases</a>
            <a v-bind:class="{ 'text-bluish-500 border-bluish-500 border-b-2': route == '/recipients', 'text-gray-500 border-b-2 border-white': route != '/recipients' }"
                class="py-5 font-medium hover:text-bluish-500" href="/recipients">Recipients</a>
            <a v-bind:class="{ 'text-bluish-500 border-bluish-500 border-b-2': route == '/stats', 'text-gray-500 border-b-2 border-white': route != '/stats' }"
                class="py-5 font-medium hover:text-bluish-500" href="/stats">Stats</a>
            <a v-bind:class="{ 'text-bluish-500 border-bluish-500 border-b-2': route == '/settings', 'text-gray-500 border-b-2 border-white': route != '/settings' }"
                class="py-5 font-medium hover:text-bluish-500" href="/settings">Settings</a>
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