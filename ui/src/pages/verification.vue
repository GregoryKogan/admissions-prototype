<template>
  <v-container class="fill-height d-flex align-center justify-center">
    <v-card v-if="token.length == 0">
      <v-alert
        title="Токен верификации не указан"
        text="Перейдите по ссылке из письма"
        variant="tonal"
        type="error"
      ></v-alert>
    </v-card>
    <div v-else>
      <v-card v-if="loading" elevation="0">
        <v-card-title>Подтверждаем почту</v-card-title>
        <div class="d-flex justify-center align-center pa-4">
          <v-progress-circular
            indeterminate
            color="primary"
          ></v-progress-circular>
        </div>
      </v-card>
      <v-card v-else-if="success" elevation="0">
        <v-card-title>Почта подтверждена</v-card-title>
        <v-card-text>
          <v-alert variant="tonal" type="success">
            Вы сможете войти в систему когда учетная запись будет активирована
            администратором.
            <br />
            Письмо с логином и паролем придет на указанную почту.
          </v-alert>
        </v-card-text>
      </v-card>
      <v-card v-else elevation="0">
        <v-card-title>Почта не подтверждена</v-card-title>
        <v-card-text>
          <v-alert variant="tonal" type="error">
            Некорректный токен верификации.
            <br />
            Перейдите по ссылке из последнего письма.
          </v-alert>
        </v-card-text>
      </v-card>
    </div>
  </v-container>
</template>

<script lang="ts">
import RegistrationService from '@/api.registration'
import { defineComponent } from 'vue'

export default defineComponent({
  data: () => ({
    loading: true,
    success: false,
  }),
  async mounted() {
    if (this.token.length == 0) return
    const response = await RegistrationService.verify(this.token)
    this.success = response.status === 200
    this.loading = false
  },
  computed: {
    token() {
      return (this.$route.query.token as string) || ''
    },
  },
})
</script>

<style scoped>
.fill-height {
  height: 100vh;
}
</style>
