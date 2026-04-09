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
                                                <td>{{ mx_host }}.</td>
                                            </tr>
                                        </template>
                                        <tr>
                                            <td>TXT</td>
                                            <td>@</td>
                                            <td>v=spf1 include:spf.{{ config.domain }} -all</td>
                                        </tr>
                                        <template v-for="selector in config.dkim_selectors" :key="selector">
                                            <tr>
                                                <td>CNAME</td>
                                                <td>{{ selector }}._domainkey</td>
                                                <td>{{ selector }}._domainkey.{{ config.domain }}.</td>
                                            </tr>
                                        </template>
                                        <tr>
                                            <td>TXT</td>
                                            <td>_dmarc</td>
                                            <td>v=DMARC1; p=quarantine; adkim=s</td>
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

const props = defineProps(['domain'])
const domain = ref(props.domain)
const error = ref('')

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
}

onMounted(() => {
    overlay.autoInit()
    addEvents()
    getConfig()
})
</script>