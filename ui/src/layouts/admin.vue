<template>
  <v-app>
    <v-navigation-drawer permanent>
      <v-list>
        <v-list-item
          to="/admin/dashboard"
          prepend-icon="mdi-view-dashboard"
          title="Панель"
        />
        <v-list-item
          to="/admin/registrations"
          prepend-icon="mdi-format-list-bulleted"
          title="Регистрации"
        />
      </v-list>
    </v-navigation-drawer>

    <v-app-bar>
      <v-app-bar-title>Панель администратора</v-app-bar-title>
      <v-spacer></v-spacer>
      <v-btn icon="mdi-logout" variant="text"></v-btn>
    </v-app-bar>

    <v-main>
      <router-view />
    </v-main>
  </v-app>
</template>

<script lang="ts" setup>
//
</script>
