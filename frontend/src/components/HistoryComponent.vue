<script lang="ts" setup>
import {reactive} from 'vue'
import {LoadHistory} from '../../wailsjs/go/app/App'
import {history, player} from "../../wailsjs/go/models";
import HistoryItemComponent from "./HistoryItemComponent.vue";

interface Props {
  selectedPlayer: player.Player
}

const props = defineProps<Props>()

interface Data {
  history: Array<history.GameHistory>,
}

const data = reactive<Data>({
  history: []
})

LoadHistory(props.selectedPlayer).then(result => {
  data.history = result
}).catch(error => console.error(error))

</script>

<template>
  <main>
    <table v-if="data.history" class="table table-dark table-striped table-hover m-0">
      <thead>
      <tr>
        <td class="text-center">#</td>
        <td>Date</td>
        <td>Elo</td>
        <td>Outcome</td>
      </tr>
      </thead>
      <tbody>
      <HistoryItemComponent v-for="game in data.history" :game="game"></HistoryItemComponent>
      </tbody>
    </table>
  </main>
</template>
