<template>
  <div class="flex-col items-center">
    <h1 class="text-center text-gray-700 text-6xl">
      <b>Endpoint #2</b>
    </h1>
    <p
      class="text-gray-600 text-xl text-center"
    >Aquí podrás ver los dominios que han sido consultados en el Endpoint #1</p>
    <component-button class="mt-5">
      <slot>Consultar dominios</slot>
    </component-button>
    <router-link
      to="/"
      class="mt-5 text-xl text-green-600 hover:underline"
    >Volver a la página de Inicio (Home)</router-link>
  </div>
</template>

<script>
import ComponentButton from "@/components/ComponentButton";
export default {
  name: "Endpoint2",
  components: { ComponentButton }
};
</script>