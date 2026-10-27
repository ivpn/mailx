<template>
    <div class="container mx-auto max-w-screen-lg sm:p-10 p-5">
        <h1>Trackers Report</h1>
        <div v-if="report">
            <p>
                Mailx removed {{ report.total }} {{ report.total === 1 ? 'tracker' : 'trackers' }} from this email before forwarding it.
                Trackers like this let senders see when you open an email and learn your IP address, device and approximate location.
            </p>
            <p v-if="report.total > report.trackers.length" class="text-tertiary">
                Showing the first {{ report.trackers.length }} of {{ report.total }} trackers.
            </p>
            <div class="table-container mb-6">
                <table class="text-auto">
                    <thead>
                        <tr>
                            <th class="start">Domain</th>
                            <th>Element</th>
                            <th>Reason</th>
                        </tr>
                    </thead>
                    <tbody>
                        <tr v-for="(tracker, index) in report.trackers" :key="index">
                            <td class="start text-wrap break-all">{{ tracker.domain || 'Unknown' }}</td>
                            <td>{{ elementLabels[tracker.element] ?? tracker.element }}</td>
                            <td>{{ reasonLabels[tracker.reason] ?? tracker.reason }}</td>
                        </tr>
                    </tbody>
                </table>
            </div>
        </div>
        <p v-else class="error">
            This report link is invalid or incomplete. Open it again from the "View report" link in the forwarded email.
        </p>
        <p class="text-sm text-tertiary">
            The details of this report are contained in the link in the forwarded email and are only shown in your browser.
            They are never sent to Mailx servers, and Mailx doesn't keep a copy.
            You can turn off tracker blocking in Settings.
        </p>
        <br>
        <Footer />
    </div>
</template>

<script setup lang="ts">
import { onMounted, onUnmounted, ref } from 'vue'
import Footer from './Footer.vue'
import { decodeTrackersReport, elementLabels, reasonLabels, type TrackersReport } from '../lib/trackersReport.ts'

const report = ref<TrackersReport | null>(null)

const readReport = () => {
    report.value = decodeTrackersReport(window.location.hash)
}

onMounted(() => {
    readReport()
    window.addEventListener('hashchange', readReport)
})

onUnmounted(() => {
    window.removeEventListener('hashchange', readReport)
})
</script>

<style scoped>
p {
    @apply leading-7 mb-4;
}
</style>
