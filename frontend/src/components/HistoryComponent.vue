<script lang="ts" setup>
import {reactive} from 'vue'
import {LoadHistory} from '../../wailsjs/go/app/App'
import {history, player} from "../../wailsjs/go/models";

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
      <tr v-for="game in data.history">
        <td class="text-center">{{ game.gameId }}</td>
        <td>{{ game.date }}</td>
        <td>
          {{ game.elo }}
          <span v-if="['WIN', 'LOSS'].includes(game.outcome)" class="badge badge rounded-pill" :class="game.outcome == 'WIN' ? 'bg-success' : 'bg-danger'">{{`${game.outcome == 'WIN' ? '+' : '-'}${game.eloDiff}`}}</span>
        </td>
        <td>
          {{ game.outcome }}
          <span v-if="game.winsStreak > 1">&#128293; {{game.winsStreak}}</span>
        </td>
      </tr>
      </tbody>
    </table>
  </main>
</template>
