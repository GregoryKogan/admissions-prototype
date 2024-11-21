<template>
  <v-container class="fill-height d-flex align-center justify-center">
    <v-card class="pa-5" max-width="500">
      <v-card-title>
        <v-btn icon color="grey" to="/" variant="text">
          <v-icon>mdi-arrow-left</v-icon>
          Назад
        </v-btn>
        <h1 class="text-center flex-grow-1">Регистрация</h1>
      </v-card-title>
      <v-form @submit.prevent="handleSubmit">
        <v-text-field
          v-model="username"
          label="Имя пользователя"
          required
        ></v-text-field>
        <v-text-field
          v-model="email"
          label="Электронная почта"
          required
        ></v-text-field>
        <v-text-field
          v-model="password"
          label="Пароль"
          type="password"
          required
        ></v-text-field>
        <v-btn color="primary" type="submit" class="mt-4 mx-auto d-block"
          >Зарегистрироваться</v-btn
        >
      </v-form>
      <v-btn
        color="secondary"
        to="/login"
        class="mt-4 mx-auto d-block"
        variant="text"
        >Уже есть аккаунт? Войти</v-btn
      >
    </v-card>
  </v-container>
</template>

<script lang="ts">
import { defineComponent } from 'vue'

export default defineComponent({
  data() {
    return {
      username: '',
      email: '',
      password: '',
    }
  },
  methods: {
    handleSubmit() {
      console.log({
        username: this.username,
        email: this.email,
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
.top-left {
  position: absolute;
  top: 16px;
  left: 16px;
}
</style>
