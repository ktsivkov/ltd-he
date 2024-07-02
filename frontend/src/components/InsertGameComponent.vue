<script lang="ts" setup>
import {reactive} from 'vue'
import {Insert} from '../../wailsjs/go/app/App'
import {history, player} from "../../wailsjs/go/models";

interface Props {
  selectedPlayer: player.Player
}

const props = defineProps<Props>()

interface Data {
  insertionInProgress: boolean
  input: history.InsertRequest
}

const data = reactive<Data>({
  insertionInProgress: false,
  input: new history.InsertRequest()
})

function insertGame() {
  data.insertionInProgress = true
  Insert(props.selectedPlayer, data.input).then(()=>{data.input = new history.InsertRequest()}).finally(() => {
    data.insertionInProgress = false
  })
}

</script>

<template>
  <div class="row">
    <div class="col-6">
      <fieldset v-bind:disabled="data.insertionInProgress">
        <form class="row gy-2 gx-3 align-items-center">
          <div class="col">
            <label class="visually-hidden" for="elo">ELO</label>
            <input type="number" class="form-control" id="elo" placeholder="ELO" v-model="data.input.elo">
          </div>
          <div class="col-auto">
            <div class="form-check">
              <input class="form-check-input" type="checkbox" id="mvp" v-model="data.input.mvp">
              <label class="form-check-label" for="mvp">
                MVP
              </label>
            </div>
          </div>
          <div class="col-auto">
            <button type="submit" class="btn btn-outline-success" v-on:click="insertGame">Insert Game</button>
          </div>
        </form>
      </fieldset>
    </div>
  </div>
</template>
