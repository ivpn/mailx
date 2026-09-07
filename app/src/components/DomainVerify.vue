<template>
    <div>
        <div v-bind:id="'modal-verify-domain' + domain.id" class="hs-overlay hidden">
            <div>
                <div>
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
                                                <td>
                                                    <div class="hs-tooltip inline-block">
                                                        <div class="hs-tooltip-toggle">
                                                            <button class="plain truncate max-w-[320px] text-[13px] p-0" @click="copyToClipboard('@')">
                                                                @
                                                            </button>
                                                            <span class="hs-tooltip-content hs-tooltip-shown:opacity-100 hs-tooltip-shown:visible" role="tooltip">
                                                                {{ copyText }}
                                                            </span>
                                                        </div>
                                                    </div>
                                                </td>
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
                                            <td>
                                                <div class="hs-tooltip inline-block">
                                                    <div class="hs-tooltip-toggle">
                                                        <button class="plain truncate max-w-[320px] text-[13px] p-0" @click="copyToClipboard('@')">
                                                            @
                                                        </button>
                                                        <span class="hs-tooltip-content hs-tooltip-shown:opacity-100 hs-tooltip-shown:visible" role="tooltip">
                                                            {{ copyText }}
                                                        </span>
                                                    </div>
                                                </div>
                                            </td>
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
                                                <td>
                                                    <div class="hs-tooltip inline-block">
                                                        <div class="hs-tooltip-toggle">
                                                            <button class="plain truncate max-w-[320px] text-[13px] p-0" @click="copyToClipboard(selector + '._domainkey')">
                                                                {{ selector }}._domainkey
                                                            </button>
                                                            <span class="hs-tooltip-content hs-tooltip-shown:opacity-100 hs-tooltip-shown:visible" role="tooltip">
                                                                {{ copyText }}
                                                            </span>
                                                        </div>
                                                    </div>
                                                </td>
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
                                            <td>
                                                <div class="hs-tooltip inline-block">
                                                    <div class="hs-tooltip-toggle">
                                                        <button class="plain truncate max-w-[320px] text-[13px] p-0" @click="copyToClipboard('_dmarc')">
                                                            _dmarc
                                                        </button>
                                                        <span class="hs-tooltip-content hs-tooltip-shown:opacity-100 hs-tooltip-shown:visible" role="tooltip">
                                                            {{ copyText }}
                                                        </span>
                                                    </div>
                                                </div>
                                            </td>
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
                                        TXT @ v=spf1 include:spf.{{ config.domain }} -all <br>
                                        <template v-for="selector in config.dkim_selectors" :key="selector">
                                            CNAME {{ selector }}._domainkey {{ selector }}._domainkey.{{ config.domain }}. <br>
                                        </template>
                                        TXT _dmarc v=DMARC1; p=quarantine; adkim=s
                                    </p>
                                </div>
                            </div>
                            <div class="mb-5" v-if="checks.length">
                                <h5 class="text-sm">Verification Results:</h5>
                                <div class="flex flex-col gap-2">
                                    <div v-for="row in checkRows" :key="row.key" class="flex items-center gap-2 text-sm">
                                        <i v-if="row.passed === true" class="icon check icon-success text-sm"></i>
                                        <i v-if="row.passed === false" class="icon close icon-error text-sm"></i>
                                        <span v-bind:class="{ 'text-tertiary': row.passed === null }">{{ row.label }}</span>
                                    </div>
                                </div>
                            </div>
                        </div>
                    </article>
                    <footer>
                        <nav>
                            <button @click.stop="verifyDomain" class="cta">
                                Verify DNS Records
                            </button>
                            <button @click="close" class="cancel">
                                {{ verified ? 'Done' : 'Cancel' }}
                            </button>
                        </nav>
                        <p v-if="message" class="success px-5">{{ message }}</p>
                        <p v-if="error" class="error px-5">Error: {{ error }}</p>
                    </footer>
                </div>
            </div>
        </div>
    </div> 
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import overlay from '@preline/overlay'
import { domainApi } from '../api/domain.ts'
import axios from 'axios'
import events from '../events.ts'
import tooltip from '@preline/tooltip'

interface DnsCheck {
    name: string
    passed: boolean
    error?: string
}

const props = defineProps(['domain'])
const domain = ref(props.domain)
const error = ref('')
const message = ref('')
const verified = ref(false)
const checks = ref<DnsCheck[]>([])
const copyText = ref('Click to copy')

const config = ref({
    verify: '',
    domain: '',
    dkim_selectors: [] as string[],
    mx_hosts: [] as string[],
})

const checkRows = computed(() => {
    const rows = [
        { key: 'mx', label: 'MX' },
        { key: 'spf', label: 'SPF' },
        ...config.value.dkim_selectors.map((selector) => ({ key: 'dkim:' + selector, label: 'DKIM (' + selector + ')' })),
        { key: 'dmarc', label: 'DMARC' },
    ]
    return rows.map((row) => {
        const check = checks.value.find((c) => c.name === row.key)
        return { ...row, passed: check ? check.passed : null }
    })
})

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

const verifyDomain = async () => {
    try {
        const res = await domainApi.verifyDns(domain.value.id)
        error.value = ''
        message.value = res.data.message
        checks.value = res.data.checks || []
        verified.value = true
    } catch (err) {
        message.value = ''
        verified.value = false
        if (axios.isAxiosError(err)) {
            error.value = err.response?.data.error || err.message
            checks.value = err.response?.data.checks || []
        }
    }
}

const close = () => {
    error.value = ''
    message.value = ''
    verified.value = false
    checks.value = []
    const modal = document.querySelector('#modal-verify-domain' + domain.value.id) as any
    overlay.close(modal)
    events.emit('domain.reload', {})
}

const addEvents = () => {
    const modal = overlay.getInstance('#modal-verify-domain' + domain.value.id as any, true) as any
    modal.element.on('close', () => {
        close()
    })
    modal.element.on('open', () => {
        tooltip.autoInit()
    })
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
    getConfig()
    addEvents()
})
</script>