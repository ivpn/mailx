<template>
    <header class="bg-secondary flex flex-col justify-between h-full">
        <nav>
            <router-link to="/account" class="p-0">
                <h1 class="pl-6 pr-5 m-0 text-accent head flex items-center justify-between">
                    <span class="logo"></span>
                </h1>
            </router-link>
            <div class="flex flex-col items-center">
                <router-link v-bind:class="{ 'active': route == '/account' && !route.startsWith('/account/') }" to="/account">
                    <i class="icon at icon-primary"></i>
                    Aliases
                </router-link>
                <router-link v-bind:class="{ 'active': route == '/account/wildcard' }" to="/account/wildcard">
                    <i class="icon scan icon-primary"></i>
                    Wildcard
                </router-link>
                <router-link v-bind:class="{ 'active': route == '/account/recipients' }" to="/account/recipients">
                    <i class="icon inbox icon-primary"></i>
                    Recipients
                </router-link>
                <router-link v-bind:class="{ 'active': route == '/account/domains' }" to="/account/domains">
                    <i class="icon global icon-primary"></i>
                    Domains
                </router-link>
                <router-link v-bind:class="{ 'active': route == '/account/stats' }" to="/account/stats">
                    <i class="icon chart icon-primary"></i>
                    Stats
                </router-link>
                <router-link v-bind:class="{ 'active': route == '/account/diagnostics' }" to="/account/diagnostics">
                    <i class="icon alert icon-primary"></i>
                    Diagnostics
                </router-link>
                <router-link v-bind:class="{ 'active': route == '/account/settings' }" to="/account/settings">
                    <i class="icon settings icon-primary"></i>
                    Settings
                </router-link>
                <router-link v-bind:class="{ 'active': route == '/account/profile' }" to="/account/profile">
                    <i class="icon user icon-primary"></i>
                    Account
                </router-link>
            </div>
        </nav>
        <div>
            <nav class="footer">
                <div class="[@media(max-height:715px)]:hidden">
                    <router-link to="/announcements">
                        <i class="icon logout icon-primary"></i>
                        Announcements
                    </router-link>
                    <router-link to="/faq">
                        <i class="icon logout icon-primary"></i>
                        FAQ
                    </router-link>
                    <router-link to="/privacy">
                        <i class="icon logout icon-primary"></i>
                        Privacy
                    </router-link>
                    <a href="mailto:mailx@ivpn.net">
                        <i class="icon logout icon-primary"></i>
                        Support
                    </a>
                </div>
                <div class="flex items-center py-5 pb-3 pr-5">
                    <a @click.stop="logout">
                        <i class="icon logout icon-primary"></i>
                        Log out
                    </a>
                    <ThemeSwitch />
                </div>
            </nav>
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