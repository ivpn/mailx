<template>
    <div>
        <button ref="trigger" type="button" class="hidden" data-hs-overlay="#modal-alias-pause"></button>
        <div id="modal-alias-pause" class="hs-overlay hidden">
            <div>
                <div>
                    <header>
                        <button @click="close" class="close">
                            <i class="icon arrow-left-line icon-primary"></i>
                        </button>
                        <h4>PAUSE ALIAS</h4>
                    </header>
                    <article>
                        <p class="mb-7">
                            Emails forwarded to <strong class="break-all">{{ alias?.name }}</strong> will be discarded until the pause ends.
                            Replies and sends from this alias will keep working.
                        </p>
                        <label for="alias_pause_value">
                            Pause for:
                        </label>
                        <div class="flex gap-3">
                            <div class="w-1/3">
                                <select v-model.number="value" id="alias_pause_value" name="alias_pause_value">
                                    <option v-for="n in pauseValues" :key="n" :value="n">{{ n }}</option>
                                </select>
                            </div>
                            <div class="flex-1">
                                <select v-model="unit" id="alias_pause_unit" name="alias_pause_unit" aria-label="Unit">
                                    <option v-for="u in pauseUnits" :key="u.value" :value="u.value">{{ value === 1 ? u.singular : u.plural }}</option>
                                </select>
                            </div>
                        </div>
                        <p class="text-sm text-tertiary mt-3" data-test="pause-preview">Paused until {{ preview }}</p>
                    </article>
                    <footer>
                        <nav>
                            <button @click="pause" v-bind:disabled="loading" class="cta">
                                Pause
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
import { computed, ref, nextTick, onMounted } from 'vue'
import overlay from '@preline/overlay'
import axios from 'axios'
import { aliasApi } from '../api/alias.ts'
import events from '../events.ts'
import { initOverlays } from '../lib/preline.ts'
import { type PauseUnit, formatPausedUntil, pauseEnd, pauseUnits, pauseValues, useNow } from '../lib/aliasPause.ts'

const defaultValue = 1
const defaultUnit: PauseUnit = 'days'

// One instance is shared by every row, so the alias being paused is set by open().
const alias = ref<any>(null)
const value = ref(defaultValue)
const unit = ref<PauseUnit>(defaultUnit)
const loading = ref(false)
const error = ref('')
const trigger = ref<HTMLButtonElement | null>(null)
const now = useNow()

const preview = computed(() => formatPausedUntil(pauseEnd(value.value, unit.value, now.value)))

const open = async (target: any) => {
    alias.value = target
    value.value = defaultValue
    unit.value = defaultUnit
    error.value = ''
    now.value = new Date()
    await nextTick()
    trigger.value?.click()
}

const pause = async () => {
    if (!alias.value || loading.value) return

    loading.value = true
    try {
        await aliasApi.pause(alias.value.id, value.value, unit.value)
        error.value = ''
        close()
        events.emit('alias.update', {})
    } catch (err) {
        if (axios.isAxiosError(err)) {
            error.value = err.response?.data?.error || err.message
        }
    } finally {
        loading.value = false
    }
}

const close = () => {
    overlay.close(document.querySelector('#modal-alias-pause') as any)
}

defineExpose({ open })

onMounted(() => {
    initOverlays()
})
</script>
