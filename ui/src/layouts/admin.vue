<template>
  <v-app style="height: 100vh">
    <v-navigation-drawer
      v-model="drawer"
      :rail="!mobile && rail"
      @click="rail = false"
      :temporary="mobile"
      :permanent="!mobile"
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

    <v-app-bar v-if="mobile">
      <v-app-bar-nav-icon @click="drawer = !drawer"></v-app-bar-nav-icon>
      <v-toolbar-title>Админ панель</v-toolbar-title>
    </v-app-bar>

    <v-main
      @click="!mobile && (rail = true)"
      style="overflow: auto; max-height: 100vh"
      :class="{ 'pt-15': mobile }"
    >
      <router-view />
    </v-main>

    <v-dialog v-model="logoutDialog" :width="mobile ? '90%' : 'auto'">
      <v-card>
        <v-card-title>Подтверждение</v-card-title>
        <v-card-text> Вы уверены, что хотите выйти? </v-card-text>
        <v-card-actions>
          <v-btn color="error" @click="confirmLogout">Выйти</v-btn>
          <v-btn color="grey" @click="logoutDialog = false">Отмена</v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>
  </v-app>
</template>

<script lang="ts" setup>
import { useAuthStore } from '@/stores/auth'
import { useRouter } from 'vue-router'
import { ref } from 'vue'
import { useDisplay } from 'vuetify'

const drawer = ref(false)
const rail = ref(false)
const { mobile } = useDisplay()

const auth = useAuthStore()
const router = useRouter()

const logoutDialog = ref(false)

const handleLogout = () => {
  logoutDialog.value = true
}

const confirmLogout = async () => {
  await auth.logout()
  logoutDialog.value = false
  router.push('/')
}

// Close drawer when route changes on mobile
router.afterEach(() => {
  if (mobile.value) {
    drawer.value = false
  }
})
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
