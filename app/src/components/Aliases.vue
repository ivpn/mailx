<template>
    <div class="card-container">
        <header class="head">
            <h2>Aliases</h2>
            <div class="flex gap-3 items-center justify-between">
                <div class="max-lg:hidden relative">
                    <form v-if="loaded" @submit.prevent="getList" autocomplete="off">
                        <input class="search" type="text" v-model="search" placeholder="Search aliases..." id="input_search">
                    </form>
                    <button v-if="searchQuery" @click.prevent="clearSearch" class="absolute top-0 right-0 bottom-0 px-2 flex items-center justify-center">
                        <i class="icon close icon-tertiary text-base"></i>
                    </button>
                </div>
                <button v-if="recipients.length" class="cta text-nowrap" data-hs-overlay="#modal-create-alias-false">
                    New Alias
                </button>
            </div>
        </header>
        <div class="mb-7 tablet-lg">
            <div class="relative">
                <form v-if="loaded" @submit.prevent="getList" autocomplete="off">
                    <input class="search" type="text" v-model="search" placeholder="Search aliases...">
                </form>
                <button v-if="searchQuery" @click.prevent="clearSearch" class="absolute top-0 right-0 bottom-0 px-2 flex items-center justify-center">
                    <i class="icon close icon-tertiary text-base"></i>
                </button>
            </div>
        </div>
        <div v-if="!list.length && loaded && status === 'active_inactive'" class="card-empty">
            <span class="bg-secondary rounded flex items-center justify-center p-2 mb-5">
                <i class="icon at icon-accent text-2xl"></i>
            </span>
            <h4 class="mb-6">
                <span v-if="!searchQuery && !loading">You have no aliases yet</span>
                <span v-if="searchQuery || loading">No aliases found</span>
            </h4>
            <p v-if="!recipients.length" class="text-tertiary mb-6">
                To get started, first <router-link to="/account/profile">verify</router-link> your primary email address.
            </p>
            <button v-if="!searchQuery && !loading && recipients.length" class="cta" data-hs-overlay="#modal-create-alias-false">
                New Alias
            </button>
        </div>
        <div v-bind:class="{ 'hidden': (!list.length && status === 'active_inactive') || !loaded }">
            <div class="tablet-lg">
                <div class="hs-dropdown [--placement:bottom-left] mb-2">
                    <button id="hs-dropdown-alias-status-mobile" class="sort">
                        <span class="flex mr-1">Show:</span>
                        <span class="text-accent flex items-center">
                            {{ statusLabel }}
                        </span>
                    </button>
                    <div
                        class="hs-dropdown-menu hs-dropdown-open:opacity-100 hidden"
                        aria-labelledby="hs-dropdown-alias-status-mobile"
                    >
                        <button @click="setStatus('active_inactive')">Active/Inactive</button>
                        <button @click="setStatus('active')">Active</button>
                        <button @click="setStatus('inactive')">Inactive</button>
                        <button @click="setStatus('deleted')">Deleted</button>
                        <button @click="setStatus('all')">All</button>
                    </div>
                </div>
            </div>
            <div class="card-primary">
                <div class="table-container">
                    <table>
                        <thead class="desktop-lg">
                            <tr>
                                <th class="w-10 py-6">
                                    <input
                                        type="checkbox"
                                        class="checkbox-plain"
                                        ref="selectAllCheckbox"
                                        v-bind:checked="allSelected"
                                        @change="toggleSelectAll"
                                    >
                                </th>
                                <template v-if="selectedCount === 0">
                                    <th>
                                        <div class="hs-dropdown [--placement:bottom-left]">
                                            <button id="hs-dropdown-alias-status" class="sort">
                                                <span class="flex mr-1">Show:</span>
                                                <span class="text-accent flex items-center">
                                                    {{ statusLabel }}
                                                </span>
                                            </button>
                                            <div
                                                class="hs-dropdown-menu hs-dropdown-open:opacity-100 hidden"
                                                aria-labelledby="hs-dropdown-alias-status"
                                            >
                                                <button @click="setStatus('active_inactive')">Active/Inactive</button>
                                                <button @click="setStatus('active')">Active</button>
                                                <button @click="setStatus('inactive')">Inactive</button>
                                                <button @click="setStatus('deleted')">Deleted</button>
                                                <button @click="setStatus('all')">All</button>
                                            </div>
                                        </div>
                                    </th>
                                    <th>Description</th>
                                    <th>
                                        <button
                                        @click="sort"
                                        data-sort="name"
                                        class="sort">
                                            Alias
                                            <i
                                                data-sort="name"
                                                v-if="sortBy !== 'name'"
                                                v-bind:class="{'rotate-180': sortOrder === 'ASC' && sortBy === 'name' }"
                                                class="icon arrow-down text-xl icon-tertiary"
                                            ></i>
                                            <i
                                                data-sort="name"
                                                v-if="sortBy === 'name'"
                                                v-bind:class="{'rotate-180': sortOrder === 'ASC' && sortBy === 'name' }"
                                                class="icon arrow-down text-xl icon-accent"
                                            ></i>
                                        </button>
                                    </th>
                                    <th>Count</th>
                                    <th>
                                        <button
                                        @click="sort"
                                        data-sort="created_at"
                                        class="sort">
                                            Created
                                            <i
                                                data-sort="created_at"
                                                v-if="sortBy !== 'created_at'"
                                                v-bind:class="{'rotate-180': sortOrder === 'ASC' && sortBy === 'created_at' }"
                                                class="icon arrow-down text-xl icon-tertiary"
                                            ></i>
                                            <i
                                                data-sort="created_at"
                                                v-if="sortBy === 'created_at'"
                                                v-bind:class="{'rotate-180': sortOrder === 'ASC' && sortBy === 'created_at' }"
                                                class="icon arrow-down text-xl icon-accent"
                                            ></i>
                                        </button>
                                    </th>
                                    <th>Actions</th>
                                </template>
                                <th v-else colspan="6">
                                    <div class="flex items-center gap-3 flex-nowrap">
                                        <button v-bind:disabled="!canActivate || bulkLoading" @click="bulkActivate">Activate</button>
                                        <button v-bind:disabled="!canDeactivate || bulkLoading" @click="bulkDeactivate">Deactivate</button>
                                        <button v-bind:disabled="!canPin || bulkLoading" @click="bulkPin">Pin</button>
                                        <button v-bind:disabled="!canUnpin || bulkLoading" @click="bulkUnpin">Unpin</button>
                                        <button v-bind:disabled="!canDelete || bulkLoading" @click="bulkDelete" class="delete">Delete</button>
                                        <button v-bind:disabled="!canForget || bulkLoading" @click="bulkForget" class="delete">Forget</button>
                                        <button v-bind:disabled="!canRestore || bulkLoading" @click="bulkRestore">Restore</button>
                                        <span class="text-tertiary text-sm ml-auto">{{ selectedCount }} selected</span>
                                    </div>
                                </th>
                            </tr>
                        </thead>
                        <tbody>
                            <AliasRow
                                v-for="alias in list"
                                :alias="alias"
                                :key="rowKey"
                                :recipients.sync="recipients"
                                :wildcard=false
                                :selectable="true"
                                :selected="selectedIds.has(alias.id)"
                                @onToggleSelect="toggleSelectOne"
                            />
                        </tbody>

                    </table>
                </div>
                <p v-if="error" class="error">Error: {{ error }}</p>
                <Pagination v-if="list.length" :list.sync="list" :limit="limit" :page="page" :total="total" :key="rowKey" @onUpdatePage="onUpdatePage" />
            </div>
        </div>
    </div>
    <AliasCreate v-if="recipients.length && settings.id && loaded" :recipients.sync="recipients" :settings.sync="settings" :wildcard=false :label="'New Alias'" />
</template>

<script setup lang="ts">
import { onMounted, onUnmounted, ref, computed, watch, watchEffect, nextTick } from 'vue'
import axios from 'axios'
import { aliasApi } from '../api/alias'
import { settingsApi } from '../api/settings.ts'
import AliasRow from './AliasRow.vue'
import AliasCreate from './AliasCreate.vue'
import Pagination from './Pagination.vue'
import events from '../events.ts'
import { RouterLink } from 'vue-router'
import dropdown from '@preline/dropdown'

const alias = {
    id: '',
    created_at: '',
    deleted_at: null as string | null,
    name: '',
    enabled: false,
    description: '',
    recipients: '',
    from_name: '',
    pinned: false,
    is_custom_domain: false,
    stats: {
        forwards: 0,
        blocks: 0,
        replies: 0,
        sends: 0
    }
}

const list = ref([] as typeof alias[])
const recipients = ref([])
const settings = ref({
    id: '',
    domain: '',
    domains: [],
    custom_domains: [],
    recipient: '',
    from_name: ''
})
const error = ref('')
const loaded = ref(false)
const loading = ref(false)
const rowKey = ref(0)
const limit = ref(25)
const page = ref(1)
const total = ref(0)
const sortBy = ref('created_at')
const sortOrder = ref('DESC')
const search = ref('')
const searchQuery = ref('')
const status = ref('active_inactive')
const statusLabel = computed(() => {
    if (status.value === 'active') return 'Active'
    if (status.value === 'inactive') return 'Inactive'
    if (status.value === 'deleted') return 'Deleted'
    if (status.value === 'all') return 'All'
    return 'Active/Inactive'
})

const selectedIds = ref(new Set<string>())
const selectAllCheckbox = ref<HTMLInputElement | null>(null)
const bulkLoading = ref(false)

const selectedAliases = computed(() => list.value.filter(a => selectedIds.value.has(a.id)))
const selectedCount = computed(() => selectedAliases.value.length)
const allSelected = computed(() => list.value.length > 0 && selectedIds.value.size === list.value.length)
const someSelected = computed(() => selectedIds.value.size > 0 && !allSelected.value)

const canActivate = computed(() => selectedAliases.value.some(a => !a.enabled))
const canDeactivate = computed(() => selectedAliases.value.some(a => a.enabled))
const canPin = computed(() => selectedAliases.value.some(a => !a.pinned))
const canUnpin = computed(() => selectedAliases.value.some(a => a.pinned))
const canDelete = computed(() => selectedCount.value > 0 && selectedAliases.value.every(a => !a.deleted_at))
const canRestore = computed(() => selectedCount.value > 0 && selectedAliases.value.every(a => !!a.deleted_at))
const canForget = computed(() => selectedCount.value > 0 && selectedAliases.value.every(a => a.is_custom_domain))

// HTML has no declarative attribute for the indeterminate checkbox state.
watchEffect(() => {
    if (selectAllCheckbox.value) {
        selectAllCheckbox.value.indeterminate = someSelected.value
    }
})

// The status dropdown <th> is destroyed/recreated when the bulk toolbar toggles, so
// Preline's dropdown widget (bound at autoInit() time) needs to be re-bound for it to work.
watch(selectedCount, () => {
    nextTick(() => dropdown.autoInit())
})

const getList = async () => {
    loading.value = true
    selectedIds.value = new Set()
    searchQuery.value = search.value.trim()
    if (searchQuery.value) {
        page.value = 1 // Reset to first page on search
    }

    try {
        const res = await aliasApi.getList({
            limit: limit.value,
            page: page.value,
            sort_by: sortBy.value,
            sort_order: sortOrder.value,
            wildcard: false,
            search: searchQuery.value,
            status: status.value
        })
        list.value = res.data.aliases
        total.value = res.data.total
        loaded.value = true
        loading.value = false
        error.value = ''
        renderRow()
    } catch (err) {
        if (axios.isAxiosError(err)) {
            error.value = err.message
        }
    }
}

const getSettings = async () => {
    try {
        const res = await settingsApi.getDefaults()
        settings.value = res.data
        recipients.value = res.data.recipients
        error.value = ''
    } catch (err) {
        if (axios.isAxiosError(err)) {
            error.value = err.message
        }
    }
}

const deleteAlias = async (payload: any) => {
    try {
        await aliasApi.delete(payload.id)
        error.value = ''
        getList()
    } catch (err) {
        if (axios.isAxiosError(err)) {
            error.value = err.message
        }
    }
}

const forgetAlias = async (payload: any) => {
    try {
        await aliasApi.forget(payload.id)
        error.value = ''
        getList()
    } catch (err) {
        if (axios.isAxiosError(err)) {
            error.value = err.message
        }
    }
}

const renderRow = () => {
    rowKey.value++
}

const onUpdatePage = (obj: any) => {
    limit.value = obj.limit
    page.value = obj.page
    getList()
}

const onDeleteAlias = (payload: { id: string, wildcard: boolean }) => {
    deleteAlias(payload)
}

const onForgetAlias = (payload: { id: string }) => {
    forgetAlias(payload)
}

const sort = (e: any) => {
    const sort = e.target.dataset.sort
    if (sort === sortBy.value) {
        sortOrder.value = sortOrder.value === 'ASC' ? 'DESC' : 'ASC'
    } else {
        sortBy.value = sort
        sortOrder.value = 'DESC'
    }

    getList()
}

const clearSearch = () => {
    search.value = ''
    searchQuery.value = ''
    getList()
}

const setStatus = (value: string) => {
    status.value = value
    page.value = 1
    getList()
}

const handleKeydown = (event: KeyboardEvent) => {
    // Only trigger if not focused on an input or textarea
    const activeElement = document.activeElement
    if (activeElement && (activeElement.tagName === 'INPUT' || activeElement.tagName === 'TEXTAREA' || activeElement.getAttribute('contenteditable') === 'true')) {
        return
    }
    
    if (event.key === 's' || event.key === 'S') {
        event.preventDefault()
        const input = document.getElementById('input_search')
        input?.focus()
    }
    if (event.key === 'n' || event.key === 'N') {
        event.preventDefault()
        const modalTrigger = document.querySelector('[data-hs-overlay="#modal-create-alias-false"]') as HTMLElement
        modalTrigger?.click()
    }
}

const toggleSelectAll = () => {
    selectedIds.value = allSelected.value ? new Set() : new Set(list.value.map(a => a.id))
}

const toggleSelectOne = (id: string) => {
    const next = new Set(selectedIds.value)
    if (next.has(id)) {
        next.delete(id)
    } else {
        next.add(id)
    }
    selectedIds.value = next
}

const bulkActionError = (err: unknown) => {
    if (axios.isAxiosError(err)) {
        error.value = err.response?.data?.error || err.message
    }
}

const bulkUpdateEnabled = async (enabled: boolean) => {
    bulkLoading.value = true
    try {
        await aliasApi.bulkEnabled(Array.from(selectedIds.value), enabled)
        error.value = ''
        getList()
    } catch (err) {
        bulkActionError(err)
    } finally {
        bulkLoading.value = false
    }
}

const bulkUpdatePinned = async (pinned: boolean) => {
    bulkLoading.value = true
    try {
        await aliasApi.bulkPinned(Array.from(selectedIds.value), pinned)
        error.value = ''
        getList()
    } catch (err) {
        bulkActionError(err)
    } finally {
        bulkLoading.value = false
    }
}

const bulkActivate = () => bulkUpdateEnabled(true)
const bulkDeactivate = () => bulkUpdateEnabled(false)
const bulkPin = () => bulkUpdatePinned(true)
const bulkUnpin = () => bulkUpdatePinned(false)

const bulkDelete = async () => {
    if (!confirm(`Are you sure you want to delete ${selectedCount.value} alias(es)? A deleted email alias can be restored within 90 days.`)) return

    bulkLoading.value = true
    try {
        await aliasApi.bulkDelete(Array.from(selectedIds.value))
        error.value = ''
        getList()
    } catch (err) {
        bulkActionError(err)
    } finally {
        bulkLoading.value = false
    }
}

const bulkRestore = async () => {
    bulkLoading.value = true
    try {
        await aliasApi.bulkRestore(Array.from(selectedIds.value))
        error.value = ''
        getList()
    } catch (err) {
        bulkActionError(err)
    } finally {
        bulkLoading.value = false
    }
}

const bulkForget = async () => {
    if (!confirm(`WARNING: This operation cannot be undone. You will not be able to restore these ${selectedCount.value} alias(es). Are you sure you want to delete?`)) return

    bulkLoading.value = true
    try {
        await aliasApi.bulkForget(Array.from(selectedIds.value))
        error.value = ''
        getList()
    } catch (err) {
        bulkActionError(err)
    } finally {
        bulkLoading.value = false
    }
}

onMounted(async () => {
    await getSettings()
    getList()
    dropdown.autoInit()
    events.on('alias.create', getList)
    events.on('alias.update', getList)
    events.on('alias.delete', onDeleteAlias)
    events.on('alias.forget', onForgetAlias)
    document.addEventListener('keydown', handleKeydown)
})

onUnmounted(() => {
    document.removeEventListener('keydown', handleKeydown)
})
</script>