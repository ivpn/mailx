<template>
    <header class="card flex flex-col w-full p-0 m-0">
        <div class="container mx-auto max-w-screen-lg px-0">
            <div class="flex flex-col max-w-screen-lg">
                <MainMenu />
                <div class="hs-dropdown">
                    <button id="hs-dropdown-default">
                        {{ email }}
                        <svg class="ms-1 flex-shrink-0 size-5" xmlns="http://www.w3.org/2000/svg" width="24" height="24"
                            viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"
                            stroke-linecap="round" stroke-linejoin="round">
                            <path d="m6 9 6 6 6-6" />
                        </svg>
                    </button>
                    <div class="hs-dropdown-menu hs-dropdown-open:opacity-100 hidden" aria-labelledby="hs-dropdown-default">
                        <button @click="$router.push('/account')">
                            Account
                        </button>
                        <button @click.stop="logout" class="delete">
                            Log out
                        </button>
                    </div>
                    <!-- <ThemeSwitch /> -->
                </div>
            </div>
        </div>
    </header>
</template>

<script setup lang="ts">
import MainMenu from './MainMenu.vue'
// import ThemeSwitch from './ThemeSwitch.vue'
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