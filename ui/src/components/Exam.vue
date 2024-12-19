<template>
  <v-container>
    <v-card>
      <v-card-item>
        <v-row align="center" no-gutters>
          <v-col>
            <div class="text-h6 mb-1">{{ capitalizedTitle }}</div>
            <div class="text-subtitle-1">{{ props.data.location }}</div>
          </v-col>
          <v-col cols="auto">
            <v-btn
              color="primary"
              variant="text"
              icon
              @click="downloadRegistrations"
            >
              <v-tooltip activator="parent"> Скачать регистрации </v-tooltip>
              <v-icon>mdi-download</v-icon>
            </v-btn>
            <v-btn
              color="error"
              variant="text"
              icon="mdi-delete"
              @click="deleteExam"
            ></v-btn>
          </v-col>
        </v-row>
      </v-card-item>

      <v-card-text>
        <v-row>
          <v-col cols="12" sm="6">
            <div class="d-flex align-center mb-2">
              <v-icon icon="mdi-calendar" class="mr-2"></v-icon>
              <span>{{ formattedDate }}</span>
            </div>
            <div class="d-flex align-center mb-2">
              <v-icon icon="mdi-clock-outline" class="mr-2"></v-icon>
              <span>{{ formattedTime }}</span>
            </div>
          </v-col>
          <v-col cols="12" sm="6">
            <div class="d-flex align-center">
              <v-icon icon="mdi-school" class="mr-2"></v-icon>
              <span>{{ props.data.grade }} класс</span>
            </div>
            <div class="d-flex flex-column mb-2">
              <div class="d-flex align-center mb-1">
                <v-icon icon="mdi-account-group" class="mr-2"></v-icon>
                <span>
                  {{
                    allocation
                      ? `${allocation.occupied}/${allocation.capacity}`
                      : '...'
                  }}
                  мест
                </span>
              </div>
              <v-progress-linear
                v-if="allocation"
                :model-value="(allocation.occupied / allocation.capacity) * 100"
                :color="progressColor"
                height="8"
              ></v-progress-linear>
              <v-progress-linear v-else indeterminate></v-progress-linear>
            </div>
          </v-col>
        </v-row>
      </v-card-text>
    </v-card>
  </v-container>

  <v-dialog v-model="deleteDialog" width="auto">
    <v-card>
      <v-card-title>Подтверждение удаления</v-card-title>
      <v-card-text>Вы уверены, что хотите удалить этот экзамен?</v-card-text>
      <v-card-actions>
        <v-spacer></v-spacer>
        <v-btn color="error" @click="confirmDelete">Удалить</v-btn>
        <v-btn color="grey" variant="text" @click="deleteDialog = false">
          Отмена
        </v-btn>
      </v-card-actions>
    </v-card>
  </v-dialog>
</template>

<script lang="ts" setup>
import { computed, ref, onMounted } from 'vue'
import { Exam, Allocation } from '@/api/api.exams'
import ExamsService from '@/api/api.exams'

const props = defineProps<{
  data: Exam
}>()

const emit = defineEmits<{
  (e: 'status-changed'): void
}>()

const formattedDate = computed(() => {
  const start = new Date(props.data.start)
  return start.toLocaleDateString('ru-RU', {
    day: 'numeric',
    month: 'long',
    year: 'numeric',
  })
})

const formattedTime = computed(() => {
  const start = new Date(props.data.start)
  const end = new Date(props.data.end)
  return `${start.toLocaleTimeString('ru-RU', {
    hour: '2-digit',
    minute: '2-digit',
  })} - ${end.toLocaleTimeString('ru-RU', {
    hour: '2-digit',
    minute: '2-digit',
  })}`
})

const capitalizedTitle = computed(() => {
  const title = props.data.type.title
  return title.charAt(0).toUpperCase() + title.slice(1)
})

const deleteDialog = ref(false)

const allocation = ref<Allocation | null>(null)

const progressColor = computed(() => {
  if (!allocation.value) return 'primary'
  const ratio = allocation.value.occupied / allocation.value.capacity
  if (ratio >= 0.9) return 'error'
  if (ratio >= 0.7) return 'warning'
  return 'success'
})

onMounted(async () => {
  const response = await ExamsService.allocation(props.data.ID)
  allocation.value = response.data
})

async function deleteExam() {
  deleteDialog.value = true
}

async function confirmDelete() {
  await ExamsService.delete(props.data.ID)
  deleteDialog.value = false
  emit('status-changed')
}

async function downloadRegistrations() {
  try {
    const response = await ExamsService.downloadRegistrations(props.data.ID)
    const url = window.URL.createObjectURL(
      new Blob([response.data], { type: 'text/csv' })
    )
    const link = document.createElement('a')
    link.href = url
    link.setAttribute('download', `exam_registrations_${props.data.ID}.csv`)
    document.body.appendChild(link)
    link.click()
    document.body.removeChild(link)
  } catch (error) {
    console.error(error)
  }
}
</script>
