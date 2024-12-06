<template>
  <v-container>
    <v-card>
      <v-card-title>Список регистраций</v-card-title>
      <v-card-text>
        <p>Здесь будет таблица с регистрациями</p>
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
