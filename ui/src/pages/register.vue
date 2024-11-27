<template>
  <v-container class="fill-height d-flex align-center justify-center">
    <v-card v-if="!finished" class="pa-5" width="700">
      <v-card-title>
        <v-btn icon color="grey" to="/" variant="text">
          <v-icon>mdi-arrow-left</v-icon>
          Назад
        </v-btn>
        <h1 class="text-center flex-grow-1">Регистрация</h1>
      </v-card-title>
      <v-form @submit.prevent="handleSubmit" ref="form">
        <v-text-field
          v-model="email"
          label="Электронная почта"
          :rules="[rules.required, rules.email]"
        ></v-text-field>
        <v-text-field
          v-model="first_name"
          label="Имя"
          :rules="[rules.required]"
        ></v-text-field>
        <v-text-field
          v-model="last_name"
          label="Фамилия"
          :rules="[rules.required]"
        ></v-text-field>
        <v-text-field v-model="patronymic" label="Отчество"></v-text-field>
        <v-select
          v-model="gender"
          :items="['Мужской', 'Женский']"
          label="Пол"
          :rules="[rules.required]"
        ></v-select>
        <v-date-picker
          v-model="birth_date"
          title="Дата рождения"
          :rules="[rules.required]"
        ></v-date-picker>
        <v-select
          v-model="grade"
          :items="[6, 7, 8, 9, 10, 11]"
          label="Класс поступления"
          :rules="[rules.required]"
        ></v-select>
        <v-text-field
          v-model="old_school"
          label="Предыдущая школа"
          :rules="[rules.required]"
        ></v-text-field>
        <v-text-field
          v-model="parent_first_name"
          label="Имя родителя"
          :rules="[rules.required]"
        ></v-text-field>
        <v-text-field
          v-model="parent_last_name"
          label="Фамилия родителя"
          :rules="[rules.required]"
        ></v-text-field>
        <v-text-field
          v-model="parent_patronymic"
          label="Отчество родителя"
        ></v-text-field>
        <v-text-field
          v-model="parent_phone"
          label="Телефон родителя"
          :rules="[rules.required, rules.phone]"
        ></v-text-field>
        <v-checkbox
          v-model="june_exam"
          label="Буду сдавать экзамен в июне"
        ></v-checkbox>
        <v-checkbox v-model="vmsh" label="Учился в ВМШ"></v-checkbox>
        <v-textarea
          v-model="source"
          label="Откуда узнали о лицее?"
        ></v-textarea>
        <v-btn color="primary" type="submit" class="mt-4 mx-auto d-block"
          >Зарегистрироваться</v-btn
        >
      </v-form>
      <div class="d-flex justify-end">
        <v-btn color="secondary" to="/login" class="mt-4" variant="text"
          >Уже есть аккаунт? Войти</v-btn
        >
      </div>
    </v-card>
    <v-card v-else elevation="0">
      <v-card-title>Вы успешно заполнили анкету</v-card-title>
      <v-card-text>
        <v-alert
          text="Для завершения регистрации подтвердите почту, перейдя по ссылке из письма"
          variant="tonal"
          type="success"
        ></v-alert>
      </v-card-text>
    </v-card>
  </v-container>
  <v-snackbar :timeout="5000" v-model="errorSnackbar" color="error">
    {{ errorText }}
  </v-snackbar>
</template>

<script lang="ts">
import { defineComponent } from 'vue'
import { VForm } from 'vuetify/components'

export default defineComponent({
  data() {
    return {
      email: '',
      first_name: '',
      last_name: '',
      patronymic: '',
      gender: '', // 'Мужской' or 'Женский'
      birth_date: new Date(),
      grade: 6, // 6, 7, 8, 9, 10, 11
      old_school: '',
      parent_first_name: '',
      parent_last_name: '',
      parent_patronymic: '',
      parent_phone: '',
      june_exam: false,
      vmsh: false,
      source: '',
      errorSnackbar: false,
      errorText: '',
      finished: false,
      rules: {
        required: (value: string) => !!value || 'Обязательное поле.',
        email: (value: string) =>
          /.+@.+\..+/.test(value) || 'Некорректный email.',
        phone: (value: string) =>
          /^((8|\+7)[- ]?)?(\(?\d{3}\)?[- ]?)?[\d\- ]{7,10}$/.test(value) ||
          'Некорректный номер телефона.',
      },
    }
  },
  computed: {
    formattedGender(): string | null {
      if (this.gender === 'Мужской') return 'M'
      if (this.gender === 'Женский') return 'F'
      return null
    },
  },
  methods: {
    async handleSubmit() {
      const isValid = await (this.$refs.form as VForm).validate()
      if (!isValid.valid) return

      const response = await fetch('/api/regdata', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
          email: this.email,
          first_name: this.first_name,
          last_name: this.last_name,
          patronymic: this.patronymic,
          gender: this.formattedGender,
          birth_date: this.birth_date,
          grade: this.grade,
          old_school: this.old_school,
          parent_first_name: this.parent_first_name,
          parent_last_name: this.parent_last_name,
          parent_patronymic: this.parent_patronymic,
          parent_phone: this.parent_phone,
          june_exam: this.june_exam,
          vmsh: this.vmsh,
          source: this.source,
        }),
      })

      if (response.ok) {
        this.finished = true
      } else {
        const responseText = await response.text()
        const responseBody = JSON.parse(responseText)
        this.errorText = responseBody.message
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
.top-left {
  position: absolute;
  top: 16px;
  left: 16px;
}
</style>
