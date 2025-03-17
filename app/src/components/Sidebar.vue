<template>
    <header class="bg-secondary flex flex-col justify-between h-full">
        <nav>
            <h1 class="py-5 pl-8 m-0 mb-3 text-accent">MailX</h1>
            <div class="flex flex-col items-center">
                <router-link v-bind:class="{ 'active': route == '/' }" to="/">
                    <i class="icon at icon-primary"></i>
                    Aliases
                </router-link>
                <router-link v-bind:class="{ 'active': route == '/recipients' }" to="/recipients">
                    <i class="icon mailbox icon-primary"></i>
                    Recipients
                </router-link>
                <router-link v-bind:class="{ 'active': route == '/stats' }" to="/stats">
                    <i class="icon chart icon-primary"></i>
                    Stats
                </router-link>
                <router-link v-bind:class="{ 'active': route == '/settings' }" to="/settings">
                    <i class="icon settings icon-primary"></i>
                    Settings
                </router-link>
                <router-link v-bind:class="{ 'active': route == '/account' }" to="/account">
                    <i class="icon user icon-primary"></i>
                    Account
                </router-link>
            </div>
        </nav>
        <nav>
            <div class="flex flex-col items-center py-5">
                <a @click.stop="logout">
                    <i class="icon logout icon-primary"></i>
                    Log out
                </a>
            </div>
        </nav>
        <!-- <div class="hs-dropdown">
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
            <ThemeSwitch />
        </div> -->
    </header>
</template>

<script setup lang="ts">
// import ThemeSwitch from './ThemeSwitch.vue'
// import dropdown from '@preline/dropdown'
import { ref, onMounted, watch } from 'vue'
import { useRoute } from 'vue-router'
import { userApi } from '../api/user.ts'
import events from '../events.ts'

const route = ref('/')
const currentRoute = useRoute()
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
    events.on('user.update', onUpdateEmail)
})

watch(currentRoute, (newRoute) => {
    route.value = newRoute.path
}, { immediate: true })
</script>