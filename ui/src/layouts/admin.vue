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
          to="/admin/profile"
          prepend-icon="mdi-account"
          title="Профиль"
        />
        <v-list-item
          to="/admin/registrations"
          prepend-icon="mdi-format-list-bulleted"
          title="Регистрации"
        />
      </v-list>
    </v-navigation-drawer>

    <v-app-bar v-if="mobile" position="fixed">
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
  </v-app>
</template>

<script lang="ts" setup>
import { useRouter } from 'vue-router'
import { ref } from 'vue'
import { useDisplay } from 'vuetify'

const drawer = ref(false)
const rail = ref(false)
const { mobile } = useDisplay()
const router = useRouter()

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
