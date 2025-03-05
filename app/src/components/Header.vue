<template>
    <header class="flex flex-col justify-between w-full bg-white dark:bg-neutral-800">
        <div class="container mx-auto max-w-screen-lg px-5">
            <div class="flex flex-row justify-between max-w-screen-lg">
                <HeaderMenu />
                <div class="hs-dropdown relative flex items-center [--placement:bottom-right] my-3">
                    <button id="hs-dropdown-default"
                        class="flex items-center hs-dropdown-toggle text-gray-500 dark:text-gray-400 pl-4 pr-3 hover:text-gray-800 dark:hover:text-gray-100">
                        {{ email }}
                        <svg class="ms-1 flex-shrink-0 size-5" xmlns="http://www.w3.org/2000/svg" width="24" height="24"
                            viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"
                            stroke-linecap="round" stroke-linejoin="round">
                            <path d="m6 9 6 6 6-6" />
                        </svg>
                    </button>
                    <div class="hs-dropdown-menu transition-[opacity,margin] duration hs-dropdown-open:opacity-100 opacity-0 hidden min-w-60 bg-white dark:dark:bg-neutral-800 border border-gray-200 dark:border-neutral-600 shadow-sm p-2 mt-2 after:h-4 after:absolute after:-bottom-4 after:start-0 after:w-full before:h-4 before:absolute before:-top-4 before:start-0 before:w-full"
                        aria-labelledby="hs-dropdown-default">
                        <router-link class="flex items-center gap-x-3.5 py-2 px-3 text-gray-500 dark:text-gray-400 hover:bg-gray-100 dark:hover:bg-neutral-700 focus:outline-none"
                            to="/account">
                            Account
                        </router-link>
                        <button @click.stop="logout"
                            class="flex items-center gap-x-3.5 py-2 px-3 text-red-600 hover:bg-gray-100 dark:hover:bg-neutral-700 focus:outline-none w-full"
                            href="#">
                            Log out
                        </button>
                    </div>
                    <ThemeSwitch />
                </div>
            </div>
        </div>
    </header>
</template>

<script setup lang="ts">
import HeaderMenu from './HeaderMenu.vue'
import ThemeSwitch from './ThemeSwitch.vue'
import dropdown from '@preline/dropdown'
import { userApi } from '../api/user.ts'
import { onMounted, ref } from 'vue'
import events from '../events.ts'

const email = ref(localStorage.getItem('email'))

const logout = async () => {
    if (!confirm('Do you want to proceed with the logout?')) return

    try {
        await userApi.logout()
        userApi.clearSession()
    } catch { }
}

const onUpdateEmail = (event: any) => {
    email.value = event.email
}

onMounted(() => {
    dropdown.autoInit()
    events.on('user.update', onUpdateEmail)
})
</script>