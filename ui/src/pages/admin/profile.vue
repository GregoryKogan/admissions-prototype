<template>
  <v-container class="pa-2 pa-sm-4">
    <v-row justify="center">
      <v-col cols="12" sm="11" md="10" lg="8">
        <v-card v-if="admin" class="elevation-3">
          <v-card-title class="text-h5 text-sm-h4 pa-3 pa-sm-4">
            Администратор
          </v-card-title>
          <v-card-subtitle class="px-3 px-sm-4 pb-0">
            Регистрация: {{ createdAt }}
          </v-card-subtitle>

          <v-card-text>
            <v-list
              :density="$vuetify.display.smAndDown ? 'compact' : 'comfortable'"
            >
              <v-list-item>
                <v-list-item-title>{{ admin.login }}</v-list-item-title>
                <v-list-item-subtitle>Login</v-list-item-subtitle>
              </v-list-item>
            </v-list>

            <v-divider class="my-4 my-sm-6"></v-divider>

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
import { useAuthStore } from '@/stores/auth'
import { onMounted, ref } from 'vue'
import { useRouter } from 'vue-router'
import { useTheme } from 'vuetify'

const authStore = useAuthStore()
const router = useRouter()
const theme = useTheme()
interface Admin {
  login: string
  CreatedAt: string
}

const admin = ref<Admin | null>(null)
const createdAt = ref('')
const logoutDialog = ref(false)

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
  admin.value = me
  createdAt.value = new Date(me.CreatedAt).toLocaleString('ru-RU')
})
</script>

<route lang="yaml">
meta:
  layout: admin
</route>
