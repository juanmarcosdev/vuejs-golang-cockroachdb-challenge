<template>
  <div class="flex-col items-center">
    <h1 class="text-center text-gray-700 text-6xl">
      <b>Endpoint #2</b>
    </h1>
    <p
      class="text-gray-600 text-xl text-center"
    >Aquí podrás ver los dominios que han sido consultados en el Endpoint #1</p>
    <component-button class="mt-5" @custom-click="doQuery">
      <slot>Consultar dominios</slot>
    </component-button>
    <div v-if="showInformation" class="m-15">
      <p class="text-gray-600 text-2xl text-center">
        A continuación se despliegan
        <b>todos</b>
        los dominios que han sido consultados:
      </p>
    </div>
    <table v-if="showInformation">
      <thead>
        <tr class="bg-gray-100 border-b-2 border-gray-400">
          <th>Dominios</th>
        </tr>
      </thead>
      <tbody>
        <tr
          v-for="i in items"
          :key="i.domain"
          class="border-b border-gray-200 hover:bg-gray-100 hover:bg-orange-100"
        >
          <td>{{ i.domain }}</td>
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
  name: "Endpoint2",
  components: { ComponentButton },
  data() {
    return {
      items: [],
      showInformation: false
    };
  },
  methods: {
    doQuery() {
      Swal.fire({
        title: "Se van a desplegar todos los dominios consultados",
        showCancelButton: true,
        confirmButtonText: "Desplegar",
        cancelButtonText: "Cancelar",
        showLoaderOnConfirm: true,
        preConfirm: () => {
          return axios
            .get(`http://localhost:5000/queried/domains`)
            .then(response => {
              this.items = response.data.items;
              this.showInformation = true;
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