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
              v-if="!props.registered"
              color="primary"
              variant="tonal"
              :loading="isRegistering"
              @click="register"
            >
              Записаться
            </v-btn>
            <v-chip v-else color="success">
              <v-icon start icon="mdi-check"></v-icon>
              Вы записаны
            </v-chip>
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
        </v-row>
        <v-row>
          <v-col cols="12" sm="6">
            <div class="d-flex align-center mb-2">
              <v-icon icon="mdi-school" class="mr-2"></v-icon>
              <span>{{ props.data.grade }} класс</span>
            </div>
            <div class="d-flex flex-column">
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
</template>

<script lang="ts" setup>
import { computed, ref, onMounted } from 'vue'
import { Exam, Allocation } from '@/api/api.exams'
import ExamsService from '@/api/api.exams'

const props = defineProps<{
  data: Exam
  registered?: boolean
}>()

const emit = defineEmits<{
  (e: 'status-changed'): void
}>()

const isRegistering = ref(false)
const allocation = ref<Allocation | null>(null)

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

const progressColor = computed(() => {
  if (!allocation.value) return 'primary'
  const ratio = allocation.value.occupied / allocation.value.capacity
  if (ratio >= 0.9) return 'error'
  if (ratio >= 0.7) return 'warning'
  return 'success'
})

async function register() {
  try {
    isRegistering.value = true
    await ExamsService.register(props.data.ID)
    emit('status-changed')
  } finally {
    isRegistering.value = false
  }
}

onMounted(async () => {
  const response = await ExamsService.allocation(props.data.ID)
  allocation.value = response.data
})
</script>
