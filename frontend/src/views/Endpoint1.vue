<template>
  <div class="flex-col items-center">
    <h1 class="text-center text-gray-700 text-6xl">
      <b>Endpoint #1</b>
    </h1>
    <p class="text-gray-600 text-xl text-center">
      En este Endpoint ingresas un dominio como
      <i>google.com</i> y obtienes la información acerca de sus servidores,
      y a su vez información acerca de ellos, como su IP, grado de seguridad SSL, país y organización (owner) que lo administra.
      También podrás ver si los servidores han cambiado desde su último registro de 1 hora o más antes (si existe) y de su grado SSL menor,
      el grado SSL del anterior registro, junto con su title en el HTML y su posible link de logo.
    </p>
    <router-link
      to="/"
      class="mt-5 text-xl text-green-600 hover:underline"
    >Volver a la página de Inicio (Home)</router-link>
  </div>
</template>