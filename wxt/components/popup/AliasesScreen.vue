<template>
    <div class="page p-4">
        <p v-if="isLoading" class="text-secondary">Loading...</p>
        <p v-else-if="error" class="error">{{ error }}</p>
        <div v-else class="w-full">
            <header class="pb-5">
                <div class="flex gap-3 items-center justify-between">
                    <div class="relative grow">
                        <form @submit.prevent="" autocomplete="off">
                            <input class="search w-full" type="text" placeholder="Search aliases..." id="input_search">
                        </form>
                        <button @click.prevent=""
                            class="absolute top-0 right-0 bottom-0 px-2 flex items-center justify-center">
                            <i class="icon close icon-tertiary text-base"></i>
                        </button>
                    </div>
                    <button class="cta sm text-nowrap" data-hs-overlay="#modal-create-alias-false">
                        New Alias
                    </button>
                </div>
            </header>
            <hr class="m-0">
            <div v-for="alias in list" :key="alias.id" class="py-3 border-b border-secondary flex items-center gap-x-4">
                <div class="flex items-center">
                    <input @change="updateAlias(alias)" v-bind:checked="alias.enabled" type="checkbox" class="xs">
                </div>
                <div class="grow font-medium text-sm">
                    {{ alias.name }}
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
</template>

<script lang="ts" setup>
import { ref, onMounted, onUpdated } from 'vue'
import { api } from '@/lib/api'
import { Alias } from '@/lib/types'
import tooltip from '@preline/tooltip'

const props = defineProps<{ apiToken: string }>()
const list = ref([] as Alias[])
const isLoading = ref(false)
const error = ref<string | null>(null)
const copyText = ref('Click to copy')

const fetchAliases = async () => {
    error.value = null

    try {
        isLoading.value = true
        const res = await api.fetchAliases(props.apiToken)
        list.value = res.aliases
        console.log('Fetched aliases:', res.aliases)
    } catch (err) {
        error.value = 'An unexpected error occurred'
        console.error('Fetch aliases error:', err)
    } finally {
        isLoading.value = false
    }
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

onMounted(() => {
    fetchAliases()
})

onUpdated(() => {
    tooltip.autoInit()
})
</script>
