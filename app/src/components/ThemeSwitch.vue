<template>
    <div>
        <button
            v-if="isLight()"
            :key="rowKey"
            @click="toggleTheme"
            class="hs-dark-mode inline-flex items-center gap-x-2 py-2 px-2 rounded-full text-sm cancel focus:outline-none"
            data-hs-theme-click-value="dark">
            <svg class="shrink-0 size-4" xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24"
                fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                <path d="M12 3a6 6 0 0 0 9 9 9 9 0 1 1-9-9Z"></path>
            </svg>
        </button>
        <button
            v-if="isDark()"
            :key="rowKey"
            @click="toggleTheme"
            class="hs-dark-mode inline-flex items-center gap-x-2 py-2 px-2 rounded-full text-sm cancel focus:outline-none"
            data-hs-theme-click-value="light">
            <svg class="shrink-0 size-4" xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24"
                fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                <circle cx="12" cy="12" r="4"></circle>
                <path d="M12 2v2"></path>
                <path d="M12 20v2"></path>
                <path d="m4.93 4.93 1.41 1.41"></path>
                <path d="m17.66 17.66 1.41 1.41"></path>
                <path d="M2 12h2"></path>
                <path d="M20 12h2"></path>
                <path d="m6.34 17.66-1.41 1.41"></path>
                <path d="m19.07 4.93-1.41 1.41"></path>
            </svg>
        </button>
    </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'

const rowKey = ref(0)

const toggleTheme = () => {
    if (localStorage.theme === 'dark' || (!('theme' in localStorage) && window.matchMedia('(prefers-color-scheme: dark)').matches)) {
        document.documentElement.classList.remove('dark')
        document.documentElement.classList.add('light')
        localStorage.theme = 'light'
    } else {
        document.documentElement.classList.remove('light')
        document.documentElement.classList.add('dark')
        localStorage.theme = 'dark'
    }
    renderRow()
}

const isDark = () => {
    return localStorage.theme === 'dark' || (!('theme' in localStorage) && window.matchMedia('(prefers-color-scheme: dark)').matches)
}

const isLight = () => {
    return localStorage.theme === 'light' || (!('theme' in localStorage) && window.matchMedia('(prefers-color-scheme: light)').matches)
}

const renderRow = () => {
    rowKey.value++
}

</script>