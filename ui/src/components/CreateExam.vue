<template>
  <v-form @submit.prevent="handleSubmit" ref="form">
    <v-row>
      <v-col cols="12" md="4">
        <v-text-field
          v-model="formattedSelectedDate"
          label="Дата проведения"
          prepend-icon="mdi-calendar"
          :rules="[(v) => !!selectedDate || 'Выберите дату']"
          :active="datePickerDialog"
          readonly
        >
          <v-dialog v-model="datePickerDialog" activator="parent" width="auto">
            <v-date-picker
              v-if="datePickerDialog"
              v-model="selectedDate"
              title="Дата проведения"
              required
            ></v-date-picker>
          </v-dialog>
        </v-text-field>
      </v-col>
      <v-col cols="12" md="4">
        <v-text-field
          v-model="start"
          label="Начало экзамена"
          prepend-icon="mdi-clock-time-ten-outline"
          placeholder="ЧЧ:ММ"
          clearable
          :rules="[
            (v) => !!v || 'Выберите время начала',
            (v) =>
              /^([01]?[0-9]|2[0-3]):[0-5][0-9]$/.test(v) ||
              'Неверный формат времени (ЧЧ:ММ)',
            (v) =>
              isStartTimeValid ||
              'Время начала должно быть раньше времени окончания',
          ]"
        >
        </v-text-field>
      </v-col>
      <v-col cols="12" md="4">
        <v-text-field
          v-model="end"
          label="Конец экзамена"
          prepend-icon="mdi-clock-time-three-outline"
          placeholder="ЧЧ:ММ"
          clearable
          :rules="[
            (v) => !!v || 'Выберите время окончания',
            (v) =>
              /^([01]?[0-9]|2[0-3]):[0-5][0-9]$/.test(v) ||
              'Неверный формат времени (ЧЧ:ММ)',
            (v) =>
              isEndTimeValid ||
              'Время окончания должно быть позже времени начала',
          ]"
        >
        </v-text-field>
      </v-col>
    </v-row>
    <v-row>
      <v-col cols="12" :md="showCustomLocation ? 6 : 12">
        <v-select
          v-model="selectedPredefinedLocation"
          :items="predefinedLocations"
          label="Место проведения"
          :rules="[(v) => !!v || 'Выберите место проведения']"
          required
        />
      </v-col>
      <v-col v-if="showCustomLocation" cols="12" md="6">
        <v-text-field
          v-model="form.location"
          label="Введите место проведения"
          :rules="[(v) => !!v || 'Введите место проведения']"
          required
        />
      </v-col>
    </v-row>
    <v-text-field
      v-model.number="form.capacity"
      label="Количество мест"
      type="number"
      min="1"
      max="999"
      :rules="[
        (v) => !!v || 'Введите количество мест',
        (v) => v > 0 || 'Количество мест должно быть больше 0',
        (v) => v <= 999 || 'Максимальное количество мест - 999',
        (v) =>
          Number.isInteger(v) || 'Количество мест должно быть целым числом',
      ]"
      required
    />
    <v-select
      v-model="form.grade"
      :items="grades"
      label="Класс"
      :rules="[(v) => !!v || 'Выберите класс']"
      required
    />
    <v-select
      v-model="form.type_id"
      :items="examTypes"
      item-text="title"
      item-value="ID"
      label="Тип экзамена"
      :rules="[(v) => !!v || 'Выберите тип экзамена']"
      required
    />
    <v-btn type="submit" color="primary"> Создать экзамен </v-btn>
  </v-form>
</template>

<script lang="ts">
import ExamsService, { ExamRequest, ExamType } from '@/api/api.exams'
import { defineComponent } from 'vue'
import { VForm } from 'vuetify/components'

export default defineComponent({
  data: () => ({
    datePickerDialog: false,
    startPickerDialog: false,
    endPickerDialog: false,
    start: '10:00',
    end: '12:00',
    selectedPredefinedLocation: null as string | null,
    predefinedLocations: ['Новое здание', 'Старое здание', 'Другое'],
    form: {
      location: '',
      capacity: 30,
      grade: 6,
      type_id: 0,
    },
    selectedDate: new Date(),
    grades: [6, 7, 8, 9, 10, 11],
    examTypes: [] as ExamType[],
  }),
  async mounted() {
    const response = await ExamsService.types()
    this.examTypes = response.data
    this.form.type_id = this.examTypes[0].ID
  },
  methods: {
    async handleSubmit() {
      const isValid = await (this.$refs.form as VForm).validate()
      if (!isValid.valid) return

      const examRequest: ExamRequest = {
        location: this.form.location,
        capacity: this.form.capacity,
        grade: this.form.grade,
        type_id: this.form.type_id,
        start: this.startTimestamp,
        end: this.endTimestamp,
      }

      await ExamsService.create(examRequest)
      this.$emit('exam-created')
    },
    updateLocation(value: string | null) {
      if (!value || value === 'Другое') {
        this.form.location = ''
      } else {
        this.form.location = value
      }
    },
  },
  computed: {
    formattedSelectedDate() {
      return this.selectedDate.toLocaleDateString('ru-RU')
    },
    startTimestamp(): Date {
      const [hours, minutes] = this.start.split(':')
      const date = new Date(this.selectedDate)
      date.setHours(parseInt(hours), parseInt(minutes), 0)
      return date
    },
    endTimestamp(): Date {
      const [hours, minutes] = this.end.split(':')
      const date = new Date(this.selectedDate)
      date.setHours(parseInt(hours), parseInt(minutes), 0)
      return date
    },
    showCustomLocation(): boolean {
      return this.selectedPredefinedLocation === 'Другое'
    },
    isEndTimeValid(): boolean {
      if (!this.start || !this.end) return true
      const [startHours, startMinutes] = this.start.split(':').map(Number)
      const [endHours, endMinutes] = this.end.split(':').map(Number)

      if (endHours > startHours) return true
      if (endHours === startHours) return endMinutes > startMinutes
      return false
    },
    isStartTimeValid(): boolean {
      if (!this.start || !this.end) return true
      const [startHours, startMinutes] = this.start.split(':').map(Number)
      const [endHours, endMinutes] = this.end.split(':').map(Number)

      if (startHours < endHours) return true
      if (startHours === endHours) return startMinutes < endMinutes
      return false
    },
  },
  watch: {
    selectedPredefinedLocation(newVal) {
      this.updateLocation(newVal)
    },
  },
})
</script>
