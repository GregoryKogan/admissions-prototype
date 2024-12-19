<template>
  <v-container>
    <v-card :class="{ 'registered-border': registered }">
      <v-card-item>
        <v-row align="center" no-gutters justify="space-between">
          <v-col cols="8" sm="6">
            <div class="text-h6 mb-1">{{ capitalizedTitle }}</div>
            <div class="text-subtitle-1">{{ data.location }}</div>
          </v-col>
          <v-btn
            v-if="!registered"
            color="primary"
            variant="tonal"
            :loading="isRegistering"
            @click="register"
            :disabled="registeredSameType"
          >
            {{ registeredSameType ? 'Выбран другой слот' : 'Записаться' }}
          </v-btn>
          <v-btn
            v-if="registered && canCancel"
            color="error"
            variant="tonal"
            :loading="isUnregistering"
            @click="unregister"
          >
            Отказаться
          </v-btn>
          <v-btn
            v-if="registered && new Date(data.start) < new Date()"
            color="info"
            variant="tonal"
            disabled
          >
            Результат
          </v-btn>
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
          <v-col cols="8" sm="6">
            <div class="d-flex align-center mb-2">
              <v-icon icon="mdi-school" class="mr-2"></v-icon>
              <span>{{ data.grade }} класс</span>
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

<script lang="ts">
import { defineComponent } from 'vue'
import ExamsService, { Exam } from '@/api/api.exams'
import { useExamsStore } from '@/stores/exams'

export default defineComponent({
  props: {
    data: {
      type: Object as () => Exam,
      required: true,
    },
  },

  data: () => ({
    isRegistering: false,
    isUnregistering: false,
  }),

  computed: {
    formattedDate(): string {
      const start = new Date(this.data.start)
      return start.toLocaleDateString('ru-RU', {
        day: 'numeric',
        month: 'long',
        year: 'numeric',
      })
    },

    formattedTime(): string {
      const start = new Date(this.data.start)
      const end = new Date(this.data.end)
      return `${start.toLocaleTimeString('ru-RU', {
        hour: '2-digit',
        minute: '2-digit',
      })} - ${end.toLocaleTimeString('ru-RU', {
        hour: '2-digit',
        minute: '2-digit',
      })}`
    },

    capitalizedTitle(): string {
      const title = this.data.type.title
      return title.charAt(0).toUpperCase() + title.slice(1)
    },

    progressColor(): string {
      if (!this.allocation) return 'primary'
      const ratio = this.allocation.occupied / this.allocation.capacity
      if (ratio >= 0.9) return 'error'
      if (ratio >= 0.7) return 'warning'
      return 'success'
    },

    allocation() {
      const examsStore = useExamsStore()
      return examsStore.allocations.get(this.data.ID) || null
    },
    registrationStatus() {
      const examsStore = useExamsStore()
      return (
        examsStore.registrationStatuses.get(this.data.ID) || {
          registered: false,
          registered_to_same_type: false,
        }
      )
    },
    registered() {
      return this.registrationStatus.registered
    },
    registeredSameType() {
      return this.registrationStatus.registered_to_same_type
    },
    canCancel(): boolean {
      return new Date(this.data.start) > new Date()
    },
  },

  methods: {
    async register() {
      const examsStore = useExamsStore()
      try {
        this.isRegistering = true
        await ExamsService.register(this.data.ID)
        await Promise.all([
          examsStore.reloadAll(),
          examsStore.reloadAllExamData(),
        ])
      } finally {
        this.isRegistering = false
      }
    },
    async unregister() {
      const examsStore = useExamsStore()
      try {
        this.isUnregistering = true
        await ExamsService.unregister(this.data.ID)
        await examsStore.reloadAll()
        await examsStore.reloadAllExamData()
      } finally {
        this.isUnregistering = false
      }
    },
  },

  mounted() {
    const examsStore = useExamsStore()
    examsStore.reloadExamData(this.data.ID)
  },
})
</script>

<style scoped>
.registered-border {
  border-left: 6px solid rgb(var(--v-theme-success)) !important;
}
</style>
