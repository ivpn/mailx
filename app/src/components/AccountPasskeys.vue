<template>
    <div class="mb-5">
        <h2>Passkeys</h2>
        <div v-if="passkeySupported">
            <p>
                Add or remove Passkeys.<br>
            </p>
            <div class="flex justify-start items-center gap-x-3 mb-3">
                <button @click="addPasskey" class="cta">
                    <svg class="flex-shrink-0 size-3.5" xmlns="http://www.w3.org/2000/svg" width="24" height="24"
                        viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="3" stroke-linecap="round"
                        stroke-linejoin="round">
                        <path d="M5 12h14"></path>
                        <path d="M12 5v14"></path>
                    </svg>
                    Add Passkey
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
                <thead>
                    <tr>
                        <th>CREATED</th>
                        <th>ID</th>
                        <th>ACTIONS</th>
                    </tr>
                </thead>
                <tbody>
                    <tr v-for="cred in list" :key="rowKey">
                        <td>
                            {{ new Date(cred.created_at).toDateString() }}
                        </td>
                        <td>
                            {{ cred.id }}
                        </td>
                        <td>
                            <button @click="deleteCred(cred.id)" class="delete">
                                Delete
                            </button>
                        </td>
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
    const data = {
        email: localStorage.getItem('email')
    }

    try {
        var res = await userApi.registerAdd(data)
        startAddPasskey(res)
    } catch (err) {
        if (axios.isAxiosError(err)) {
            error.value = err.response?.data.error || err.message

            if (err.response?.status === 429) {
                error.value = 'Too many requests, please try again later'
            }
        }
    }
}

const startAddPasskey = async (res: any) => {
    try {
        const creds = await startRegistration({ optionsJSON: res.data['publicKey'] })
        res = await userApi.registerFinish(creds)
        error.value = ''
        getList()
    } catch (err: any) {
        if (axios.isAxiosError(err)) {
            error.value = err.response?.data.error || err.message

            if (err.response?.status === 429) {
                error.value = 'Too many requests, please try again later'
            }
        } else {
            error.value = 'The operation was aborted or failed'
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