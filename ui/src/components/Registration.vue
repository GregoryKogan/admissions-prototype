<template>
  <v-card density="compact">
    <v-card-title class="d-flex flex-sm-row justify-space-between align-center">
      <div class="text-wrap">
        {{ props.data.last_name }} {{ props.data.first_name }}
        {{ props.data.patronymic }}
      </div>
      <div v-if="isPending" class="d-flex flex-column flex-sm-row mt-2 mt-sm-0">
        <v-btn
          density="compact"
          variant="tonal"
          color="success"
          class="mb-2 mb-sm-0 mr-sm-2 ml-sm-2"
          @click="handleApprove"
          >Одобрить</v-btn
        >
        <v-btn
          density="compact"
          variant="tonal"
          color="error"
          @click="handleReject"
          >Отклонить</v-btn
        >
      </div>
    </v-card-title>
    <v-card-subtitle>{{ createdAt }}</v-card-subtitle>
    <v-card-text class="pa-2">
      <v-list density="compact" class="pa-0">
        <v-row dense class="ma-0">
          <v-col cols="12" sm="6" class="pa-0">
            <v-list-item>
              <v-list-item-title>{{ birthDate }}</v-list-item-title>
              <v-list-item-subtitle>Дата рождения</v-list-item-subtitle>
            </v-list-item>
            <v-list-item>
              <v-list-item-title>{{ props.data.grade }}</v-list-item-title>
              <v-list-item-subtitle>Класс поступления</v-list-item-subtitle>
            </v-list-item>
            <v-list-item>
              <v-list-item-title>{{ gender }}</v-list-item-title>
              <v-list-item-subtitle>Пол</v-list-item-subtitle>
            </v-list-item>
          </v-col>
          <v-col cols="12" sm="6" class="pa-0">
            <v-list-item>
              <v-list-item-title>{{
                props.data.parent_phone
              }}</v-list-item-title>
              <v-list-item-subtitle>Телефон родителя</v-list-item-subtitle>
            </v-list-item>
            <v-list-item>
              <v-list-item-title class="text-wrap text-break">
                {{ props.data.parent_last_name }}
                {{ props.data.parent_first_name }}
                {{ props.data.parent_patronymic }}
              </v-list-item-title>
              <v-list-item-subtitle>ФИО родителя</v-list-item-subtitle>
            </v-list-item>
            <v-list-item>
              <template v-slot:title>
                <div class="text-wrap text-break">
                  {{ props.data.old_school }}
                </div>
              </template>
              <v-list-item-subtitle>Предыдущая школа</v-list-item-subtitle>
            </v-list-item>
          </v-col>
        </v-row>
      </v-list>
    </v-card-text>

    <v-dialog v-model="approveDialog" width="auto">
      <v-card>
        <v-card-title>Подтверждение</v-card-title>
        <v-card-text>
          Вы уверены, что хотите одобрить заявку "{{ props.data.last_name }}
          {{ props.data.first_name }}"?
        </v-card-text>
        <v-card-actions>
          <v-btn color="success" @click="confirmApprove">Одобрить</v-btn>
          <v-btn color="grey" @click="approveDialog = false">Отмена</v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>

    <v-dialog v-model="rejectDialog" width="auto">
      <v-card>
        <v-card-title>Подтверждение</v-card-title>
        <v-card-text>
          <p class="mb-4">
            Вы уверены, что хотите отклонить заявку "{{ props.data.last_name }}
            {{ props.data.first_name }}"?
          </p>
          <v-text-field
            v-model="rejectReason"
            label="Причина отклонения"
            variant="outlined"
            density="compact"
          ></v-text-field>
        </v-card-text>
        <v-card-actions>
          <v-btn color="error" @click="confirmReject" :disabled="!rejectReason"
            >Отклонить</v-btn
          >
          <v-btn color="grey" @click="closeRejectDialog">Отмена</v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>
  </v-card>
</template>

<script lang="ts" setup>
import RegistrationService, { Registration } from '@/api/api.registration'
import { ref, withDefaults, defineProps, defineEmits } from 'vue'

const props = withDefaults(
  defineProps<{
    data: Registration
    isPending?: boolean
  }>(),
  {
    isPending: true,
  }
)

const isPending = props.isPending

const createdAt = new Date(props.data.CreatedAt).toLocaleString('ru-RU')

const birthDate = new Date(props.data.birth_date).toLocaleDateString('ru-RU', {
  year: 'numeric',
  month: 'long',
  day: 'numeric',
})

const gender =
  props.data.gender === 'M'
    ? 'Мужской'
    : props.data.gender === 'F'
    ? 'Женский'
    : 'Не указан'

const approveDialog = ref(false)
const rejectDialog = ref(false)
const rejectReason = ref('')

const handleApprove = () => {
  approveDialog.value = true
}

const handleReject = () => {
  rejectReason.value = ''
  rejectDialog.value = true
}

const closeRejectDialog = () => {
  rejectDialog.value = false
  rejectReason.value = ''
}

const emit = defineEmits(['statusChanged'])

const confirmApprove = async () => {
  approveDialog.value = false
  await RegistrationService.accept(props.data.ID)
  emit('statusChanged')
}

const confirmReject = async () => {
  await RegistrationService.reject(props.data.ID, rejectReason.value)
  closeRejectDialog()
  emit('statusChanged')
}
</script>
