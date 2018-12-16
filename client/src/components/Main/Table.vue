<template>
  <div class="box">
    <table class="table is-striped">
      <thead>
        <tr>
          <th v-for="(header, index) in headers" :key="index">{{ header }}</th>
        </tr>
      </thead>
      <tbody>
        <tr v-for="(itemrow, index) in itemrows" :key="index">
          <td v-for="(item, index) in itemrow" :key="index" v-html="item"></td>
        </tr>
      </tbody>
    </table>
  </div>
</template>

<script>
export default {
  name: 'Table',
  props: {
    headers: {
      type: Array,
      required: true
    },
    itemrows: {
      type: Array,
      required: true
    }
  }
}
</script>
