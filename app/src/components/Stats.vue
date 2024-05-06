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

interface Message {
    created_at: string
    type: number
}

interface CountData {
    name: string
    data: number[]
}

const stats = ref({
    forwards: 0,
    blocks: 0,
    replies: 0,
    sends: 0,
    bandwidth: 0,
    aliases: 0,
    messages: [],
})
const error = ref('')

const getStats = async () => {
    try {
        const response = await userApi.stats()
        stats.value = response.data
        error.value = ''
        initChart()
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
            enabled: false,
        },
        series: getLast7DaysCounts(stats.value.messages),
        xaxis: {
            categories: getLast7Days(),
        },
        yaxis: {
            forceNiceScale: true,
        },
    }

    const chart = new ApexCharts(document.querySelector('#chart'), options)
    chart.render()
}

function getLast7Days(): string[] {
    const result: string[] = []

    for (let i = 6; i >= 0; i--) {
        const today = new Date()
        today.setDate(today.getDate() - i)
        result.push((today).toLocaleDateString('en-US', { weekday: 'short' }))
    }

    return result
}

function getLast7DaysCounts(messages: Message[]): CountData[] {
    const typeNames = ['Forwards', 'Blocks', 'Replies', 'Sends']

    const days: { [key: string]: number[] } = {
        Forwards: Array(7).fill(0),
        Blocks: Array(7).fill(0),
        Replies: Array(7).fill(0),
        Sends: Array(7).fill(0),
    }

    const now = new Date()
    const sevenDaysAgo = new Date(now)
    sevenDaysAgo.setDate(sevenDaysAgo.getDate() - 7)

    messages.forEach((msg) => {
        const messageDate = new Date(msg.created_at)
        if (messageDate >= sevenDaysAgo) {
            const diffTime = Math.abs(now.getTime() - messageDate.getTime())
            const dayIndex = Math.floor(diffTime / (1000 * 60 * 60 * 24))
            const typeIndex = msg.type

            if (typeIndex >= 0 && typeIndex < 4 && dayIndex >= 0 && dayIndex < 7) {
                days[typeNames[typeIndex]][5 - dayIndex]++
            }
        }
    })

    const result: CountData[] = []
    for (let i = 0; i < 4; i++) {
        result.push({
            name: typeNames[i],
            data: days[typeNames[i]],
        })
    }

    return result
}

onMounted(() => {
    getStats()
})
</script>