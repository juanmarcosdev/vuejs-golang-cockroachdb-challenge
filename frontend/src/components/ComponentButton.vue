<template>
  <button
    @click="buttonClick"
    class="bg-transparent hover:bg-purple-700 text-purple-700 font-semibold hover:text-white py-2 px-2 border border-purple-600 hover:border-transparent rounded"
  >
    <p>
      <slot></slot>
    </p>
  </button>
</template>

<script>
export default {
  name: "ComponentButton",
  methods: {
    buttonClick() {
      this.$emit("custom-click");
    }
  }
};
</script>