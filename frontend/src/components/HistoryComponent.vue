<script lang="ts" setup>
import {reactive} from 'vue'
import {LoadHistory, BackupFolder} from '../../wailsjs/go/app/App'
import {history, player} from "../../wailsjs/go/models";
import HistoryItemComponent from "./HistoryItemComponent.vue";

interface Props {
  selectedPlayer: player.Player
}

const props = defineProps<Props>()

interface Data {
  history: Array<history.GameHistory>,
  backupFolder?: string,
}

const data = reactive<Data>({
  history: [],
})

BackupFolder(props.selectedPlayer).then(result => data.backupFolder=result)

setInterval(() => {
  LoadHistory(props.selectedPlayer).then(result => {
    data.history = result
  }).catch(error => console.error(error))
}, 1000)

</script>

<template>
  <div class="container-fluid">
    <div class="row" v-if="data.backupFolder">
      <div class="col">
        <div class="alert alert-dark fade show" role="alert">
          <h4 class="alert-heading fw-bold">Hi {{props.selectedPlayer.battleTag}}!</h4>
          <p class="mb-0">This application will let you restore your game history to any checkpoint you would want to.</p>
          <p class="mb-0">Upon restoring the application will automatically create backups of your history before executing any operations.</p>
          <p>in case of an error you can restore.</p>
          <hr>
          <p class="mb-0"><strong>The backup location for your account is:</strong><br /><u>{{data.backupFolder}}</u></p>
        </div>
      </div>
    </div>
    <div class="row">
      <table v-if="data.history" class="table table-dark table-striped table-hover m-0">
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
