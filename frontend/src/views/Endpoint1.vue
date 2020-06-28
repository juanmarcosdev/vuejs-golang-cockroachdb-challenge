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
    <component-button class="m-10" @custom-click="doQuery">
      <slot>Realizar consulta</slot>
    </component-button>
    <div v-if="showInformation" class="m-15">
      <p class="text-gray-600 text-2xl text-center">
        Resultado de consulta sobre:
        <i>
          <b>{{ domain }}</b>
        </i>
      </p>
    </div>
    <table v-if="showInformation">
      <thead>
        <tr class="bg-gray-100 border-b-2 border-gray-400">
          <th>Dirección IP</th>
          <th>Grado SSL</th>
          <th>País</th>
          <th>Owner</th>
        </tr>
      </thead>
      <tbody>
        <tr
          v-for="s in servers"
          :key="s.address"
          class="border-b border-gray-200 hover:bg-gray-100 hover:bg-orange-100"
        >
          <td>{{ s.address }}</td>
          <td>{{ s.ssl_grade }}</td>
          <td>{{ s.country }}</td>
          <td>{{ s.owner }}</td>
        </tr>
      </tbody>
    </table>
    <table v-if="showInformation">
      <thead>
        <tr class="bg-gray-100 border-b-2 border-gray-400">
          <th>Los servidores han cambiado?</th>
          <th>Grado SSL (más bajo)</th>
          <th>Anterior Grado SSL</th>
          <th>Logo</th>
          <th>Title</th>
          <th>Está caído?</th>
        </tr>
      </thead>
      <tbody>
        <tr class="border-b border-gray-200 hover:bg-gray-100 hover:bg-orange-100">
          <td>{{ servers_changed }}</td>
          <td>{{ ssl_grade }}</td>
          <td>{{ previous_ssl_grade }}</td>
          <td>{{ logo }}</td>
          <td>{{ title }}</td>
          <td>{{ is_down }}</td>
        </tr>
      </tbody>
    </table>
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
      domain: "",
      servers: [],
      servers_changed: false,
      ssl_grade: "",
      previous_ssl_grade: "",
      logo: "",
      title: "",
      is_down: false,
      showInformation: false
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
        cancelButtonText: "Cancelar",
        showLoaderOnConfirm: true,
        preConfirm: domain => {
          if (domain === "") {
            Swal.fire(
              "Dominio vacío",
              "No has ingresado un dominio para consultar",
              "error"
            );
          } else {
            return axios
              .post(`http://localhost:5000/info/servers/${domain}`)
              .then(response => {
                this.domain = domain;
                this.servers = response.data.servers;
                this.servers_changed = response.data.servers_changed;
                this.ssl_grade = response.data.ssl_grade;
                this.previous_ssl_grade = response.data.previous_ssl_grade;
                this.logo = response.data.logo;
                this.title = response.data.title;
                this.is_down = response.data.is_down;
                this.showInformation = true;
              })
              .catch(error => {
                Swal.showValidationMessage(`Falló la solicitud: ${error}`);
              });
          }
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

<style scoped>
td {
  padding: 10px;
  text-align: center;
}
</style>