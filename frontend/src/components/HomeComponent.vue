<script lang="ts" setup>
import {reactive} from 'vue'
import {ListPlayers} from '../../wailsjs/go/app/App'
import {player} from "../../wailsjs/go/models";
import HistoryComponent from "./HistoryComponent.vue";
import SystemAlertsComponent from "./SystemAlertsComponent.vue";
import InsertGameComponent from "./InsertGameComponent.vue";
import AppInfoComponent from "./AppInfoComponent.vue";

interface Data {
  players: Array<player.Player>,
  selectedPlayer?: player.Player
}

const data = reactive<Data>({
  players: [],
  selectedPlayer: undefined
})

setInterval(() => {
  ListPlayers().then(result => {
    data.players = result
  }).catch(error => console.error(error))
}, 1000)

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
                <li><a class="dropdown-item" href="#" v-on:click="data.selectedPlayer = undefined">Select Player</a></li>
                <li v-for="player in data.players"><a class="dropdown-item" href="#" v-on:click="data.selectedPlayer = player">{{
                    player.battleTag
                  }}</a>
                </li>
              </ul>
            </li>
          </ul>
        </div>
      </div>
    </nav>
    <SystemAlertsComponent></SystemAlertsComponent>
    <div class="container-fluid" v-if="data.selectedPlayer">
      <div class="d-flex align-items-start">
        <div class="nav flex-column nav-pills col-2" id="v-pills-tab" role="tablist" aria-orientation="vertical">
          <button class="nav-link active" id="v-pills-home-tab" data-bs-toggle="pill" data-bs-target="#v-pills-home" type="button" role="tab" aria-controls="v-pills-home" aria-selected="true">Home</button>
          <button class="nav-link" id="v-pills-history-tab" data-bs-toggle="pill" data-bs-target="#v-pills-history" type="button" role="tab" aria-controls="v-pills-history" aria-selected="false">History</button>
          <button class="nav-link" id="v-pills-insert-game-tab" data-bs-toggle="pill" data-bs-target="#v-pills-insert-game" type="button" role="tab" aria-controls="v-pills-insert-game" aria-selected="false">Insert Game</button>
        </div>
        <div class="tab-content col-10 p-2" id="v-pills-tabContent">
          <div class="tab-pane fade show active" id="v-pills-home" role="tabpanel" aria-labelledby="v-pills-home-tab" tabindex="0"><AppInfoComponent :selected-player="data.selectedPlayer"></AppInfoComponent></div>
          <div class="tab-pane fade" id="v-pills-history" role="tabpanel" aria-labelledby="v-pills-history-tab" tabindex="0"><HistoryComponent :selected-player="data.selectedPlayer"></HistoryComponent></div>
          <div class="tab-pane fade" id="v-pills-insert-game" role="tabpanel" aria-labelledby="v-pills-insert-game-tab" tabindex="0"><InsertGameComponent :selected-player="data.selectedPlayer"></InsertGameComponent></div>
        </div>
      </div>
    </div>
  </main>
</template>
