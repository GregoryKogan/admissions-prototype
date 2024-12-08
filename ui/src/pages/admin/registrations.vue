<template>
  <v-container>
    <div class="d-flex align-center justify-space-between mb-4">
      <h2 class="text-h5">Регистрации</h2>
      <v-btn
        color="primary"
        @click="reload"
        prepend-icon="mdi-refresh"
        variant="tonal"
      >
        Обновить
      </v-btn>
    </div>

    <template v-if="registrations.length">
      <Registration
        v-for="registration in registrations"
        :key="registration"
        :data="registration"
        style="margin-bottom: 16px"
        @status-changed="reload"
      />
    </template>
    <v-alert v-else type="info" text="Нет актуальных регистраций"></v-alert>
  </v-container>
</template>

<script lang="ts">
import RegistrationService from '@/api.registration'
import { defineComponent } from 'vue'

export default defineComponent({
  data: () => ({
    registrations: [],
  }),
  methods: {
    async reload() {
      const registrations = await RegistrationService.list()
      this.registrations = registrations.data
    },
  },
  async mounted() {
    await this.reload()
  },
})
</script>

<style scoped></style>

<route lang="yaml">
meta:
  layout: admin
</route>
