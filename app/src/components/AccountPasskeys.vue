<template>
    <div class="mb-5">
        <h2>Passkeys</h2>
        <div v-if="passkeySupported">
            <p>
                Add or remove Passkeys associated with your account.<br>
            </p>
            <div class="flex justify-start items-center gap-x-3 mb-3">
                <button @click="addPasskey" class="cta">
                    New Passkey
                </button>
            </div>
            <p v-if="error" class="error mt-6 mb-4">Error: {{ error }}</p>
        </div>
        <div v-if="!passkeySupported">
            <p>
                Your browser/device does not support adding Passkeys.<br>
            </p>
        </div>
        <div v-if="list.length" class="table-container">
            <table>
                <thead class="desktop">
                    <tr>
                        <th>Created</th>
                        <th>ID</th>
                        <th>Actions</th>
                    </tr>
                </thead>
                <tbody>
                    <tr v-for="cred in list" :key="rowKey" class="desktop">
                        <td>
                            {{ new Date(cred.created_at).toDateString() }}
                        </td>
                        <td>
                            {{ cred.id }}
                        </td>
                        <td>
                            <button @click.stop="deleteCred(cred.id)" class="delete w-full flex items-center gap-x-2 py-2 place-content-end">
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
                                    <button @click.stop="deleteCred(cred.id)" class="delete w-full flex items-center gap-x-2 py-2 place-content-end">
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
import { startRegistration, browserSupportsWebAuthn } from '@simplewebauthn/browser'

const credential = {
    id: '',
    created_at: '',
}

const list = ref([] as typeof credential[])
const error = ref('')
const passkeySupported = ref(false)
const rowKey = ref(0)

const getList = async () => {
    try {
        const res = await userApi.getCredentials()
        list.value = res.data
        error.value = ''
        renderRow()
    } catch (err) {
        if (axios.isAxiosError(err)) {
            error.value = err.message
        }
    }
}

const deleteCred = async (id: string) => {
    if (!confirm('Are you sure you want to delete Passkey?')) return

    try {
        await userApi.deleteCredential(id)
        list.value = list.value.filter((cred: any) => cred.id !== id)
        error.value = ''
        renderRow()
    } catch (err) {
        if (axios.isAxiosError(err)) {
            error.value = err.message
        }
    }
}

const addPasskey = async () => {
    try {
        var res = await userApi.registerAdd()
        startAddPasskey(res)
    } catch (err) {
        if (axios.isAxiosError(err)) {
            error.value = err.response?.data.error || err.message

            if (err.response?.status === 429) {
                error.value = 'Too many requests, please try again later.'
            }
        }
    }
}

const startAddPasskey = async (res: any) => {
    try {
        const creds = await startRegistration({ optionsJSON: res.data['publicKey'] })
        res = await userApi.registerAddFinish(creds)
        error.value = ''
        getList()
    } catch (err: any) {
        if (axios.isAxiosError(err)) {
            error.value = err.response?.data.error || err.message

            if (err.response?.status === 429) {
                error.value = 'Too many requests, please try again later.'
            }
        } else {
            error.value = 'The operation was aborted or failed.'
        }
    }
}

const renderRow = () => {
    rowKey.value++
}

onMounted(() => {
    getList()
    passkeySupported.value = browserSupportsWebAuthn()
})
</script>