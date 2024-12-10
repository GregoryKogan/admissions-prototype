<template>
  <v-app style="height: 100vh">
    <v-navigation-drawer
      v-model="drawer"
      :rail="!smAndDown && rail"
      @click="rail = false"
      :temporary="smAndDown"
      :permanent="!smAndDown"
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

    <v-app-bar v-if="smAndDown" position="fixed">
      <v-app-bar-nav-icon @click="drawer = !drawer"></v-app-bar-nav-icon>
      <v-toolbar-title>Админ панель</v-toolbar-title>
    </v-app-bar>

    <v-main
      @click="!smAndDown && (rail = true)"
      style="overflow: auto; max-height: 100vh"
      :class="{ 'pt-15': smAndDown }"
    >
      <router-view />
    </v-main>
  </v-app>
</template>

<script lang="ts" setup>
import { useRouter } from 'vue-router'
import { ref } from 'vue'
import { useDisplay } from 'vuetify'

const { smAndDown } = useDisplay()
const drawer = ref(!smAndDown.value)
const rail = ref(false)
const router = useRouter()

// Close drawer when route changes on mobile
router.afterEach(() => {
  if (smAndDown.value) {
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
