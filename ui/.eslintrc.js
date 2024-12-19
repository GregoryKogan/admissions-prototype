/**
 * .eslint.js
 *
 * ESLint configuration file.
 */

// eslint-disable-next-line no-undef
module.exports = {
  root: true,
  env: {
    node: true,
  },
  extends: [
    'vuetify',
    '@vue/eslint-config-typescript',
    './.eslintrc-auto-import.json',
  ],
  rules: {
    'vue/multi-word-component-names': 'off',
  },
}
