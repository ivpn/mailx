<template>
    <div class="flex flex-col bg-white shadow-sm rounded-xl p-5 pb-4 my-8">
        <h1 class="text-xl font-bold text-gray-800 mb-5">Aliases</h1>
        <div class="flex items-center justify-between mb-5">
            <button
                class="bg-violet-600 hover:bg-violet-700 text-white font-medium text-sm py-2 px-3 rounded-md focus:outline-none focus:shadow-outline"
                type="submit">
                Create Alias
            </button>
        </div>
        <div class="grid grid-cols-6 gap-2 mb-5">
            <div class="font-medium text-sm text-gray-500">Created</div>
            <div class="font-medium text-sm text-gray-500">Alias</div>
            <div class="font-medium text-sm text-gray-500">Recipients</div>
            <div class="font-medium text-sm text-gray-500" title="Forwards/Blocks/Replies/Sends">F/B/R/S</div>
            <div class="font-medium text-sm text-gray-500">Active</div>
        </div>
        <div v-if="list.length" class="mb-4">
            <AliasCard v-for="alias in list" :alias="alias" />
        </div>
        <p v-if="error" class="text-red-600 text-sm mb-4">{{ error }}</p>
    </div>
</template>

<script setup lang="ts">
import { onMounted, ref } from 'vue'
import axios from 'axios'
import { aliasApi } from '../api/alias'
import AliasCard from './AliasCard.vue'

const alias = ref({
    id: '',
    created: '',
    name: '',
    enabled: false,
    description: '',
    recipients: '',
    from_name: '',
    stats: {
        forwards: 0,
        blocks: 0,
        replies: 0,
        sends: 0
    }
})

const list = ref([] as typeof alias[])
const error = ref('')

const getList = async () => {
    try {
        const response = await aliasApi.getList()
        list.value = response.data
        error.value = ''
    } catch (err) {
        if (axios.isAxiosError(err)) {
            error.value = err.message
        }
    }
}

onMounted(() => {
    getList()
})

</script>