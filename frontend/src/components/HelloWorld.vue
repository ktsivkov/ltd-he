<script lang="ts" setup>
import {reactive} from 'vue'
import {ListPlayers, LoadHistory} from '../../wailsjs/go/app/App'
import {player, history} from "../../wailsjs/go/models";

interface Data {
  players: Array<player.Player>,
  history: Array<history.GameHistory>,
}

const data = reactive<Data>({
  players: [],
  history: [],
})

ListPlayers().then(result => {
  data.players = result
}).catch(error => console.error(error))

function selectPlayer(p: player.Player) {
  LoadHistory(p).then(result => {
    data.history = result
  }).catch(error => console.error(error))
}

</script>

<template>
  <main>
    <div v-for="player in data.players">
      <h1 v-on:click="selectPlayer(player)">{{player.battleTag}}</h1>
    </div>
    <div v-for="game in data.history">
      {{game.date}}
      {{game.eloDiff}}
      {{game.outcome}}
    </div>
  </main>
</template>

<style scoped>

</style>
