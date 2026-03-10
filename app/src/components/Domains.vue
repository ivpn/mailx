<template>
    <div class="card-container">
        <header class="head">
            <h2>Domains</h2>
            <div class="flex items-center justify-between">
                <DomainCreate />
            </div>
        </header>
        <div v-if="!list.length && loaded" class="card-empty">
            <span class="bg-secondary rounded flex items-center justify-center p-2 mb-5">
                <i class="icon global icon-accent text-2xl"></i>
            </span>
            <h4 class="mb-6">You have no domains yet</h4>
            <DomainCreate />
        </div>
        <div v-bind:class="{ 'hidden': !list.length || !loaded }" class="card-primary">
            <div class="table-container">
                <table>
                    <thead class="desktop">
                        <tr>
                            <th>Created</th>
                            <th>Domain</th>
                            <th>Description</th>
                            <th>Default Recipient</th>
                            <th>Active</th>
                            <th>Actions</th>
                        </tr>
                    </thead>
                    <tbody>
                        <!-- <DomainRow v-for="domain in list" :domain="domain" :key="domain.id" /> -->
                    </tbody>
                </table>
                <p v-if="error" class="error">Error: {{ error }}</p>
            </div>
        </div>
    </div>
</template>

<script setup lang="ts">
import { onMounted, ref } from 'vue'
import axios from 'axios'
import { domainApi } from '../api/domain.ts'
import DomainCreate from './DomainCreate.vue'
import events from '../events.ts'

const domain = {
    id: '',
    created_at: '',
    name: '',
}

const list = ref([] as typeof domain[])
const error = ref('')
const loaded = ref(false)
const rowKey = ref(0)

const getList = async () => {
    try {
        const response = await domainApi.getList()
        list.value = response.data
        error.value = ''
        renderRow()
    } catch (err) {
        if (axios.isAxiosError(err)) {
            error.value = err.response?.data.error || err.message
        }
    } finally {
        loaded.value = true
    }
}

const renderRow = () => {
    rowKey.value++
}

onMounted(() => {
    getList()
    events.on('domain.create', getList)
})
</script>