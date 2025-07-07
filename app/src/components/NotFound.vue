<template>
    <div class="page center">
        <div></div>
        <div class="card-tertiary center text-center">
            <article>
                <h1 class="text-5xl">404</h1>
                <p>Page not found</p>
                <router-link v-if="isLoggedIn()" to="/" tag="button" class="cta full mb-5">
                    Go to Aliases
                </router-link>
                <router-link v-if="!isLoggedIn()" to="/login" tag="button" class="cta full mb-5">
                    Go to Login
                </router-link>
            </article>
        </div>
        <div v-if="isLoggedIn()"></div>
        <Footer v-if="!isLoggedIn()" />
    </div>
</template>

<script setup lang="ts">
import Footer from './Footer.vue'

const isLoggedIn = (): boolean => {
    const email = localStorage.getItem('email')
    return email !== null && email.trim() !== ''
}
</script>