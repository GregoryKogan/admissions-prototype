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

        <template v-if="examsStore.availableExams.length">
          <div v-for="exam in examsStore.availableExams" :key="exam.ID">
            <ExamCard :data="exam" />
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

        <template v-if="examsStore.historyExams.length">
          <div v-for="exam in examsStore.historyExams" :key="exam.ID">
            <ExamCard :data="exam" />
          </div>
        </template>
        <v-alert v-else type="info" text="Нет истории экзаменов"></v-alert>
      </v-tabs-window-item>
    </v-tabs-window>
  </v-container>
</template>

<script lang="ts" setup>
import { onMounted, ref } from 'vue'
import { useExamsStore } from '@/stores/exams'
import ExamCard from '@/components/ExamCard.vue'

const activeTab = ref('available')
const examsStore = useExamsStore()

async function reloadAll() {
  await examsStore.reloadAll()
  await examsStore.reloadAllExamData()
}

onMounted(reloadAll)
</script>

<route lang="yaml">
meta:
  layout: user
</route>
