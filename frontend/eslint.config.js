// For more info, see https://github.com/storybookjs/eslint-plugin-storybook#configuration-flat-config-format
import storybook from "eslint-plugin-storybook";

import pluginVue from 'eslint-plugin-vue'
import vueTsEslintConfig from '@vue/eslint-config-typescript'
import pluginVitest from '@vitest/eslint-plugin'
import skipFormatting from '@vue/eslint-config-prettier/skip-formatting'

export default [{
  name: 'app/files-to-lint',
  files: ['**/*.{ts,mts,tsx,vue}'],
}, {
  name: 'app/files-to-ignore',
  ignores: ['**/dist/**', '**/dist-ssr/**', '**/coverage/**'],
}, ...pluginVue.configs['flat/essential'], ...vueTsEslintConfig(), {
  ...pluginVitest.configs.recommended,
  files: ['src/**/__tests__/*'],
}, skipFormatting, ...storybook.configs["flat/recommended"], {
  name: 'app/storybook-overrides',
  files: ['**/*.stories.ts', '**/*.stories.tsx'],
  rules: {
    '@typescript-eslint/no-explicit-any': 'off', // Allow 'any' in story files for mocking
  }
}, {
  name: 'app/typescript-overrides',
  files: ['**/*.{ts,tsx,vue}'],
  rules: {
    '@typescript-eslint/no-unused-vars': ['error', {
      argsIgnorePattern: '^_',
      varsIgnorePattern: '^_',
      caughtErrorsIgnorePattern: '^_'
    }],
    '@typescript-eslint/ban-ts-comment': ['error', {
      'ts-expect-error': 'allow-with-description',
      'ts-ignore': false,
      'ts-nocheck': true,
      'ts-check': false,
    }],
  }
}];
