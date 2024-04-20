<template>
    <header class="flex flex-col justify-between w-full text-sm border-b-2 border-gray-100">
        <div class="flex flex-row justify-between pt-6 pb-2">
            <a class="inline-flex items-center gap-x-2 text-xl font-semibold text-black" href="#">
                <span class="bg-black w-10 h-10 text-white flex items-center justify-center rounded-md">B</span>
                Brand
            </a>
            <button type="button"
                class="flex items-center text-gray-500 bg-gray-100 px-3 rounded-md hover:text-blue-5000 font-medium">
                user.name@example.net
                <svg class="ms-1 flex-shrink-0 size-4" xmlns="http://www.w3.org/2000/svg" width="24" height="24"
                    viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round"
                    stroke-linejoin="round">
                    <path d="m6 9 6 6 6-6" />
                </svg>
            </button>
        </div>
        <nav>
            <div class="flex flex-row items-center gap-6">
                <a
                v-bind:class="{ 'text-blue-600 border-blue-600 border-b-2': route == '/', 'text-gray-500': route != '/'}"
                class="py-4 font-medium hover:text-blue-600" href="/">Dashboard</a>
                <a
                v-bind:class="{ 'text-blue-600 border-blue-600 border-b-2': route == '/aliases', 'text-gray-500': route != '/aliases'}"
                class="py-4 font-medium hover:text-blue-600" href="/aliases">Aliases</a>
                <a
                v-bind:class="{ 'text-blue-600 border-blue-600 border-b-2': route == '/recipients', 'text-gray-500': route != '/recipients'}"
                class="py-4 font-medium hover:text-blue-600" href="/recipients">Recipients</a>
                <a
                v-bind:class="{ 'text-blue-600 border-blue-600 border-b-2': route == '/stats', 'text-gray-500': route != '/stats'}"
                class="py-4 font-medium hover:text-blue-600" href="/stats">Stats</a>
                <a
                v-bind:class="{ 'text-blue-600 border-blue-600 border-b-2': route == '/settings', 'text-gray-500': route != '/settings'}"
                class="py-4 font-medium hover:text-blue-600"  href="/settings">Settings</a>
                <a
                v-bind:class="{ 'text-blue-600 border-blue-600 border-b-2': route == '/subscription', 'text-gray-500': route != '/subscription'}"
                class="py-4 font-medium hover:text-blue-600"  href="/subscription">Subscription</a>
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