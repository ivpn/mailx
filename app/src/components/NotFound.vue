<template>
    <div class="page center">
        <div></div>
        <div class="card-tertiary center text-center">
            <article>
                <h1 class="text-5xl">404</h1>
                <p>Page not found</p>
                <router-link to="/" tag="button" class="cta full mb-5">
                    Go to Dashboard
                </router-link>
            </article>
            
        </div>
        <Footer />
    </div>
</template>

<script setup lang="ts">
import Footer from './Footer.vue'
</script>