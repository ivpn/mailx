<template>
    <div>
        <div v-bind:id="'modal-verify-domain' + domain.id" class="hs-overlay hidden">
            <div>
                <div>
                    <header>
                        <button @click="close" class="close">
                            <i class="icon arrow-left-line icon-primary"></i>
                        </button>
                        <h4>VERIFY DNS RECORDS</h4>
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
                            <button @click.stop="verifyDomain" class="cta">
                                Verify DNS Records
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
import { domainApi } from '../api/domain.ts'
import axios from 'axios'
import events from '../events.ts'
import tooltip from '@preline/tooltip'

const props = defineProps(['domain'])
const domain = ref(props.domain)
const error = ref('')
const copyText = ref('Click to copy')

const config = ref({
    verify: '',
    domain: '',
    dkim_selectors: [] as string[],
    mx_hosts: [] as string[],
})

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

const verifyDomain = async () => {
    try {
        await domainApi.verifyDns(domain.value.id)
        error.value = ''
        events.emit('domain.reload', {})
        close()
    } catch (err) {
        if (axios.isAxiosError(err)) {
            error.value = err.response?.data.error || err.message
        }
    }
}

const close = () => {
    error.value = ''
    const modal = document.querySelector('#modal-verify-domain' + domain.value.id) as any
    overlay.close(modal)
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
    addEvents()
    getConfig()
})
</script>