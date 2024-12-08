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
        <v-btn color="primary" to="/login" class="mx-auto sized-button"
          >Войти</v-btn
        >
      </v-col>
    </v-row>
  </v-container>
</template>

<script lang="ts">
import { defineComponent } from 'vue'
import { useAuthStore } from '@/stores/auth'
import { useRouter } from 'vue-router'

export default defineComponent({
  async mounted() {
    const authStore = useAuthStore()
    const router = useRouter()

    await authStore.checkAuth()
    if (authStore.isAuth) {
      try {
        const me = await authStore.me()
        if (me?.role?.admin) {
          router.push('/admin/dashboard')
        } else {
          router.push('/home')
        }
      } catch {
        router.push('/home')
      }
    }
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

<route lang="yaml">
meta:
  layout: public
</route>
