<template>
  <v-container class="pa-2 pa-sm-4">
    <v-row justify="center">
      <v-col cols="12" sm="11" md="10" lg="8">
        <v-card v-if="registration" class="elevation-3">
          <v-card-title class="text-h5 text-sm-h4 pa-3 pa-sm-4 text-wrap">
            {{ registration.last_name }} {{ registration.first_name }}
            {{ registration.patronymic }}
          </v-card-title>

          <div class="px-3 px-sm-4 pb-4 d-flex flex-column align-center">
            <v-avatar
              :size="$vuetify.display.smAndDown ? 100 : 140"
              class="mb-3 elevation-2"
            >
              <v-icon
                :size="$vuetify.display.smAndDown ? 80 : 112"
                color="grey-darken-2"
                disabled
              >
                mdi-account-circle
              </v-icon>
            </v-avatar>
          </div>

          <v-card-subtitle class="px-3 px-sm-4 pb-0">
            Регистрация: {{ createdAt }}
          </v-card-subtitle>

          <v-card-text>
            <v-row dense>
              <v-col cols="12" sm="6">
                <v-list
                  :density="
                    $vuetify.display.smAndDown ? 'compact' : 'comfortable'
                  "
                >
                  <v-list-item>
                    <v-list-item-title class="text-wrap">{{
                      registration.email
                    }}</v-list-item-title>
                    <v-list-item-subtitle>Email</v-list-item-subtitle>
                  </v-list-item>
                  <v-list-item>
                    <v-list-item-title class="text-wrap">{{
                      formatBirthDate
                    }}</v-list-item-title>
                    <v-list-item-subtitle>Дата рождения</v-list-item-subtitle>
                  </v-list-item>
                  <v-list-item>
                    <v-list-item-title class="text-wrap">{{
                      registration.grade
                    }}</v-list-item-title>
                    <v-list-item-subtitle
                      >Класс поступления</v-list-item-subtitle
                    >
                  </v-list-item>
                  <v-list-item>
                    <v-list-item-title class="text-wrap">{{
                      formatGender
                    }}</v-list-item-title>
                    <v-list-item-subtitle>Пол</v-list-item-subtitle>
                  </v-list-item>
                </v-list>
              </v-col>
              <v-col cols="12" sm="6">
                <v-list
                  :density="
                    $vuetify.display.smAndDown ? 'compact' : 'comfortable'
                  "
                >
                  <v-list-item>
                    <v-list-item-title class="text-wrap">{{
                      registration.parent_phone
                    }}</v-list-item-title>
                    <v-list-item-subtitle
                      >Телефон родителя</v-list-item-subtitle
                    >
                  </v-list-item>
                  <v-list-item>
                    <v-list-item-title class="text-wrap">
                      {{ registration.parent_last_name }}
                      {{ registration.parent_first_name }}
                      {{ registration.parent_patronymic }}
                    </v-list-item-title>
                    <v-list-item-subtitle>Родитель</v-list-item-subtitle>
                  </v-list-item>
                  <v-list-item>
                    <v-list-item-title class="text-wrap">{{
                      registration.old_school
                    }}</v-list-item-title>
                    <v-list-item-subtitle
                      >Предыдущая школа</v-list-item-subtitle
                    >
                  </v-list-item>
                </v-list>
              </v-col>
            </v-row>

            <div class="px-3 px-sm-4 mb-4 mb-sm-6">
              <div class="text-subtitle-1 mb-2">Пометки</div>
              <div class="d-flex flex-wrap gap-2">
                <v-chip
                  :color="registration.june_exam ? 'success' : 'grey'"
                  class="mr-2"
                >
                  Экзамен в июне
                </v-chip>
                <v-chip :color="registration.vmsh ? 'success' : 'grey'">
                  ВМШ
                </v-chip>
              </div>
            </div>

            <v-divider class="mb-4 mb-sm-6"></v-divider>

            <div class="px-3 px-sm-4 mb-4">
              <div class="text-subtitle-1 mb-2">Настройки</div>
              <div class="d-flex flex-wrap align-center gap-4">
                <span class="d-flex align-center">
                  <span class="text-body-1">Тема оформления</span>
                  <v-btn
                    :icon="
                      theme.global.current.value.dark
                        ? 'mdi-weather-sunny'
                        : 'mdi-weather-night'
                    "
                    @click="toggleTheme"
                    variant="text"
                    class="ml-2"
                    size="small"
                  ></v-btn>
                </span>
              </div>
            </div>

            <div class="px-3 px-sm-4">
              <div class="text-subtitle-1 mb-2">Действия</div>
              <v-btn
                color="error"
                @click="handleLogout"
                variant="outlined"
                :size="$vuetify.display.smAndDown ? 'small' : 'default'"
              >
                Выйти из аккаунта
              </v-btn>
            </div>
          </v-card-text>
        </v-card>

        <v-skeleton-loader v-else type="card" class="mt-4"></v-skeleton-loader>
      </v-col>
    </v-row>

    <v-dialog
      v-model="logoutDialog"
      :width="$vuetify.display.smAndDown ? '90%' : 'auto'"
    >
      <v-card>
        <v-card-title class="text-body-1 text-sm-h6"
          >Подтверждение</v-card-title
        >
        <v-card-text>Вы уверены, что хотите выйти?</v-card-text>
        <v-card-actions class="gap-2">
          <v-btn
            color="error"
            @click="confirmLogout"
            :size="$vuetify.display.smAndDown ? 'small' : 'default'"
          >
            Выйти
          </v-btn>
          <v-btn
            color="grey"
            @click="logoutDialog = false"
            :size="$vuetify.display.smAndDown ? 'small' : 'default'"
          >
            Отмена
          </v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>
  </v-container>
</template>

<script lang="ts" setup>
import RegistrationService, { Registration } from '@/api/api.registration'
import { useAuthStore } from '@/stores/auth'
import { onMounted, ref, computed } from 'vue'
import { useRouter } from 'vue-router'
import { useTheme } from 'vuetify'

const authStore = useAuthStore()
const createdAt = ref('')
const registration = ref<Registration | null>(null)
const router = useRouter()
const logoutDialog = ref(false)
const theme = useTheme()

const formatBirthDate = computed(() => {
  if (!registration.value) return ''
  return new Date(registration.value.birth_date).toLocaleDateString('ru-RU', {
    year: 'numeric',
    month: 'long',
    day: 'numeric',
  })
})

const formatGender = computed(() => {
  if (!registration.value) return ''
  return registration.value.gender === 'M'
    ? 'Мужской'
    : registration.value.gender === 'F'
    ? 'Женский'
    : 'Не указан'
})

const handleLogout = () => {
  logoutDialog.value = true
}

const confirmLogout = async () => {
  await authStore.logout()
  logoutDialog.value = false
  router.push('/')
}

const toggleTheme = () => {
  theme.global.name.value = theme.global.current.value.dark ? 'light' : 'dark'
  localStorage.setItem('theme', theme.global.name.value)
}

onMounted(async () => {
  const me = await authStore.me()
  createdAt.value = new Date(me.CreatedAt).toLocaleString('ru-RU')
  registration.value = (await RegistrationService.mine()).data
})
</script>

<route lang="yaml">
meta:
  layout: user
</route>
