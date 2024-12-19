<template>
  <v-container>
    <v-tabs v-model="activeTab">
      <v-tab value="pending">Ожидающие</v-tab>
      <v-tab value="accepted">Принятые</v-tab>
    </v-tabs>

    <v-tabs-window v-model="activeTab" class="pt-4">
      <v-tabs-window-item value="pending">
        <div class="d-flex align-center justify-space-between mb-4">
          <h2 class="text-h5">Ожидающие регистрации</h2>
          <v-btn
            color="primary"
            @click="reloadPending"
            :icon="$vuetify.display.xs"
            variant="tonal"
          >
            <v-icon>mdi-refresh</v-icon>
            <span class="d-none d-sm-inline ml-2">Обновить</span>
          </v-btn>
        </div>

        <template v-if="pendingRegistrations.length">
          <Registration
            v-for="registration in pendingRegistrations"
            :key="registration.ID"
            :data="registration"
            style="margin-bottom: 16px"
            @status-changed="statusChanged"
          />
        </template>
        <v-alert v-else type="info" text="Нет ожидающих регистраций"></v-alert>
      </v-tabs-window-item>

      <v-tabs-window-item value="accepted">
        <div class="d-flex align-center justify-space-between mb-4">
          <h2 class="text-h5">Принятые регистрации</h2>
          <div>
            <v-btn
              color="primary"
              @click="reloadAccepted"
              :icon="$vuetify.display.xs"
              variant="tonal"
            >
              <v-icon>mdi-refresh</v-icon>
              <span class="d-none d-sm-inline ml-2">Обновить</span>
            </v-btn>
          </div>
        </div>

        <v-container class="d-flex justify-start">
          <v-btn
            color="primary"
            @click="downloadAccepted"
            variant="tonal"
            class="mr-2"
          >
            <v-icon>mdi-download</v-icon>
            <span class="d-none d-sm-inline ml-2">Скачать CSV</span>
          </v-btn>
        </v-container>

        <template v-if="acceptedRegistrations.length">
          <Registration
            v-for="registration in acceptedRegistrations"
            :key="registration.ID"
            :data="registration"
            :isPending="false"
            style="margin-bottom: 16px"
          />
        </template>
        <v-alert v-else type="info" text="Нет принятых регистраций"></v-alert>
      </v-tabs-window-item>
    </v-tabs-window>
  </v-container>
</template>

<script lang="ts">
import { defineComponent } from 'vue'
import RegistrationService, { Registration } from '@/api/api.registration'

export default defineComponent({
  data() {
    return {
      activeTab: 'pending',
      pendingRegistrations: [] as Registration[],
      acceptedRegistrations: [] as Registration[],
    }
  },
  mounted() {
    this.reload()
  },
  methods: {
    async reload() {
      await Promise.all([this.reloadPending(), this.reloadAccepted()])
    },
    async reloadPending() {
      const response = await RegistrationService.list()
      this.pendingRegistrations = response.data
    },
    async reloadAccepted() {
      const response = await RegistrationService.accepted()
      this.acceptedRegistrations = response.data
    },
    statusChanged() {
      this.reload()
    },
    async downloadAccepted() {
      try {
        const response = await RegistrationService.downloadAccepted()
        const url = window.URL.createObjectURL(
          new Blob([response.data], { type: 'text/csv' })
        )
        const link = document.createElement('a')
        link.href = url
        link.setAttribute('download', 'accepted_registrations.csv')
        document.body.appendChild(link)
        link.click()
        document.body.removeChild(link)
      } catch (error) {
        console.error(error)
      }
    },
  },
})
</script>

<style scoped></style>

<route lang="yaml">
meta:
  layout: admin
</route>
