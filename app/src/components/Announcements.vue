<template>
    <div class="container mx-auto max-w-screen-lg sm:p-10 p-5">
        <p class="py-5">
            <button @click="goBack" class="flex items-center gap-2">
                <i class="icon arrow-left-line icon-accent"></i>
                Back
            </button>
        </p>
        <h1>Announcements</h1>
        <hr>

        <div v-if="!loaded && !error">
            <p>Loading...</p>
        </div>

        <div v-if="error">
            <p class="text-red-500">{{ error }}</p>
        </div>

        <div v-if="loaded && list.length === 0">
            <p>No announcements.</p>
        </div>

        <template v-if="loaded">
            <div v-for="(item, index) in list" :key="index">
                <h2>{{ item.title }}</h2>
                <div v-html="item.body"></div>
                <p v-if="item.link">
                    <a :href="item.link" target="_blank" rel="noopener noreferrer">{{ item.link }}</a>
                </p>
                <hr v-if="index < list.length - 1">
            </div>
        </template>
    </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import axios from 'axios'
import { announcementsApi } from '../api/announcements.ts'

interface Announcement {
    title: string
    body: string
    link: string
}

const router = useRouter()
const goBack = () => router.back()

const list = ref<Announcement[]>([])
const error = ref('')
const loaded = ref(false)

onMounted(() => {
    getList()
})

const getList = async () => {
    try {
        const response = await announcementsApi.getList()
        list.value = response.data
        loaded.value = true
        error.value = ''
    } catch (err) {
        if (axios.isAxiosError(err)) {
            error.value = err.response?.data.error || err.message

            if (err.response?.status === 429) {
                error.value = 'Too many requests, please try again later.'
            }
        }
    }
}
</script>
