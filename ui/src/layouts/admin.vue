<template>
  <v-app>
    <v-navigation-drawer permanent>
      <v-list>
        <v-list-item
          to="/admin/dashboard"
          prepend-icon="mdi-view-dashboard"
          title="Панель"
        />
        <v-list-item
          to="/admin/registrations"
          prepend-icon="mdi-format-list-bulleted"
          title="Регистрации"
        />
      </v-list>
    </v-navigation-drawer>

    <v-app-bar>
      <v-app-bar-title>Панель администратора</v-app-bar-title>
      <v-spacer></v-spacer>
      <v-btn icon="mdi-logout" variant="text" @click="handleLogout"></v-btn>
    </v-app-bar>

    <v-main>
      <router-view />
    </v-main>
  </v-app>
</template>

<script lang="ts" setup>
import { useAuthStore } from '@/stores/auth'
import { useRouter } from 'vue-router'

const auth = useAuthStore()
const router = useRouter()

const handleLogout = async () => {
  await auth.logout()
  router.push('/')
}
</script>
