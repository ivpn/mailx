import { ref } from 'vue'

export interface StepUpMethods {
    password: boolean
    passkey: boolean
}

const isOpen = ref(false)
const methods = ref<StepUpMethods>({ password: false, passkey: false })

export const stepUpIsOpen = isOpen
export const stepUpMethods = methods

// A single pending promise is shared across concurrent callers so that simultaneous requests
// needing step-up all wait on the same modal instead of opening one each.
let pending: Promise<void> | null = null
let resolvePending: (() => void) | null = null
let rejectPending: ((reason?: any) => void) | null = null

export function requestStepUp(requestedMethods: StepUpMethods): Promise<void> {
    if (pending) return pending

    methods.value = requestedMethods
    isOpen.value = true

    pending = new Promise<void>((resolve, reject) => {
        resolvePending = resolve
        rejectPending = reject
    }).finally(() => {
        pending = null
        resolvePending = null
        rejectPending = null
    })

    return pending
}

export function completeStepUp() {
    isOpen.value = false
    resolvePending?.()
}

export function cancelStepUp() {
    isOpen.value = false
    rejectPending?.(new Error('Step-up verification was cancelled.'))
}
