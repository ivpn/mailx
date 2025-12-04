<template>
    <div>
        <button v-bind:data-hs-overlay="'#modal-create-accesskey'" class="cta">
            New Access Key
        </button>
        <div v-bind:id="'modal-create-accesskey'" class="hs-overlay hidden">
            <div>
                <div>
                    <header>
                        <button @click="close" class="close">
                            <i class="icon arrow-left-line icon-primary"></i>
                        </button>
                        <h4>NEW ACCESS KEY</h4>
                    </header>
                    <div v-if="!isCreated">
                        <article>
                            <div class="mb-5">
                                <p>
                                    Create a new Access Key associated with your account.
                                </p>
                            </div>
                            <div class="mb-5">
                                <label for="accesskey_name">
                                    Name:
                                </label>
                                <input
                                    v-model="accessKey.name"
                                    v-bind:class="{ 'error': nameError }"
                                    id="accesskey_name"
                                    type="text"
                                >
                                <p v-if="nameError" class="error">Required</p>
                            </div>
                            <div class="mb-5">
                                <label for="accesskey_expires_at">
                                    Expires After:
                                </label>
                                <select v-model="accessKey.expires_at" name="accesskey_expires_at" id="accesskey_expires_at">
                                    <option value="">Never</option>
                                    <option value="7d">1 week</option>
                                    <option value="30d">1 month</option>
                                    <option value="90d">3 months</option>
                                    <option value="180d">6 months</option>
                                    <option value="365d">1 year</option>
                                </select>
                            </div>
                        </article>
                        <footer>
                            <nav>
                                <button @click="postAccessKey" class="cta">
                                    Create Access Key
                                </button>
                                <button @click="close" class="cancel">
                                    Cancel
                                </button>
                            </nav>
                            <p v-if="error" class="error px-5">Error: {{ error }}</p>
                        </footer>
                    </div>
                    <div v-else>
                        <article>
                            <div class="mb-5">
                                <p>
                                    Access Key has been created successfully.
                                </p>
                            </div>
                            <div class="mb-5">
                                <p>
                                    Please make sure to copy and store the Access Key securely as it will not be shown again:
                                </p>
                                <p class="py-4 px-5 bg-secondary text-primary break-all font-mono">
                                    {{ accessKey.token }}
                                </p>
                            </div>
                        </article>
                        <footer>
                            <nav>
                                <button @click="copyAccessKey" class="cta">
                                    {{ copyText }}
                                </button>
                                <button @click="close" class="cancel">
                                    Close
                                </button>
                            </nav>
                        </footer>
                    </div>
                </div>
            </div>
        </div>
    </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import overlay from '@preline/overlay'
import axios from 'axios'
import { userApi } from '../api/user.ts'
import events from '../events.ts'

const accessKey = ref({
    name: '',
    expires_at: '',
    token: '',
})
const error = ref('')
const nameError = ref(false)
const isCreated = ref(false)
const copyText = ref('Copy To Clipboard')

const validateName = () => {
    nameError.value = !accessKey.value.name
    return !nameError.value
}

const postAccessKey = async () => {
    if (!validateName()) {
        return
    }

    try {
        const req = {
            name: accessKey.value.name,
            expires_at: parseExpiry(accessKey.value.expires_at),
        }
        const res = await userApi.accessKeyCreate(req)
        accessKey.value.token = res.data.token
        error.value = ''
        isCreated.value = true
        events.emit('accesskey.create', {})
    } catch (err) {
        if (axios.isAxiosError(err)) {
            error.value = err.response?.data.error || err.message

            if (err.response?.status === 429) {
                error.value = 'Too many requests, please try again later.'
            }
        }
    }
}

const parseExpiry = (value: string) => {
    if (!value) return null // "Never"

    const amount = parseInt(value)    // e.g. 7
    const unit = value.replace(String(amount), "") // "d"

    const now = new Date()

    if (unit === "d") {
        now.setDate(now.getDate() + amount)
    }

    return now
}

const close = () => {
    accessKey.value = {
        name: '',
        expires_at: '',
        token: '',
    }
    error.value = ''
    nameError.value = false
    isCreated.value = false
    document.removeEventListener('keydown', handleKeydown)
    const modal = document.querySelector('#modal-create-accesskey') as any
    overlay.close(modal)
}

const addEvents = () => {
    const modal = overlay.getInstance('#modal-create-accesskey' as any, true) as any
    modal.element.on('close', () => {
        close()
    })
    modal.element.on('open', () => {
        document.addEventListener('keydown', handleKeydown)
        focusFirstInput()
    })
}

const focusFirstInput = () => {
    const input = document.getElementById('accesskey_name')
    input?.focus()
}

const handleKeydown = (event: KeyboardEvent) => {
    if (event.key === 'Enter') {
        event.preventDefault()
        postAccessKey()
    }
}

const copyAccessKey = () => {
    navigator.clipboard.writeText(accessKey.value.token)
    copyText.value = 'Copied'
    setTimeout(() => {
        copyText.value = 'Copy To Clipboard'
    }, 2000)
}

onMounted(() => {
    overlay.autoInit()
    addEvents()
})
</script>