<script lang="ts" setup>
import {reactive} from 'vue'
import {LoadHistory} from '../../wailsjs/go/app/App'
import {history, player} from "../../wailsjs/go/models";
import HistoryItemComponent from "./HistoryItemComponent.vue";
import AppendGameComponent from "./AppendGameComponent.vue";

interface Props {
  selectedPlayer: player.Player
}

const props = defineProps<Props>()

interface Data {
  history: Array<history.GameHistory>
}

const data = reactive<Data>({
  history: []
})

setInterval(() => {
  LoadHistory(props.selectedPlayer).then(result => {
    data.history = result
  }).catch(error => console.error(error))
}, 1000)

</script>

<template>
  <AppendGameComponent :selected-player="props.selectedPlayer"></AppendGameComponent>
  <div class="row" v-if="data.history">
    <div class="col">
      <table class="table table-dark table-striped table-hover m-0">
        <thead>
        <tr>
          <td class="text-center">Outcome</td>
          <td>ELO</td>
          <td>Date</td>
          <td>Options</td>
        </tr>
        </thead>
        <tbody>
        <HistoryItemComponent v-for="game in data.history" :game="game"></HistoryItemComponent>
        </tbody>
      </table>
    </div>
  </div>
</template>
