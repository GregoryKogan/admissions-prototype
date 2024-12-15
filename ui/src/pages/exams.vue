<template>
  <v-container>
    <v-tabs v-model="activeTab">
      <v-tab value="available">Доступные</v-tab>
      <v-tab value="history">История</v-tab>
    </v-tabs>

    <v-tabs-window v-model="activeTab" class="pt-4">
      <v-tabs-window-item value="available">
        <div class="d-flex align-center justify-space-between mb-4">
          <v-btn
            color="primary"
            @click="reloadAll"
            :icon="$vuetify.display.xs"
            variant="tonal"
          >
            <v-icon>mdi-refresh</v-icon>
            <span class="d-none d-sm-inline ml-2">Обновить</span>
          </v-btn>
        </div>

        <template v-if="availableExams.length">
          <div v-for="exam in availableExams" :key="exam.ID">
            <ExamCard :data="exam" @status-changed="delayReloadAll" />
          </div>
        </template>
        <v-alert v-else type="info" text="Нет доступных экзаменов"></v-alert>
      </v-tabs-window-item>

      <v-tabs-window-item value="history">
        <div class="d-flex align-center justify-space-between mb-4">
          <v-btn
            color="primary"
            @click="reloadAll"
            :icon="$vuetify.display.xs"
            variant="tonal"
          >
            <v-icon>mdi-refresh</v-icon>
            <span class="d-none d-sm-inline ml-2">Обновить</span>
          </v-btn>
        </div>

        <template v-if="historyExams.length">
          <div v-for="exam in historyExams" :key="exam.ID">
            <ExamCard
              :data="exam"
              :registered="true"
              @status-changed="delayReloadAll"
            />
          </div>
        </template>
        <v-alert v-else type="info" text="Нет истории экзаменов"></v-alert>
      </v-tabs-window-item>
    </v-tabs-window>
  </v-container>
</template>

<script lang="ts" setup>
import { onMounted, ref } from 'vue'
import { Exam } from '@/api/api.exams'
import ExamsService from '@/api/api.exams'
import ExamCard from '@/components/ExamCard.vue'

const activeTab = ref('available')
const availableExams = ref<Exam[]>([])
const historyExams = ref<Exam[]>([])

async function reloadAvailableExams() {
  const response = await ExamsService.available()
  availableExams.value = response.data
}

async function reloadHistoryExams() {
  const response = await ExamsService.history()
  historyExams.value = response.data
}

function reloadAll() {
  reloadAvailableExams()
  reloadHistoryExams()
}

function delayReloadAll() {
  setTimeout(reloadAll, 500)
}

onMounted(reloadAll)
</script>

<route lang="yaml">
meta:
  layout: user
</route>
