<template>
    <div class="min-h-[550px] px-4 py-[50px]">
        <div class="w-full">
            <header class="fixed top-0 left-0 right-0 bg-secondary z-10 px-4 py-2">
                <div class="flex gap-3 items-center justify-between">
                    <div class="relative grow">
                        <form @submit.prevent="fetchAliases" autocomplete="off">
                            <input v-model="search" class="search w-full" type="text" placeholder="Search aliases..." id="input_search">
                        </form>
                        <button v-if="searchQuery" @click.prevent="clearSearch" class="absolute top-0 right-0 bottom-0 px-2 flex items-center justify-center">
                            <i class="icon close icon-tertiary text-base"></i>
                        </button>
                    </div>
                    <button class="cta sm text-nowrap" data-hs-overlay="#modal-create-alias">
                        New Alias
                    </button>
                </div>
            </header>
            <p v-if="isLoading" class="text-secondary py-4">Loading...</p>
            <p v-else-if="error" class="error">{{ error }}</p>
            <div v-else>
                <div v-for="(alias, index) in list" :key="alias.id" class="py-3 flex items-center gap-x-4" :class="{ 'border-t border-secondary': index > 0 }">
                    <div class="flex items-center">
                        <input @change="updateAlias(alias)" v-bind:checked="alias.enabled" type="checkbox" class="xs">
                    </div>
                    <div class="grow">
                        <p class="m-0 text-primary font-medium text-base">{{ alias.name }}</p>
                        <p class="m-0">{{ alias.description }}</p>
                    </div>
                    <div class="hs-tooltip">
                        <button class="hs-tooltip-toggle plain" @click="copyAlias(alias.name)">
                            <i class="icon icon-secondary copy text-xs"></i>
                            <span class="hs-tooltip-content hs-tooltip-shown:opacity-100 hs-tooltip-shown:visible opacity-0"
                                role="tooltip">
                                {{ copyText }}
                            </span>
                        </button>
                    </div>
                    <button @click="deleteAlias(alias.id)">
                        <i class="icon icon-secondary trash text-xs"></i>
                    </button>
                </div>
            </div>
        </div>
    </div>
    <AliasCreate :apiToken="props.apiToken" :defaults="props.defaults" />
</template>

<script lang="ts" setup>
import { ref, onMounted, onUpdated } from 'vue'
import { api } from '@/lib/api'
import { Defaults, Alias } from '@/lib/types'
import events from '@/lib/events'
import tooltip from '@preline/tooltip'
import AliasCreate from './AliasCreate.vue'

const props = defineProps<{
    apiToken: string
    defaults: Defaults
}>()
const list = ref([] as Alias[])
const isLoading = ref(false)
const error = ref<string | null>(null)
const copyText = ref('Click to copy')
const search = ref('')
const searchQuery = ref('')

const fetchAliases = async () => {
    error.value = null
    searchQuery.value = search.value.trim()

    try {
        isLoading.value = true
        const res = await api.fetchAliases(props.apiToken, searchQuery.value)
        list.value = res.aliases
        console.log('Fetched aliases:', res.aliases)
    } catch (err) {
        error.value = 'An unexpected error occurred'
        console.error('Fetch aliases error:', err)
    } finally {
        isLoading.value = false
    }
}

const clearSearch = () => {
    search.value = ''
    searchQuery.value = ''
    fetchAliases()
}

const updateAlias = async (alias: Alias) => {
    alias.enabled = !alias.enabled
    try {
        await api.updateAlias(props.apiToken, alias.id, alias)
        console.log('Updated alias:', alias)
    } catch (err) {
        console.error('Update alias error:', err)
    }
}

const deleteAlias = async (aliasId: string) => {
    if (!confirm('Are you sure you want to delete alias?')) return

    try {
        await api.deleteAlias(props.apiToken, aliasId)
        list.value = list.value.filter(alias => alias.id !== aliasId)
        console.log('Deleted alias with ID:', aliasId)
    } catch (err) {
        console.error('Delete alias error:', err)
    }
}

const copyAlias = (alias: string) => {
    navigator.clipboard.writeText(alias)
    copyText.value = 'Copied'
    setTimeout(() => {
        copyText.value = 'Click to copy'
    }, 2000)
}

const onCreateAlias = (event: { alias: Alias }) => {
    if (!event.alias) return
    console.log('Alias created event received:', event.alias)
    list.value.unshift(event.alias)
}

onMounted(() => {
    fetchAliases()
    events.on('alias.create', onCreateAlias)
})

onUpdated(() => {
    tooltip.autoInit()
})
</script>
