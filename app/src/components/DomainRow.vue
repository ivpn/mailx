<template>
    <tr class="desktop">
        <td>
            <p>{{ new Date(domain.created_at).toDateString() }}</p>
        </td>
        <td>
            <p>{{ domain.name }}</p>
        </td>
        <td>
            <p>
                <button v-if="dnsRecordsVerified()" class="cta xs success" v-bind:data-hs-overlay="'#modal-verify-domain' + domain.id">Verified</button>
                <button v-if="!dnsRecordsVerified()" class="cta xs plain" v-bind:data-hs-overlay="'#modal-verify-domain' + domain.id">Unverified</button>
            </p>
        </td>
        <td>
            <div class="flex items-center hs-tooltip">
                <input
                    @change="updateDomain"
                    v-bind:checked="domain.enabled"
                    type="checkbox"
                >
            </div>
        </td>
        <td>
            <div class="hs-dropdown [--offset:0]">
                <button v-bind:id="'hs-dropdown-domain-edit-' + domain.id">
                    <i class="icon icon-secondary more text-lg"></i>
                </button>
                <div
                    class="hs-dropdown-menu hs-dropdown-open:opacity-100 hidden"
                    v-bind:aria-labelledby="'hs-dropdown-domain-edit-' + domain.id"
                    >
                    <button v-bind:data-hs-overlay="'#modal-verify-domain' + domain.id">
                        <i class="icon icon-primary check text-xs"></i>
                        Verify
                    </button>
                    <button class="delete"
                        v-bind:data-hs-overlay="'#modal-delete-domain' + domain.id">
                        <i class="icon icon-error trash text-xs"></i>
                        Delete
                    </button>
                </div>
            </div>
        </td>
    </tr>

    <DomainVerify :domain="domain" />
    <DomainDelete :domain="domain" />
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import dropdown from '@preline/dropdown'
import { domainApi } from '../api/domain.ts'
import DomainDelete from './DomainDelete.vue'
import DomainVerify from './DomainVerify.vue'

const props = defineProps(['domain'])
const domain = ref(props.domain)

const updateDomain = async () => {
    domain.value.enabled = !domain.value.enabled
    try {
        await domainApi.update(domain.value.id, domain.value)
    } catch {}
}

const dnsRecordsVerified = () => {
    return domain.value.mx_verified_at && domain.value.send_verified_at
}

onMounted(() => {
    dropdown.autoInit()
})
</script>