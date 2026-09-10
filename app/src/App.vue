<template>
  <router-view />
  <StepUpAuth />
</template>

<script setup lang="ts">
import StepUpAuth from './components/StepUpAuth.vue'

// Evaluate the initial theme
if (localStorage.theme === 'dark' || (!('theme' in localStorage) && window.matchMedia('(prefers-color-scheme: dark)').matches)) {
  document.documentElement.classList.add('dark')
} else {
  document.documentElement.classList.remove('dark')
}
</script>