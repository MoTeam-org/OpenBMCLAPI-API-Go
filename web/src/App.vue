<template>
  <a-config-provider :locale="locale">
    <router-view></router-view>
  </a-config-provider>
</template>

<script setup lang="ts">
import zhCN from 'ant-design-vue/locale/zh_CN'

const locale = zhCN
</script>

<style>
#app {
  min-height: 100vh;
  background: #f0f2f5;
}
</style> 