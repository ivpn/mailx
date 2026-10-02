<template>
    <header class="bg-secondary fixed bottom-0 left-0 right-0 z-10">
        <nav class="mobile">
            <div class="flex flex-row items-center">
                <router-link v-bind:class="{ 'active': route == '/account' && !route.startsWith('/account/') }" to="/account" title="Aliases">
                    <i class="icon at icon-primary"></i>
                </router-link>
                <router-link v-bind:class="{ 'active': route == '/account/recipients' }" to="/account/recipients" title="Recipients">
                    <i class="icon inbox icon-primary"></i>
                </router-link>
                <router-link v-bind:class="{ 'active': route == '/account/blocklist' }" to="/account/blocklist" title="Blocklist">
                    <i class="icon block icon-primary"></i>
                </router-link>
                <router-link v-bind:class="{ 'active': route == '/account/rules' }" to="/account/rules" title="Rules">
                    <i class="icon filter icon-primary"></i>
                </router-link>
                <router-link v-bind:class="{ 'active': route == '/account/settings' }" to="/account/settings" title="Settings">
                    <i class="icon settings icon-primary"></i>
                </router-link>
                <router-link v-bind:class="{ 'active': route == '/account/profile' }" to="/account/profile" title="Account">
                    <i class="icon user icon-primary"></i>
                </router-link>
                <a v-bind:class="{ 'active': moreActive }" @click.stop="openMenu" title="More" aria-haspopup="dialog" :aria-expanded="menuOpen">
                    <i class="icon menu icon-primary"></i>
                </a>
            </div>
        </nav>
    </header>
    <Teleport to="body">
        <div v-if="menuOpen" class="mobile mobile-menu bg-secondary fixed inset-0 z-50 overflow-y-auto overscroll-contain" role="dialog" aria-modal="true" aria-label="Menu">
            <div class="flex items-center justify-between py-5 pr-3">
                <router-link to="/" class="p-0" @click="closeMenu">
                    <h1 class="px-5 m-0 text-accent head flex items-center justify-between">
                        <span class="logo"></span>
                    </h1>
                </router-link>
                <button @click="closeMenu" title="Close" class="w-10 h-10 flex items-center justify-center focus:outline-none">
                    <i class="icon close icon-primary"></i>
                </button>
            </div>
            <nav @click="closeMenu">
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
                <router-link v-bind:class="{ 'active': route == '/account/blocklist' }" to="/account/blocklist">
                    <i class="icon block icon-primary"></i>
                    Blocklist
                </router-link>
                <router-link v-bind:class="{ 'active': route == '/account/rules' }" to="/account/rules">
                    <i class="icon filter icon-primary"></i>
                    Rules
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
            </nav>
            <nav class="footer border-t border-secondary mt-3 pt-3 pb-6">
                <div @click="closeMenu">
                    <router-link to="/announcements">
                        <i class="icon megaphone icon-primary"></i>
                        Announcements
                    </router-link>
                    <router-link to="/faq">
                        <i class="icon question icon-primary"></i>
                        FAQ
                    </router-link>
                    <router-link to="/privacy">
                        <i class="icon eye-off icon-primary"></i>
                        Privacy
                    </router-link>
                    <a href="mailto:mailx@ivpn.net">
                        <i class="icon mail icon-primary"></i>
                        Support
                    </a>
                </div>
                <div class="flex items-center pr-5">
                    <a @click.stop="logout">
                        <i class="icon logout icon-primary"></i>
                        Log out
                    </a>
                    <ThemeSwitch />
                </div>
            </nav>
        </div>
    </Teleport>
</template>

<script setup lang="ts">
import { ref, computed, watch, onUnmounted } from 'vue'
import { useRoute } from 'vue-router'
import { userApi } from '../api/user.ts'
import ThemeSwitch from './ThemeSwitch.vue'

const moreRoutes = [
    '/account/domains',
    '/account/stats',
    '/account/diagnostics',
]

const route = ref('/')
const currentRoute = useRoute()
const menuOpen = ref(false)

const moreActive = computed(() => menuOpen.value || moreRoutes.includes(route.value))

const openMenu = () => {
    menuOpen.value = true
}

const closeMenu = () => {
    menuOpen.value = false
}

const onKeydown = (event: KeyboardEvent) => {
    if (event.key === 'Escape') closeMenu()
}

const logout = async () => {
    if (!confirm('Do you want to proceed with the logout?')) return

    try {
        await userApi.logout()
        userApi.clearSession()
    } catch { }
}

watch(menuOpen, (open) => {
    document.body.style.overflow = open ? 'hidden' : ''
    if (open) {
        document.addEventListener('keydown', onKeydown)
    } else {
        document.removeEventListener('keydown', onKeydown)
    }
})

watch(currentRoute, (newRoute) => {
    route.value = newRoute.path
    closeMenu()
}, { immediate: true })

onUnmounted(() => {
    document.body.style.overflow = ''
    document.removeEventListener('keydown', onKeydown)
})
</script>
