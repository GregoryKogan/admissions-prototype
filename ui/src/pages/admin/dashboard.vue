<template>
  <v-container>
    <v-card>
      <v-card-title>Панель администратора</v-card-title>
      <v-card-text>
        <p>Здесь будет статистика</p>
      </v-card-text>
    </v-card>
  </v-container>
</template>

<script lang="ts">
import { defineComponent } from 'vue'

export default defineComponent({
  data: () => ({}),
})
</script>

<style scoped></style>

<route lang="yaml">
meta:
  layout: admin
</route>
