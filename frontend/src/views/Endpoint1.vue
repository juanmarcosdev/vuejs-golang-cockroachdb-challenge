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
    <component-button class="mt-5" @custom-click="doQuery">
      <slot>Realizar consulta</slot>
    </component-button>
    <router-link
      to="/"
      class="mt-5 text-xl text-green-600 hover:underline"
    >Volver a la página de Inicio (Home)</router-link>
  </div>
</template>

<script>
import ComponentButton from "@/components/ComponentButton";
import Swal from "sweetalert2";
import axios from "axios";
export default {
  name: "Endpoint1",
  components: { ComponentButton },
  data() {
    return {
      title: ""
    };
  },
  methods: {
    doQuery() {
      Swal.fire({
        title: "Ingrese el dominio a consultar",
        input: "text",
        inputAttributes: {
          autocapitalize: "off"
        },
        showCancelButton: true,
        confirmButtonText: "Enviar",
        showLoaderOnConfirm: true,
        preConfirm: domain => {
          return axios
            .post(`http://localhost:5000/info/servers/${domain}`)
            .then(response => {
              this.title = response.data.title;
            })
            .catch(error => {
              Swal.showValidationMessage(`Falló la solicitud: ${error}`);
            });
        },
        allowOutsideClick: () => !Swal.isLoading()
      }).then(result => {
        if (result.value) {
          Swal.fire("Éxito", "Consulta realizada exitosamente", "success");
        }
      });
    }
  }
};
</script>