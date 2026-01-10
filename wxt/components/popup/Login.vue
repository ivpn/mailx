<template>
    <div class="page center">
        <form class="card-tertiary center" @submit.prevent="" autocomplete="off">
            <article>
                <h1 class="flex justify-center text-accent mb-8">
                    <span class="logo"></span>
                </h1>
                <h4 class="text-center mb-8">Log in with Access Key from Mailx</h4>
                <div>
                    <div class="mb-5">
                        <input
                            v-model="accessKey"
                            v-bind:class="{ 'error': accessKeyError }"
                            id="access_key"
                            type="text"
                            placeholder="Access Key"
                            autocomplete="off"
                            @keypress.enter.prevent>
                        <p v-if="accessKeyError" class="error">Required</p>
                    </div>
                    <div class="flex items-center w-full">
                        <button :disabled="isLoading" @click="loginWithAccessKey" class="cta full">
                            Log In
                        </button>
                    </div>
                    <p v-if="error" class="error mt-6">Error: {{ error }}</p>
                </div>
            </article>
            <footer>
                <div>
                    <p>To create Mailx access key, go to <a href="https://app.mailx.net/account">mailx.net</a> and look for "Access Keys" section.</p>
                </div>
            </footer>
        </form>
    </div>
</template>

<script lang="ts" setup>
import { ref } from 'vue'
import { api } from '@/lib/api'
import { store } from '@/lib/store'

const isLoading = ref(false)
const error = ref<string | null>(null)
const accessKey = ref('')
const accessKeyError = ref(false)

const loginWithAccessKey = async () => {
    accessKeyError.value = false
    error.value = null
    if (!accessKey.value) {
        accessKeyError.value = true
        return
    }

    try {
        isLoading.value = true
        const res = await api.authenticate(accessKey.value)
        processResponse(res)
        console.log('Login successful:', res)
    } catch (err) {
        error.value = 'An unexpected error occurred'
        console.error('Login error:', err)
    } finally {
        isLoading.value = false
    }
}

const processResponse = (res: any) => {
    const defaults = {
        domain: res.domain,
        domains: res.domains,
        recipient: res.recipient,
        recipients: res.recipients,
        alias_format: res.alias_format,
    }
    store.setApiToken(res.token)
    store.setDefaults(defaults)
}
</script>