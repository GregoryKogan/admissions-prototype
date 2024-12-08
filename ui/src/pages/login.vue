<template>
  <v-container class="fill-height d-flex align-center justify-center">
    <v-card class="pa-5" max-width="500">
      <v-card-title>
        <v-btn icon color="grey" to="/" variant="text">
          <v-icon>mdi-arrow-left</v-icon>
          Назад
        </v-btn>
        <h1 class="text-center flex-grow-1">Вход</h1>
      </v-card-title>
      <v-form @submit.prevent="handleSubmit" ref="form">
        <v-text-field
          v-model="login"
          label="Логин"
          :rules="[rules.required]"
        ></v-text-field>
        <v-text-field
          v-model="password"
          label="Пароль"
          type="password"
          :rules="[rules.required]"
        ></v-text-field>
        <v-btn color="primary" type="submit" class="mt-4 mx-auto d-block"
          >Войти</v-btn
        >
      </v-form>
      <v-btn
        color="secondary"
        to="/register"
        class="mt-4 mx-auto"
        variant="text"
        >Нет аккаунта? Зарегистрироваться</v-btn
      >
    </v-card>
  </v-container>
  <v-snackbar :timeout="5000" v-model="errorSnackbar" color="error">
    {{ errorText }}
  </v-snackbar>
</template>

<script lang="ts">
import { useAuthStore } from '@/stores/auth'
import { AxiosError } from 'axios'
import { defineComponent } from 'vue'
import { VForm } from 'vuetify/components'

export default defineComponent({
  setup() {
    const authStore = useAuthStore()
    return { authStore }
  },
  data() {
    return {
      login: '',
      password: '',
      rules: {
        required: (v: string) => !!v || 'Обязательное поле',
      },
      errorSnackbar: false,
      errorText: '',
    }
  },
  methods: {
    async handleSubmit() {
      const isValid = await (this.$refs.form as VForm).validate()
      if (!isValid.valid) return

      try {
        await this.authStore.login(this.login, this.password)
        if (this.$route.query.redirect) {
          this.$router.push(this.$route.query.redirect as string)
        } else {
          let isAdmin = false
          try {
            const me = await this.authStore.me()
            if (me.role) isAdmin = me.role.admin
          } catch {
            isAdmin = false
          }
          if (isAdmin) {
            this.$router.push('/admin/dashboard')
          } else {
            this.$router.push('/profile')
          }
        }
      } catch (e: unknown) {
        const error = e as AxiosError
        if (
          error.response?.data &&
          typeof error.response.data === 'object' &&
          'message' in error.response.data
        ) {
          this.errorText = (error.response.data as { message: string }).message
        } else {
          this.errorText = 'Ошибка при входе'
          console.error(e)
        }
        this.errorSnackbar = true
      }
    },
  },
})
</script>

<style scoped>
.fill-height {
  height: 100vh;
}
</style>
