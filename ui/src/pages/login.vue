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
</template>

<script lang="ts">
import { defineComponent } from 'vue'
import { VForm } from 'vuetify/components'

export default defineComponent({
  data() {
    return {
      login: '',
      password: '',
      rules: {
        required: (v: string) => !!v || 'Обязательное поле',
      },
    }
  },
  methods: {
    async handleSubmit() {
      const isValid = await (this.$refs.form as VForm).validate()
      if (!isValid.valid) return

      console.log({
        email: this.login,
        password: this.password,
      })
    },
  },
})
</script>

<style scoped>
.fill-height {
  height: 100vh;
}
</style>
