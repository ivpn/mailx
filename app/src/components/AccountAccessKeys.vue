<template>
    <div class="mb-5">
        <h2>Access Keys</h2>
        <div>
            <p>
                Add or remove Access Keys associated with your account.<br>
            </p>
            <div class="flex justify-start items-center gap-x-3 mb-3">
                <AccessKeysCreate />
            </div>
            <p v-if="error" class="error mt-6 mb-4">Error: {{ error }}</p>
        </div>
        <div v-if="list.length" class="table-container">
            <table>
                <thead class="desktop">
                    <tr>
                        <th>Created</th>
                        <th>Name</th>
                        <th>Expires At</th>
                        <th>Actions</th>
                    </tr>
                </thead>
                <tbody>
                    <tr v-for="cred in list" :key="rowKey" class="desktop">
                        <td>
                            {{ new Date(cred.created_at).toDateString() }}
                        </td>
                        <td>
                            {{ cred.name }}
                        </td>
                        <td>
                            {{ cred.expires_at ? new Date(cred.expires_at).toDateString() : 'Never' }}
                        </td>
                        <td>
                            <button @click.stop="deleteAccessKey(cred.id)" class="delete w-full flex items-center gap-x-2 py-2 place-content-end">
                                <i class="icon icon-error trash text-xs"></i>
                                Delete
                            </button>
                        </td>
                    </tr>
                    <tr v-for="cred in list" :key="rowKey" class="tablet">
                        <hr>
                        <div class="flex gap-2 justify-between">
                            <div class="text-start">
                                <p class="mb-4 text-sm">{{ new Date(cred.created_at).toDateString() }}</p>
                                <div>
                                    <p class="mb-1 text-sm">ID:</p>
                                    {{ cred.id }}
                                </div>
                            </div>
                            <div class="text-end">
                                    <button @click.stop="deleteAccessKey(cred.id)" class="delete w-full flex items-center gap-x-2 py-2 place-content-end">
                                        <i class="icon icon-error trash text-xs"></i>
                                    </button>
                            </div>
                        </div>
                    </tr>
                </tbody>
            </table>
        </div>
    </div>
</template>

<script setup lang="ts">
import { onMounted, ref } from 'vue'
import axios from 'axios'
import { userApi } from '../api/user.ts'
import AccessKeysCreate from './AccessKeysCreate.vue'
import events from '../events.ts'

const credential = {
    id: '',
    created_at: '',
    name: '',
    expires_at: '',
}

const list = ref([] as typeof credential[])
const error = ref('')
const rowKey = ref(0)

const getList = async () => {
    try {
        const res = await userApi.accessKeyList()
        list.value = res.data
        error.value = ''
        renderRow()
    } catch (err) {
        if (axios.isAxiosError(err)) {
            error.value = err.message
        }
    }
}

const deleteAccessKey = async (id: string) => {
    if (!confirm('Are you sure you want to delete Access Key?')) return

    try {
        await userApi.accessKeyDelete(id)
        list.value = list.value.filter((cred: any) => cred.id !== id)
        error.value = ''
        renderRow()
    } catch (err) {
        if (axios.isAxiosError(err)) {
            error.value = err.message
        }
    }
}

const renderRow = () => {
    rowKey.value++
}

onMounted(() => {
    getList()
    events.on('accesskey.create', getList)
})
</script>