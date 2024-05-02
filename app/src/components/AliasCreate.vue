<template>
    <div>
        <button v-bind:data-hs-overlay="'#hs-modal-create-alias'"
            class="mt-3 py-2 pl-2 pr-3 inline-flex justify-center items-center gap-x-2 text-sm font-medium rounded-md border border-transparent bg-violet-600 text-white hover:bg-violet-700">
            <svg class="flex-shrink-0 size-3.5" xmlns="http://www.w3.org/2000/svg" width="24" height="24"
                viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="3" stroke-linecap="round"
                stroke-linejoin="round">
                <path d="M5 12h14"></path>
                <path d="M12 5v14"></path>
            </svg>
            Create Alias
        </button>
        <div v-bind:id="'hs-modal-create-alias'"
            class="hs-overlay hidden size-full fixed top-0 start-0 z-[60] overflow-x-hidden overflow-y-auto pointer-events-none">
            <div
                class="hs-overlay-open:opacity-100 hs-overlay-open:duration-500 opacity-0 transition-all sm:max-w-lg sm:w-full m-3 sm:mx-auto">
                <div class="flex flex-col bg-white border shadow-sm rounded-xl pointer-events-auto">
                    <div class="flex justify-between items-center py-3 px-4 border-b">
                        <h3 class="font-bold text-gray-800">
                            Create Alias
                        </h3>
                        <button @click="close" type="button"
                            class="flex justify-center items-center size-7 text-sm font-semibold rounded-full border border-transparent text-gray-800 hover:bg-gray-100 disabled:opacity-50 disabled:pointer-events-none">
                            <span class="sr-only">Close</span>
                            <svg class="flex-shrink-0 size-4" xmlns="http://www.w3.org/2000/svg" width="24" height="24"
                                viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"
                                stroke-linecap="round" stroke-linejoin="round">
                                <path d="M18 6 6 18"></path>
                                <path d="m6 6 12 12"></path>
                            </svg>
                        </button>
                    </div>
                    <div class="p-4 whitespace-normal text-left text-base">
                        <div class="grid space-y-3 mb-5">
                            <p class="text-gray-500 text-sm font-semibold mb-1">
                                Alias format:
                            </p>
                            <div class="relative flex items-start">
                                <div class="flex items-center h-5 mt-1">
                                    <input v-model="alias.format" value="word" id="hs-radio-word" name="hs-radio-with-description" type="radio"
                                        class="form-radio border-gray-200 rounded-full text-violet-600 focus:ring-white"
                                        aria-describedby="hs-radio-word-description" checked>
                                </div>
                                <label for="hs-radio-word" class="ms-3">
                                    <span class="block text-sm font-semibold text-gray-800">Word</span>
                                    <span id="hs-radio-word-description" class="block text-sm text-gray-600">e.g.
                                        quiet.haze16@{{ alias.domain }}</span>
                                </label>
                            </div>
                            <div class="relative flex items-start">
                                <div class="flex items-center h-5 mt-1">
                                    <input v-model="alias.format" value="chars" id="hs-radio-chars" name="hs-radio-with-description" type="radio"
                                        class="form-radio border-gray-200 rounded-full text-violet-600 focus:ring-white"
                                        aria-describedby="hs-radio-chars-description">
                                </div>
                                <label for="hs-radio-chars" class="ms-3">
                                    <span class="block text-sm font-semibold text-gray-800">Random</span>
                                    <span id="hs-radio-chars-description" class="block text-sm text-gray-600">e.g.
                                        uf1h0hxi@{{ alias.domain }}</span>
                                </label>
                            </div>
                            <div class="relative flex items-start">
                                <div class="flex items-center h-5 mt-1">
                                    <input v-model="alias.format" value="uuid" id="hs-radio-uuid" name="hs-radio-with-description" type="radio"
                                        class="form-radio border-gray-200 rounded-full text-violet-600 focus:ring-white"
                                        aria-describedby="hs-radio-uuid-description">
                                </div>
                                <label for="hs-radio-uuid" class="ms-3">
                                    <span class="block text-sm font-semibold text-gray-800">UUID</span>
                                    <span id="hs-radio-uuid-description" class="block text-sm text-gray-600">e.g.
                                        550e8400-e29b-41d4-a716-446655440000@{{ alias.domain }}</span>
                                </label>
                            </div>
                        </div>
                        <div class="mb-5">
                            <label for="alias_description" class="block text-gray-500 text-sm font-semibold mb-3">
                                Description:
                            </label>
                            <input id="alias_description" v-model="alias.description"
                                class="appearance-none outline-none border-2 rounded-md w-full py-3 px-4 text-gray-700 leading-tight focus:border-violet-600 mb-2"
                                type="text">
                        </div>
                        <div class="mb-5">
                            <label for="alias_from_name" class="block text-gray-500 text-sm font-semibold mb-3">
                                From name:
                            </label>
                            <input id="alias_from_name" v-model="alias.from_name"
                                class="appearance-none outline-none border-2 rounded-md w-full py-3 px-4 text-gray-700 leading-tight focus:border-violet-600 mb-2"
                                type="text">
                        </div>
                        <div class="mb-6">
                            <label for="alias_recipient" class="block text-gray-500 text-sm font-semibold mb-3">
                                Recipient:
                            </label>
                            <select id="alias_recipient" :disabled="!recipients.length"
                                class="form-select py-2.5 px-4 pe-9 block w-full border-2 border-gray-200 rounded-lg text-gray-700 focus:border-violet-600 disabled:opacity-50 disabled:pointer-events-none outline-none">
                                <option v-for="(recipient, index) in recipients" v-bind:recipient
                                    :selected="recipient == settings.recipient || index === 0" :key="recipient">
                                    {{ recipient }}
                                </option>
                            </select>
                        </div>
                        <div class="mb-6">
                            <label class="block text-gray-500 text-sm font-semibold mb-3" for="alias_domain">
                                Domain:
                            </label>
                            <select id="alias_domain" :disabled="!domains.length"
                                class="form-select py-2.5 px-4 pe-9 block w-full border-2 border-gray-200 rounded-lg text-gray-700 focus:border-violet-600 disabled:opacity-50 disabled:pointer-events-none outline-none">
                                <option v-for="(domain, index) in domains" v-bind:domain
                                    :selected="domain == alias.domain || index === 0" :key="domain">
                                    {{ domain }}
                                </option>
                            </select>
                        </div>
                    </div>
                    <div class="flex justify-start items-center gap-x-2 py-3 px-4 border-t">
                        <button @click="postAlias"
                            class="py-2 px-3 inline-flex items-center gap-x-2 text-sm font-semibold rounded-md border border-transparent bg-violet-600 text-white hover:bg-violet-700 disabled:opacity-50 disabled:pointer-events-none">
                            Create Alias
                        </button>
                        <button @click="close"
                            class="py-2 px-3 inline-flex items-center gap-x-2 text-sm font-medium rounded-md border border-gray-200 bg-white text-gray-800 shadow-sm hover:bg-gray-50 disabled:opacity-50 disabled:pointer-events-none">
                            Cancel
                        </button>
                    </div>
                    <div class="flex items-start">
                        <p v-if="error" class="px-5 text-red-600 text-sm mb-3">{{ error }}</p>
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
import { aliasApi } from '../api/alias.ts'
import env from "../env.json"

const props = defineProps(['recipients', 'settings'])
const emit = defineEmits(['onCreateAlias'])
const alias = ref({
    description: '',
    enabled: true,
    format: '',
    from_name: '',
    recipients: '',
    domain: env.DOMAINS[0],
})
const recipients = ref(props.recipients)
const settings = ref(props.settings)
const domains = ref(env.DOMAINS)
const error = ref('')

const postAlias = async () => {
    alias.value.enabled = true

    const domainInput = document.getElementById('alias_domain') as HTMLInputElement
    alias.value.domain = domainInput.value

    const recipientInput = document.getElementById('alias_recipient') as HTMLInputElement
    alias.value.recipients = recipientInput.value

    try {
        await aliasApi.create(alias.value)
        error.value = ''
        emit('onCreateAlias')
        close()
    } catch (err) {
        if (axios.isAxiosError(err)) {
            error.value = err.message
        }
    }
}

const close = () => {
    alias.value = {} as any
    error.value = ''
    const modal = document.querySelector('#hs-modal-create-alias')
    if (modal instanceof HTMLElement) {
        overlay.close(modal)
    }
}

onMounted(() => {
    overlay.autoInit()
})
</script>