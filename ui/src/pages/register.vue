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
          v-model="last_name"
          label="Фамилия поступающего"
          :rules="[rules.required]"
        ></v-text-field>
        <v-text-field
          v-model="first_name"
          label="Имя поступающего"
          :rules="[rules.required]"
        ></v-text-field>
        <v-text-field
          v-model="patronymic"
          label="Отчество поступающего"
        ></v-text-field>
        <v-select
          v-model="gender"
          :items="['Мужской', 'Женский']"
          label="Пол"
          :rules="[rules.required]"
        ></v-select>
        <v-date-picker
          v-model="birth_date"
          title="Дата рождения поступающего"
          width="70%"
          min-width="300px"
          :rules="[rules.required]"
          :max="maxDate"
          :min="minDate"
        ></v-date-picker>
        <v-select
          v-model="grade"
          :items="[6, 7, 8, 9, 10, 11]"
          label="В какой класс вы хотите поступить?"
          :rules="[rules.required]"
        ></v-select>
        <v-text-field
          v-model="old_school"
          label="В какой школе вы учитесь сейчас?"
          :rules="[rules.required]"
        ></v-text-field>
        <v-text-field
          v-model="parent_last_name"
          label="Фамилия родителя"
          :rules="[rules.required]"
        ></v-text-field>
        <v-text-field
          v-model="parent_first_name"
          label="Имя родителя"
          :rules="[rules.required]"
        ></v-text-field>
        <v-text-field
          v-model="parent_patronymic"
          label="Отчество родителя"
        ></v-text-field>
        <v-text-field
          v-model="parent_phone"
          label="Телефон родителя"
          hint="Формат: +79160000000"
          :rules="[rules.required, rules.phone]"
        ></v-text-field>
        <span>
          Отметьте этот пункт, если поступающий проживает не в Москве или
          Московской области и хочет сдавать экзамены в июне.
        </span>
        <v-checkbox v-model="june_exam" label="Экзамены в июне"></v-checkbox>
        <span>
          Посещали ли вы занятия Вечерней Математической Школы Л2Ш в текущем
          учебном году?
        </span>
        <v-checkbox v-model="vmsh" label="Да"></v-checkbox>
        <span>Ответьте, пожалуйста, как Вы узнали о Лицее "Вторая школа"?</span>
        <v-select
          v-model="selectedPredefinedSource"
          label="Источник"
          :items="predefinedSources"
          :rules="[(v) => !!v || 'Выберите источник']"
          required
          class="mt-4"
        />
        <v-text-field
          v-if="showCustomSource"
          v-model="source"
          label="Введите источник"
          :rules="[(v) => !!v || 'Введите источник']"
          required
        />
        <v-checkbox
          :rules="[(v) => !!v || 'Необходимо ознакомиться с положением']"
        >
          <template v-slot:label>
            Ознакомлен с "Положением о приеме в ГБОУ Лицей "Вторая школа"
          </template>
        </v-checkbox>
        <v-btn
          color="primary"
          type="submit"
          class="mt-4 mx-auto d-block"
          :loading="registering"
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
          text="Для завершения регистрации подтвердите почту, перейдя по ссылке из письма. Проверьте папку «Спам», если письмо не пришло."
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
import RegistrationService from '@/api/api.registration'
import { AxiosError } from 'axios'
import { defineComponent } from 'vue'
import { VForm } from 'vuetify/components'

export default defineComponent({
  data() {
    const defaultDate = new Date()
    defaultDate.setFullYear(defaultDate.getFullYear() - 14)

    return {
      email: '',
      first_name: '',
      last_name: '',
      patronymic: '',
      gender: '', // 'Мужской' or 'Женский'
      birth_date: defaultDate,
      grade: 6, // 6, 7, 8, 9, 10, 11
      old_school: '',
      parent_first_name: '',
      parent_last_name: '',
      parent_patronymic: '',
      parent_phone: '',
      june_exam: false,
      vmsh: false,
      source: '',
      selectedPredefinedSource: null as string | null,
      predefinedSources: [
        'Здесь учился (учится) кто-то из моей семьи',
        'Нашел в рейтинге школ Москвы',
        'Порекомендовали друзья/знакомые',
        'Узнал из социальных сетей',
        'Другое',
      ],
      errorSnackbar: false,
      errorText: '',
      registering: false,
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
    formattedGender(): string {
      if (this.gender === 'Мужской') return 'M'
      if (this.gender === 'Женский') return 'F'
      return 'N'
    },
    maxDate(): Date {
      const date = new Date()
      date.setFullYear(date.getFullYear() - 9)
      return date
    },
    minDate(): Date {
      const date = new Date()
      date.setFullYear(date.getFullYear() - 18)
      return date
    },
    showCustomSource(): boolean {
      return this.selectedPredefinedSource === 'Другое'
    },
  },
  methods: {
    async handleSubmit() {
      const isValid = await (this.$refs.form as VForm).validate()
      if (!isValid.valid) return

      try {
        this.registering = true
        await RegistrationService.register({
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
        })
        this.finished = true
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
          console.warn(e)
        }
        this.errorSnackbar = true
      } finally {
        this.registering = false
      }
    },
    updateSource(value: string | null) {
      if (!value || value === 'Другое') {
        this.source = ''
      } else {
        this.source = value
      }
    },
  },
  watch: {
    selectedPredefinedSource(newVal) {
      this.updateSource(newVal)
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
