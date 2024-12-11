<template>
  <v-container>
    <v-tabs v-model="activeTab">
      <v-tab value="available">Доступные экзамены</v-tab>
      <v-tab value="mine">Мои экзамены</v-tab>
    </v-tabs>

    <v-tabs-window v-model="activeTab" class="pt-4">
      <v-tabs-window-item value="available">
        <div class="d-flex align-center justify-space-between mb-4">
          <v-btn
            color="primary"
            @click="reloadAvailable"
            :icon="$vuetify.display.xs"
            variant="tonal"
          >
            <v-icon>mdi-refresh</v-icon>
            <span class="d-none d-sm-inline ml-2">Обновить</span>
          </v-btn>
        </div>

        <template v-if="availableExams.length">
          <div v-for="exam in availableExams" :key="exam.ID" class="mb-4">
            <ExamCard :data="exam" @status-changed="reloadAll" />
          </div>
        </template>
        <v-alert v-else type="info" text="Нет доступных экзаменов"></v-alert>
      </v-tabs-window-item>

      <v-tabs-window-item value="mine">
        <div class="d-flex align-center justify-space-between mb-4">
          <v-btn
            color="primary"
            @click="reloadMine"
            :icon="$vuetify.display.xs"
            variant="tonal"
          >
            <v-icon>mdi-refresh</v-icon>
            <span class="d-none d-sm-inline ml-2">Обновить</span>
          </v-btn>
        </div>

        <template v-if="myExams.length">
          <div v-for="exam in myExams" :key="exam.ID" class="mb-4">
            <ExamCard
              :data="exam"
              :registered="true"
              @status-changed="reloadAll"
            />
          </div>
        </template>
        <v-alert
          v-else
          type="info"
          text="Вы не записаны ни на один экзамен"
        ></v-alert>
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
const myExams = ref<Exam[]>([])

async function reloadAvailable() {
  const response = await ExamsService.available()
  availableExams.value = response.data
}

async function reloadMine() {
  const response = await ExamsService.mine()
  myExams.value = response.data
}

function reloadAll() {
  setTimeout(() => {
    reloadAvailable()
    reloadMine()
  }, 500)
}

onMounted(reloadAll)
</script>

<route lang="yaml">
meta:
  layout: user
</route>
