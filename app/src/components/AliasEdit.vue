<template>
    <div>
        <button ref="trigger" type="button" class="hidden" data-hs-overlay="#modal-alias-edit"></button>
        <div id="modal-alias-edit" class="hs-overlay hidden">
            <div>
                <div>
                    <header>
                        <button @click="close" class="close">
                            <i class="icon arrow-left-line icon-primary"></i>
                        </button>
                        <h4>EDIT ALIAS</h4>
                    </header>
                    <article>
                        <h3 class="break-all">{{ alias?.name }}</h3>
                        <div class="mb-7">
                            <label for="alias_edit_description">
                                Description:
                            </label>
                            <input
                                id="alias_edit_description"
                                v-model="description"
                                type="text"
                            >
                        </div>
                        <div class="mb-7">
                            <label for="alias_edit_from_name">
                                From name:
                            </label>
                            <input
                                id="alias_edit_from_name"
                                v-model="fromName"
                                type="text"
                            >
                        </div>
                        <div class="mb-6">
                            <label for="alias_edit_recipient">
                                Recipient(s):
                            </label>
                            <select
                                v-model="selectRecipients"
                                id="alias_edit_recipient"
                                :disabled="!recipients.length"
                                :multiple="true"
                                data-hs-select='{
                                "placeholder": "Select recipient(s)",
                                "toggleTag": "<button type=\"button\" aria-expanded=\"false\"></button>",
                                "toggleClasses": "hs-select-disabled:pointer-events-none hs-select-disabled:opacity-50 relative py-2.5 ps-4 pe-9 flex gap-x-2 text-nowrap w-full cursor-pointer border border-primary text-secondary bg-secondary leading-tight focus:border-accent",
                                "dropdownClasses": "mt-2 z-50 w-full max-h-72 p-1 space-y-0.5 bg-primary border border-tertiary overflow-hidden overflow-y-auto",
                                "optionClasses": "py-2 px-4 w-full text-secondary cursor-pointer hover:bg-secondary bg-primary",
                                "optionTemplate": "<div class=\"flex justify-between items-center w-full\"><span data-title></span><span class=\"hidden hs-selected:block\"><svg class=\"shrink-0 size-3.5 text-accent \" xmlns=\"http:.w3.org/2000/svg\" width=\"24\" height=\"24\" viewBox=\"0 0 24 24\" fill=\"none\" stroke=\"currentColor\" stroke-width=\"2\" stroke-linecap=\"round\" stroke-linejoin=\"round\"><polyline points=\"20 6 9 17 4 12\"/></svg></span></div>",
                                "extraMarkup": "<div class=\"absolute top-1/2 end-3 -translate-y-1/2\"><svg class=\"shrink-0 size-3.5 text-secondary \" xmlns=\"http://www.w3.org/2000/svg\" width=\"24\" height=\"24\" viewBox=\"0 0 24 24\" fill=\"none\" stroke=\"currentColor\" stroke-width=\"2\" stroke-linecap=\"round\" stroke-linejoin=\"round\"><path d=\"m7 15 5 5 5-5\"/><path d=\"m7 9 5-5 5 5\"/></svg></div>"
                                }' class="hidden">
                                <option v-for="recipient in recipients"
                                    v-bind:value=recipient
                                    :key="recipient">
                                    {{ recipient }}
                                </option>
                            </select>
                            <p v-if="errorRecipients" class="error pt-3">{{ errorRecipients }}</p>
                        </div>
                    </article>
                    <footer>
                        <nav>
                            <button
                                v-if="!success"
                                @click="updateAlias"
                                v-bind:disabled="errorRecipients.length > 0"
                                class="cta">
                                Save
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
import { ref, onMounted, nextTick } from 'vue'
import overlay from '@preline/overlay'
import select from '@preline/select'
import axios from 'axios'
import { aliasApi } from '../api/alias.ts'
import events from '../events.ts'
import { initOverlays, initSelects } from '../lib/preline.ts'

const props = defineProps(['recipients'])
const recipients = ref(props.recipients)
// One instance is shared by every row, so the alias being edited is set by open() rather than
// passed as a prop; the Preline select is bound once at mount and must not be torn down.
const source = ref<any>(null)
const alias = ref<any>(null)
const description = ref('')
const fromName = ref('')
const selectRecipients = ref<any>([])
const trigger = ref<HTMLButtonElement | null>(null)
const success = ref('')
const error = ref('')
const errorRecipients = ref('')

const setSelectValue = (value: string) => {
    const multiselect = select.getInstance('#alias_edit_recipient' as any, true) as any
    multiselect?.element?.setValue(value ? value.split(',') : [])
}

const reset = () => {
    if (!source.value) return

    alias.value = Object.assign({}, source.value)
    description.value = source.value.description
    fromName.value = source.value.from_name
    selectRecipients.value = source.value.recipients
    success.value = ''
    error.value = ''
    errorRecipients.value = ''
    setSelectValue(source.value.recipients)
}

const open = async (target: any) => {
    source.value = target
    reset()
    await nextTick()
    trigger.value?.click()
}

const updateAlias = async () => {
    if (!alias.value) return

    alias.value.description = description.value
    alias.value.from_name = fromName.value
    alias.value.recipients = selectRecipients.value.toString()

    if (!validate(alias.value.recipients)) return

    try {
        const res = await aliasApi.update(alias.value.id, alias.value)
        success.value = res.data.message
        error.value = ''
        events.emit('alias.update', {})
        close()
    } catch (err) {
        if (axios.isAxiosError(err)) {
            success.value = ''
            error.value = err.message
        }
    }
}

const validate = (rcps: string) => {
    if (rcps.length === 0) {
        error.value = 'Select one or more recipients'
        return false
    }

    return true
}

const close = () => {
    overlay.close(document.querySelector('#modal-alias-edit') as any)
}

const addEvents = () => {
    const modal = overlay.getInstance('#modal-alias-edit' as any, true) as any
    modal.element.on('close', () => {
        reset()
    })
    modal.element.on('open', () => {
        focusFirstInput()
    })

    const multiselect = select.getInstance('#alias_edit_recipient' as any, true) as any
    multiselect.element.on('change', (val: any) => {
        errorRecipients.value = val.length === 0 ? 'Select one or more recipients' : ''
    })
}

const focusFirstInput = () => {
    const input = document.getElementById('alias_edit_description') as HTMLInputElement
    input?.focus()
}

defineExpose({ open })

onMounted(() => {
    initOverlays()
    initSelects()
    addEvents()
})
</script>