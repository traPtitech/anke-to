module.exports = {
  root: true,
  parser: 'vue-eslint-parser',
  parserOptions: {
    parser: 'babel-eslint',
    sourceTye: 'module'
  },
  env: {
    browser: true,
    node: true
  },
  extends: [
    'eslint:recommended',
    'plugin:prettier/recommended',
    'plugin:vue/recommended',
    'prettier/vue'
  ],
  rules: {
    'no-plusplus': 'off',
    'no-console': 'off',
    'func-names': 'off',
    'vue/no-template-shadow': 0,
    'vue/no-unused-components': [
      'error',
      {
        // suppresses all errors if binding has been detected in the template
        ignoreWhenBindingPresent: true
      }
    ],
    'prettier/prettier': [
      'error',
      {
        singleQuote: true,
        semi: false,
        bracketSpacing: true,
        tabWidth: 2
      }
    ]
  }
}
