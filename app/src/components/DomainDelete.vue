<template>
    <div>
        <div v-bind:id="'modal-delete-domain' + domain.id" class="hs-overlay hidden">
            <div>
                <div>
                    <header>
                        <button @click="close" class="close">
                            <i class="icon arrow-left-line icon-primary"></i>
                        </button>
                        <h4>DELETE DOMAIN</h4>
                    </header>
                    <article>
                        <div>
                            <div class="mb-5">
                                <p>
                                    <strong>WARNING:</strong> This operation cannot be undone. Deleting this domain will also remove all associated aliases.
                                </p>
                            </div>
                        </div>
                    </article>
                    <footer>
                        <nav>
                            <button @click.stop="deleteDomain" class="cta delete">
                                Delete Domain
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
import { domainApi } from '../api/domain.ts'
import axios from 'axios'
import overlay from '@preline/overlay'
import events from '../events.ts'

const props = defineProps(['domain'])
const domain = ref(props.domain)
const error = ref('')

const deleteDomain = async () => {
    if (!confirm('Are you sure you want to delete this domain?')) return

    try {
        await domainApi.delete(domain.value.id)
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
    const modal = document.querySelector('#modal-delete-domain' + domain.value.id) as any
    overlay.close(modal)
}

const addEvents = () => {
    const modal = overlay.getInstance('#modal-delete-domain' + domain.value.id as any, true) as any
    modal.element.on('close', () => {
        close()
    })
}

onMounted(() => {
    overlay.autoInit()
    addEvents()
})
</script>