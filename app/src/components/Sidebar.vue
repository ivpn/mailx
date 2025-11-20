<template>
    <header class="bg-secondary flex flex-col justify-between h-full">
        <nav>
            <router-link to="/" class="p-0">
                <h1 class="pl-6 pr-5 m-0 text-accent head flex items-center justify-between">
                    <span class="logo"></span>
                </h1>
            </router-link>
            <div class="flex flex-col items-center">
                <router-link v-bind:class="{ 'active': route == '/' }" to="/">
                    <i class="icon at icon-primary"></i>
                    Aliases
                </router-link>
                <router-link v-bind:class="{ 'active': route == '/wildcard' }" to="/wildcard">
                    <i class="icon scan icon-primary"></i>
                    Wildcard
                </router-link>
                <router-link v-bind:class="{ 'active': route == '/recipients' }" to="/recipients">
                    <i class="icon inbox icon-primary"></i>
                    Recipients
                </router-link>
                <router-link v-bind:class="{ 'active': route == '/stats' }" to="/stats">
                    <i class="icon chart icon-primary"></i>
                    Stats
                </router-link>
                <router-link v-bind:class="{ 'active': route == '/logs' }" to="/logs">
                    <i class="icon alert icon-primary"></i>
                    Logs
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
        <div>
            <nav>
                <div class="flex items-center py-5 pb-3 pr-5">
                    <a @click.stop="logout">
                        <i class="icon logout icon-primary"></i>
                        Log out
                    </a>
                    <ThemeSwitch />
                </div>
            </nav>
            <p class="px-5 pl-6 text-sm">
                Support:
                <a href="mailto:mailx@ivpn.net">Email</a> /
                <a href="/faq">FAQ</a>
            </p>
        </div>
    </header>
</template>

<script setup lang="ts">
import { ref, onMounted, watch } from 'vue'
import { useRoute } from 'vue-router'
import { userApi } from '../api/user.ts'
import events from '../events.ts'
import ThemeSwitch from './ThemeSwitch.vue'

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