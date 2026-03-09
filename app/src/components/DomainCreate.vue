<template>
    <div>
        <button v-bind:data-hs-overlay="'#modal-create-domain'" class="cta">
            New Domain
        </button>
        <div v-bind:id="'modal-create-domain'" class="hs-overlay hidden">
            <div>
                <div>
                    <header>
                        <button @click="close" class="close">
                            <i class="icon arrow-left-line icon-primary"></i>
                        </button>
                        <h4>NEW DOMAIN</h4>
                    </header>
                    <article>
                        <div class="mb-5">
                            <p>
                                To confirm that you own the domain, add the TXT record shown below and then click Add Domain. After the domain has been successfully added, you may remove the TXT record if you wish.
                            </p>
                            <p class="break-all">
                                Type: <span class="text-black dark:text-white">TXT</span><br>
                                Host: <span class="text-black dark:text-white">@</span><br>
                                Value: <span class="text-black dark:text-white">mailx-verify={{ config.verify }}</span>
                            </p>
                        </div>
                        <div class="mb-5">
                            <label for="domain_name">
                                Domain:
                            </label>
                            <input
                                v-model="domain.name"
                                v-bind:class="{ 'error': nameError }"
                                id="domain_name"
                                placeholder="example.net"
                                type="text"
                            >
                            <p v-if="nameError" class="error">Required</p>
                        </div>
                    </article>
                    <footer>
                        <nav>
                            <button @click="postDomain" class="cta">
                                Add Domain
                            </button>
                            <button @click="close" class="cancel">
                                Cancel
                            </button>
                        </nav>
                        <p v-if="error" class="error px-5">Error: {{ error }}</p>
                    </footer>
                </div>
            </div>
        </div>
    </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import overlay from '@preline/overlay'
import axios from 'axios'
import { domainApi } from '../api/domain.ts'
import events from '../events.ts'

const config = ref({
    verify: '',
})
const domain = ref({
    name: '',
})
const error = ref('')
const nameError = ref(false)

const validateName = () => {
    nameError.value = !domain.value.name
    return !nameError.value
}

const getConfig = async () => {
    try {
        const res = await domainApi.getConfig()
        config.value = res.data
    } catch (err) {
        if (axios.isAxiosError(err)) {
            error.value = err.response?.data.error || err.message
        }
    }
}

const postDomain = async () => {
    if (!validateName()) {
        return
    }

    const payload = {
        name: domain.value.name,
    }

    try {
        await domainApi.create(payload)
        error.value = ''
        events.emit('domain.create', {})
        close()
    } catch (err) {
        if (axios.isAxiosError(err)) {
            error.value = err.response?.data.error || err.message

            if (err.response?.status === 429) {
                error.value = 'Too many requests, please try again later.'
            }
        }
    }
}

const close = () => {
    domain.value = { name: '' }
    error.value = ''
    nameError.value = false
    document.removeEventListener('keydown', handleKeydown)
    const modal = document.querySelector('#modal-create-domain') as any
    overlay.close(modal)
}

const addEvents = () => {
    const modal = overlay.getInstance('#modal-create-domain' as any, true) as any
    modal.element.on('close', () => {
        close()
    })
    modal.element.on('open', () => {
        document.addEventListener('keydown', handleKeydown)
        focusFirstInput()
    })
}

const focusFirstInput = () => {
    const input = document.getElementById('domain_name')
    input?.focus()
}

const handleKeydown = (event: KeyboardEvent) => {
    if (event.key === 'Enter') {
        event.preventDefault()
        postDomain()
    }
}

onMounted(() => {
    overlay.autoInit()
    addEvents()
    getConfig()
})
</script>