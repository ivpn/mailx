<template>
    <div class="flex flex-col bg-white shadow-sm rounded-xl p-5 pb-4 my-8">
        <h1 class="text-2xl font-bold text-gray-800 mb-5">Messages</h1>
        <h2 class="font-semibold text-gray-800 mb-5">Last 7 days</h2>
        <div id="chart" class="mb-5"></div>
        <h2 class="font-semibold text-gray-800 mb-5">Last 90 days</h2>
        <div class="grid grid-cols-2 md:grid-cols-4 gap-4 text-center mb-8">
            <div class="p-4 border-r border-gray-200">
                <p class="text-4xl font-bold text-gray-800 mb-2">{{ stats.forwards }}</p>
                <p class="text-gray-500">Forwards</p>
            </div>
            <div class="p-4 border-r border-white md:border-gray-200">
                <p class="text-4xl font-bold text-gray-800 mb-2">{{ stats.blocks }}</p>
                <p class="text-gray-500">Blocks</p>
            </div>
            <div class="p-4 border-r border-gray-200">
                <p class="text-4xl font-bold text-gray-800 mb-2">{{ stats.replies }}</p>
                <p class="text-gray-500">Replies</p>
            </div>
            <div class="p-4">
                <p class="text-4xl font-bold text-gray-800 mb-2">{{ stats.sends }}</p>
                <p class="text-gray-500">Sends</p>
            </div>
        </div>
        <p v-if="error" class="text-red-600 text-sm mb-3">Error: {{ error }}</p>
    </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { userApi } from '../api/user.ts'
import axios from 'axios'
import ApexCharts from 'apexcharts'

const stats = ref({
    forwards: 0,
    blocks: 0,
    replies: 0,
    sends: 0,
    bandwidth: 0,
    aliases: 0,
    messages: {
        created_at: '',
        type: 0,
    }
})
const error = ref('')

const getStats = async () => {
    try {
        const response = await userApi.stats()
        stats.value = response.data
        error.value = ''
    } catch (err) {
        if (axios.isAxiosError(err)) {
            error.value = err.response?.data.error || err.message
        }
    }
}

const initChart = () => {
    const options = {
        chart: {
            type: 'bar',
            height: 350,
            toolbar: {
                show: false,
            },
        },
        dataLabels: {
            enabled: false
        },
        series: [
            {
                name: 'Forwards',
                data: [30, 40, 35, 50, 49, 60, 70],
            },
            {
                name: 'Blocks',
                data: [23, 12, 54, 61, 32, 56, 81],
            },
            {
                name: 'Replies',
                data: [45, 23, 56, 32, 34, 52, 41],
            },
            {
                name: 'Sends',
                data: [23, 12, 54, 61, 32, 56, 81],
            },
        ],
        xaxis: {
            categories: ['Mon', 'Tue', 'Wed', 'Thu', 'Fri', 'Sat', 'Sun'],
        },
    }

    const chart = new ApexCharts(document.querySelector('#chart'), options)
    chart.render()
}

onMounted(() => {
    getStats()
    initChart()
})
</script>