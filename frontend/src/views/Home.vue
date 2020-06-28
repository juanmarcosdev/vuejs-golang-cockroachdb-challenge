<template>
  <div class="flex-col items-center">
    <h1 class="text-center text-gray-700 text-6xl">
      Vue.js + Golang + CockroachDB
      <i>Challenge</i> v1.0.0
    </h1>
    <p class="text-gray-600 text-xl text-center">
      Miniproyecto, Challenge, Prueba que involucra conocimientos de Frontend, Backend y DB.
      <br />En pocas palabras el proyecto obtiene información de dominios y sus servidores (endpoints) de una API externa, a su vez con estos realiza persistencia de datos y clasificación entre ellos. También hace un poco de web scrapping.
    </p>

    <router-link
      to="/endpoint1/post"
      class="mt-5 text-xl text-green-600 hover:underline"
    >Ir a Endpoint #1</router-link>
    <router-link
      to="/endpoint2/get"
      class="mt-5 text-xl text-green-600 hover:underline"
    >Ir a Endpoint #2</router-link>
  </div>
</template>