<template>
  <main>
    <component-header />
  </main>
</template>

<script>
import ComponentHeader from "@/components/ComponentHeader";
export default {
  name: "App",
  components: { ComponentHeader }
};
</script>

<style>
#app {
  font-family: Avenir, Helvetica, Arial, sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  text-align: center;
  color: #2c3e50;
  margin-top: 60px;
}
</style>
