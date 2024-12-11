<template>
  <v-container>
    <v-tabs v-model="activeTab">
      <v-tab value="list">Список</v-tab>
      <v-tab value="create">Создать</v-tab>
    </v-tabs>

    <v-tabs-window v-model="activeTab" class="pt-4">
      <v-tabs-window-item value="list">
        <div class="d-flex align-center justify-space-between mb-4">
          <h2 class="text-h5">Экзамены</h2>
          <v-btn
            color="primary"
            @click="reload"
            :icon="$vuetify.display.xs"
            variant="tonal"
          >
            <v-icon>mdi-refresh</v-icon>
            <span class="d-none d-sm-inline ml-2">Обновить</span>
          </v-btn>
        </div>

        <template v-if="exams.length">
          <div v-for="exam in exams" :key="exam.ID" class="mb-4">
            <Exam :data="exam" @status-changed="statusChanged" />
          </div>
        </template>
        <v-alert v-else type="info" text="Нет актуальных экзаменов"></v-alert>
      </v-tabs-window-item>

      <v-tabs-window-item value="create">
        <CreateExam @exam-created="examCreated" />
      </v-tabs-window-item>
    </v-tabs-window>
  </v-container>
</template>

<script lang="ts">
import ExamsService, { Exam } from '@/api/api.exams'
import { defineComponent } from 'vue'

export default defineComponent({
  setup() {
    return {}
  },
  data: () => ({
    exams: [] as Exam[],
    activeTab: 0,
  }),
  mounted() {
    this.reload()
  },
  methods: {
    async reload() {
      const response = await ExamsService.list()
      this.exams = response.data
    },
    examCreated() {
      this.activeTab = 0
      this.statusChanged()
    },
    statusChanged() {
      setTimeout(() => this.reload(), 500)
    },
  },
})
</script>

<style scoped></style>

<route lang="yaml">
meta:
  layout: admin
</route>
