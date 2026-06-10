<template>
    <div>
        <button v-bind:data-hs-overlay="'#' + modalId" class="cta">
            New Domain
        </button>
        <div v-bind:id="modalId" class="hs-overlay hidden">
            <div>
                <div>
                    <!-- Step 1: Ownership verification -->
                    <template v-if="step === 1">
                        <header>
                            <button @click="close" class="close">
                                <i class="icon arrow-left-line icon-primary"></i>
                            </button>
                            <h4>ADD NEW DOMAIN · Step 1 of 2</h4>
                        </header>
                        <article>
                            <div class="mb-5">
                                <p>
                                    To confirm that you own the domain, add the TXT record shown below and then click Add Domain. After the domain has been successfully added, you may remove the TXT record if you wish. It may take some time for the DNS changes to propagate.
                                </p>
                                <p class="break-all">
                                    DNS Record:<br>
                                </p>
                            </div>
                            <div class="mb-5">
                                <table class="sm desktop">
                                    <thead>
                                        <tr>
                                            <th>Type</th>
                                            <th>Host</th>
                                            <th>Value</th>
                                        </tr>
                                    </thead>
                                    <tbody>
                                        <tr>
                                            <td>TXT</td>
                                            <td>@</td>
                                            <td>
                                                <div class="hs-tooltip break-all">
                                                    <div class="hs-tooltip-toggle">
                                                        <button class="plain max-w-[320px] text-[13px] p-0   plain truncate text-wrap text-end" @click="copyToClipboard('mailx-verify=' + config.verify)">
                                                            mailx-verify={{ config.verify }}
                                                        </button>
                                                        <span class="hs-tooltip-content hs-tooltip-shown:opacity-100 hs-tooltip-shown:visible" role="tooltip">
                                                            {{ copyText }}
                                                        </span>
                                                    </div>
                                                </div>
                                            </td>
                                        </tr>
                                    </tbody>
                                </table>
                                <div class="tablet">
                                    <p class="font-secondary text-sm leading-[2rem] text-black dark:text-white">
                                        TXT @ mailx-verify={{ config.verify }}
                                    </p>
                                </div>
                            </div>
                            <div class="mb-5">
                                <label for="domain_name">
                                    Domain:
                                </label>
                                <template v-if="ownershipPending">
                                    <div class="flex items-center gap-3 mt-1">
                                        <p class="font-secondary text-sm text-black dark:text-white m-0">{{ domain.name }}</p>
                                        <span class="badge progress small">PENDING</span>
                                    </div>
                                </template>
                                <template v-else>
                                    <input
                                        v-model="domain.name"
                                        v-bind:class="{ 'error': nameError }"
                                        id="domain_name"
                                        placeholder="example.net"
                                        type="text"
                                    >
                                    <p v-if="nameError" class="error">Required</p>
                                </template>
                            </div>
                        </article>
                        <footer>
                            <nav>
                                <button @click="postDomain" class="cta">
                                    {{ ownershipPending ? 'Verify ownership' : 'Add Domain' }}
                                </button>
                                <button @click="close" class="cancel">
                                    Cancel
                                </button>
                            </nav>
                            <p v-if="ownershipPending && error" class="error px-5">Unable to verify ownership. Ensure the correct TXT record is set or try again later.</p>
                            <p v-else-if="error" class="error px-5">{{ error }}</p>
                        </footer>
                    </template>

                    <!-- Step 2: DNS records verification -->
                    <template v-else>
                        <header>
                            <button @click="close" class="close">
                                <i class="icon arrow-left-line icon-primary"></i>
                            </button>
                            <h4>VERIFY DNS RECORDS · Step 2 of 2</h4>
                        </header>
                        <article>
                            <div>
                                <div class="mb-5">
                                    <p>
                                        Set the following DNS records for your domain. It may take some time for the DNS changes to propagate.
                                    </p>
                                </div>
                                <div class="mb-5">
                                    <table class="sm desktop">
                                        <thead>
                                            <tr>
                                                <th>Type</th>
                                                <th>Host</th>
                                                <th>Value</th>
                                            </tr>
                                        </thead>
                                        <tbody>
                                            <template v-for="(mx_host, index) in config.mx_hosts" :key="mx_host">
                                                <tr>
                                                    <td>MX {{ 10 * (index + 1) }}</td>
                                                    <td>@</td>
                                                    <td>
                                                        <div class="hs-tooltip inline-block">
                                                            <div class="hs-tooltip-toggle">
                                                                <button class="plain truncate max-w-[320px] text-[13px] p-0" @click="copyToClipboard(mx_host + '.')">
                                                                    {{ mx_host }}.
                                                                </button>
                                                                <span class="hs-tooltip-content hs-tooltip-shown:opacity-100 hs-tooltip-shown:visible" role="tooltip">
                                                                    {{ copyText }}
                                                                </span>
                                                            </div>
                                                        </div>
                                                    </td>
                                                </tr>
                                            </template>
                                            <tr>
                                                <td>TXT</td>
                                                <td>@</td>
                                                <td>
                                                    <div class="hs-tooltip inline-block">
                                                        <div class="hs-tooltip-toggle">
                                                            <button class="plain truncate max-w-[320px] text-[13px] p-0" @click="copyToClipboard('v=spf1 include:spf.' + config.domain + ' -all')">
                                                                v=spf1 include:spf.{{ config.domain }} -all
                                                            </button>
                                                            <span class="hs-tooltip-content hs-tooltip-shown:opacity-100 hs-tooltip-shown:visible" role="tooltip">
                                                                {{ copyText }}
                                                            </span>
                                                        </div>
                                                    </div>
                                                </td>
                                            </tr>
                                            <template v-for="selector in config.dkim_selectors" :key="selector">
                                                <tr>
                                                    <td>CNAME</td>
                                                    <td>{{ selector }}._domainkey</td>
                                                    <td>
                                                        <div class="hs-tooltip inline-block">
                                                            <div class="hs-tooltip-toggle">
                                                                <button class="plain truncate max-w-[320px] text-[13px] p-0" @click="copyToClipboard(selector + '._domainkey.' + config.domain + '.')">
                                                                    {{ selector }}._domainkey.{{ config.domain }}.
                                                                </button>
                                                                <span class="hs-tooltip-content hs-tooltip-shown:opacity-100 hs-tooltip-shown:visible" role="tooltip">
                                                                    {{ copyText }}
                                                                </span>
                                                            </div>
                                                        </div>
                                                    </td>
                                                </tr>
                                            </template>
                                            <tr>
                                                <td>TXT</td>
                                                <td>_dmarc</td>
                                                <td>
                                                    <div class="hs-tooltip inline-block">
                                                        <div class="hs-tooltip-toggle">
                                                            <button class="plain truncate max-w-[320px] text-[13px] p-0" @click="copyToClipboard('v=DMARC1; p=quarantine; adkim=s')">
                                                                v=DMARC1; p=quarantine; adkim=s
                                                            </button>
                                                            <span class="hs-tooltip-content hs-tooltip-shown:opacity-100 hs-tooltip-shown:visible" role="tooltip">
                                                                {{ copyText }}
                                                            </span>
                                                        </div>
                                                    </div>
                                                </td>
                                            </tr>
                                        </tbody>
                                    </table>
                                    <div class="tablet">
                                        <p class="font-secondary text-sm leading-[2rem] text-black dark:text-white">
                                            <template v-for="(mx_host, index) in config.mx_hosts" :key="mx_host">
                                                MX {{ 10 * (index + 1) }} {{ mx_host }}.<br>
                                            </template>
                                            TXT @ v=spf1 include:spf.{{ config.domain }}. -all <br>
                                            <template v-for="selector in config.dkim_selectors" :key="selector">
                                                CNAME {{ selector }}._domainkey {{ selector }}._domainkey.{{ config.domain }}. <br>
                                            </template>
                                            TXT _dmarc v=DMARC1; p=quarantine; adkim=s
                                        </p>
                                    </div>
                                </div>
                            </div>
                        </article>
                        <footer>
                            <nav>
                                <button @click="verifyDns" class="cta">
                                    Verify DNS Records
                                </button>
                                <button @click="close" class="cancel">
                                    Cancel
                                </button>
                            </nav>
                            <p v-if="step2Error" class="error px-5">Unable to verify domain DNS records. {{ step2Error }} or try again later.</p>
                        </footer>
                    </template>
                </div>
            </div>
        </div>
    </div>
</template>

<script setup lang="ts">
import { ref, onMounted, getCurrentInstance } from 'vue'
import overlay from '@preline/overlay'
import axios from 'axios'
import { domainApi } from '../api/domain.ts'
import events from '../events.ts'
import tooltip from '@preline/tooltip'

const modalId = 'modal-create-domain-' + getCurrentInstance()!.uid

const config = ref({
    verify: '',
    domain: '',
    dkim_selectors: [] as string[],
    mx_hosts: [] as string[],
})
const domain = ref({
    name: '',
})
const step = ref(1)
const ownershipPending = ref(false)
const createdDomain = ref({ id: '' })
const error = ref('')
const step2Error = ref('')
const nameError = ref(false)
const copyText = ref('Click to copy')

const validateName = () => {
    nameError.value = !domain.value.name
    return !nameError.value
}

const getConfig = async () => {
    try {
        const res = await domainApi.getConfig()
        config.value = res.data
        setTimeout(() => {
            tooltip.autoInit()
        }, 0)
    } catch (err) {
        if (axios.isAxiosError(err)) {
            error.value = err.response?.data.error || err.message
        }
    }
}

const postDomain = async () => {
    if (!ownershipPending.value && !validateName()) {
        return
    }

    const payload = {
        name: domain.value.name,
    }

    try {
        const res = await domainApi.create(payload)
        createdDomain.value = res.data
        error.value = ''
        ownershipPending.value = false
        // events.emit('domain.create', {})
        step.value = 2
    } catch (err) {
        if (axios.isAxiosError(err)) {
            if (err.response?.status === 429) {
                error.value = 'Too many requests, please try again later.'
            } else {
                ownershipPending.value = true
                error.value = err.response?.data.error || err.message
            }
        }
    }
}

const verifyDns = async () => {
    try {
        await domainApi.verifyDns(createdDomain.value.id)
        step2Error.value = ''
        close()
    } catch (err) {
        if (axios.isAxiosError(err)) {
            step2Error.value = err.response?.data.error || err.message
        }
    }
}

const close = () => {
    domain.value = { name: '' }
    error.value = ''
    step2Error.value = ''
    nameError.value = false
    step.value = 1
    ownershipPending.value = false
    createdDomain.value = { id: '' }
    document.removeEventListener('keydown', handleKeydown)
    events.emit('domain.reload', {})
    const modal = document.querySelector('#' + modalId) as any
    overlay.close(modal)
}

const addEvents = () => {
    const modal = overlay.getInstance(('#' + modalId) as any, true) as any
    modal.element.on('close', () => {
        close()
    })
    modal.element.on('open', () => {
        document.addEventListener('keydown', handleKeydown)
        focusFirstInput()
        getConfig()
        tooltip.autoInit()
    })
}

const focusFirstInput = () => {
    const input = document.getElementById('domain_name')
    input?.focus()
}

const handleKeydown = (event: KeyboardEvent) => {
    if (event.key === 'Enter') {
        event.preventDefault()
        if (step.value === 1) {
            postDomain()
        } else {
            verifyDns()
        }
    }
}

const copyToClipboard = (txt: string) => {
    navigator.clipboard.writeText(txt)
    copyText.value = 'Copied'
    setTimeout(() => {
        copyText.value = 'Click to copy'
    }, 2000)
}

onMounted(() => {
    overlay.autoInit()
    addEvents()
})
</script>