<script lang="ts" setup>
import {reactive} from 'vue'
import {ListPlayers} from '../../wailsjs/go/app/App'
import {player} from "../../wailsjs/go/models";
import HistoryComponent from "./HistoryComponent.vue";

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
<!--    <div class="container">-->
<!--      <div class="row">-->
<!--        <div class="col-8">-->
<!--          <HistoryComponent v-if="data.selectedPlayer" :selected-player="data.selectedPlayer"></HistoryComponent>-->
<!--        </div>-->
<!--        <div class="col-4">-->
<!--          <div class="card text-white bg-dark shadow">-->
<!--            <div class="card-header bg-gradient">Header</div>-->
<!--            <div class="card-body">-->
<!--              <h5 class="card-title">Secondary card title</h5>-->
<!--              <p class="card-text">Some quick example text to build on the card title and make up the bulk of the card's content.</p>-->
<!--            </div>-->
<!--          </div>-->
<!--        </div>-->
<!--      </div>-->
<!--    </div>-->
  </main>
</template>
