<script lang="ts" setup>
import {reactive} from 'vue'
import {ListPlayers, LoadHistory} from '../../wailsjs/go/app/App'
import {history, player} from "../../wailsjs/go/models";
import HistoryComponent from "./HistoryComponent.vue";

interface Data {
  players: Array<player.Player>,
  history: Array<history.GameHistory>,
  selectedPlayer?: player.Player
}

const data = reactive<Data>({
  players: [],
  history: [],
  selectedPlayer: undefined
})

ListPlayers().then(result => {
  data.players = result
}).catch(error => console.error(error))

function selectPlayer(p: player.Player) {
  data.selectedPlayer = p
}

function clearPlayerSelection() {
  data.selectedPlayer = undefined
}

</script>

<template>
  <main>
    <nav class="navbar navbar-expand-lg navbar-dark bg-dark sticky-top">
      <div class="container-fluid">
        <button class="navbar-toggler" type="button" data-bs-toggle="collapse" data-bs-target="#players"
                aria-controls="players" aria-expanded="false" aria-label="Toggle navigation">
          <span class="navbar-toggler-icon"></span>
        </button>
        <div class="collapse navbar-collapse" id="players">
          <ul class="navbar-nav">
            <li class="nav-item dropdown">
              <a class="nav-link dropdown-toggle" href="#" id="navbarDarkDropdownMenuLink" role="button"
                 data-bs-toggle="dropdown" aria-expanded="false">
                {{ data.selectedPlayer ? data.selectedPlayer.battleTag : 'Select Player' }}
              </a>
              <ul class="dropdown-menu dropdown-menu-dark" aria-labelledby="navbarDarkDropdownMenuLink">
                <li><a class="dropdown-item" href="#" v-on:click="clearPlayerSelection">Select Player</a></li>
                <li v-for="player in data.players"><a class="dropdown-item" href="#" v-on:click="selectPlayer(player)">{{
                    player.battleTag
                  }}</a>
                </li>
              </ul>
            </li>
          </ul>
        </div>
      </div>
    </nav>
    <HistoryComponent v-if="data.selectedPlayer" :selected-player="data.selectedPlayer"></HistoryComponent>
  </main>
</template>
