<template>
  <v-container
    class="text-center fill-height d-flex align-center justify-center"
  >
    <v-row>
      <v-col cols="12">
        <img
          src="@/assets/L2SH-logo.png"
          alt="L2SH.Admissions Logo"
          class="responsive-logo"
        />
      </v-col>
      <v-col cols="12">
        <h1>Лицей "Вторая школа"</h1>
        <h2>Приемная комиссия</h2>
      </v-col>
      <v-col cols="12" class="d-flex justify-center">
        <v-btn color="primary" @click="enter" class="mx-auto sized-button"
          >Войти</v-btn
        >
      </v-col>
    </v-row>
  </v-container>
</template>

<script lang="ts">
import { defineComponent } from 'vue'
import { useAuthStore } from '@/stores/auth'

export default defineComponent({
  methods: {
    async enter() {
      const authStore = useAuthStore()

      await authStore.checkAuth()
      if (authStore.isAuth) {
        try {
          const me = await authStore.me()
          if (me?.role?.admin) {
            this.$router.push('/admin/profile')
          } else {
            this.$router.push('/profile')
          }
        } catch {
          this.$router.push('/profile')
        }
      } else {
        this.$router.push('/login')
      }
    },
  },
})
</script>

<style scoped>
.responsive-logo {
  max-width: 150px;
  width: 100%;
  height: auto;
}

.sized-button {
  max-width: 200px;
  width: 100%;
}
</style>
