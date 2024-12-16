<template>
  <v-app>
    <v-navigation-drawer v-model="drawer" temporary>
      <v-list>
        <v-list-item to="/profile" prepend-icon="mdi-account"
          >Профиль</v-list-item
        >
        <v-list-item to="/chat" prepend-icon="mdi-chat">Поддержка</v-list-item>
        <v-list-item to="/exams" prepend-icon="mdi-file-document"
          >Экзамены</v-list-item
        >
      </v-list>
    </v-navigation-drawer>

    <v-app-bar style="position: fixed">
      <v-app-bar-nav-icon
        class="d-md-none"
        @click="drawer = !drawer"
      ></v-app-bar-nav-icon>
      <div class="d-flex align-center">
        <v-img
          src="@/assets/L2SH-logo.png"
          width="48"
          height="48"
          class="me-2 ml-4"
          alt="L2SH Logo"
          @click="goToIndex"
        ></v-img>
        <div class="d-none d-md-flex">
          <v-btn to="/exams" prepend-icon="mdi-file-document" class="mx-2"
            >Экзамены</v-btn
          >
          <v-btn to="/chat" prepend-icon="mdi-chat" class="mx-2"
            >Поддержка</v-btn
          >
        </div>
      </div>
      <v-spacer></v-spacer>
      <div class="d-none d-md-flex">
        <v-btn to="/profile" prepend-icon="mdi-account" class="ms-2 mr-4"
          >Профиль</v-btn
        >
      </div>
    </v-app-bar>

    <v-main>
      <router-view />
    </v-main>
  </v-app>
</template>

<script lang="ts" setup>
import { useRouter } from 'vue-router'
import { ref } from 'vue'

const router = useRouter()
const drawer = ref(false)

const goToIndex = () => {
  router.push('/')
}
</script>
