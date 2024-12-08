<template>
  <v-app style="height: 100vh">
    <v-navigation-drawer
      v-model="drawer"
      :rail="rail"
      @click="rail = false"
      permanent
    >
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
        <v-list-item
          prepend-icon="mdi-logout"
          title="Выход"
          @click="handleLogout"
        >
        </v-list-item>
      </v-list>
    </v-navigation-drawer>

    <v-main @click="rail = true" style="overflow: auto; max-height: 100vh">
      <router-view />
    </v-main>
  </v-app>
</template>

<script lang="ts" setup>
import { useAuthStore } from '@/stores/auth'
import { useRouter } from 'vue-router'
import { ref } from 'vue'

const drawer = ref(true)
const rail = ref(false)

const auth = useAuthStore()
const router = useRouter()

const handleLogout = async () => {
  await auth.logout()
  router.push('/')
}
</script>

<style>
.v-navigation-drawer__content:not(:hover)::-webkit-scrollbar {
  display: none;
}

.v-navigation-drawer__content:not(:hover) {
  -ms-overflow-style: none;
  scrollbar-width: none;
}

.v-navigation-drawer__content:hover {
  scrollbar-width: thin;
}
</style>
