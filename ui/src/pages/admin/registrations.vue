<template>
  <v-container>
    <Registration
      v-for="registration in registrations"
      :key="registration"
      :data="registration"
      style="margin-bottom: 16px"
    />
  </v-container>
</template>

<script lang="ts">
import RegistrationService from '@/api.registration'
import { defineComponent } from 'vue'

export default defineComponent({
  data: () => ({
    registrations: [],
  }),
  async mounted() {
    const registrations = await RegistrationService.list()
    this.registrations = registrations.data
  },
})
</script>

<style scoped></style>

<route lang="yaml">
meta:
  layout: admin
</route>
